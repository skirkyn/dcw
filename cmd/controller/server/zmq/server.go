package zmq

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"github.com/skirkyn/dcw/cmd/cmn/nw"
	"github.com/skirkyn/dcw/cmd/cmn/util"
	"github.com/skirkyn/dcw/cmd/controller/server"
	"log"
	"time"
)

const internalAddress = "inproc://backend"

type ServerConfig struct {
	Workers                               int
	Port                                  int
	MaxSendResponseRetries                int
	TimeToSleepBetweenSendResponseRetries time.Duration
}
type Server struct {
	handler      server.Handler
	shouldStop   *chan bool
	stopError    *chan error
	serverConfig ServerConfig
}

func NewZMQServer(handler server.Handler,
	serverConfig ServerConfig) server.Server {
	shouldStop := make(chan bool)
	stopped := make(chan error)
	return &Server{handler, &shouldStop, &stopped, serverConfig}
}

func (s *Server) Start() error {
	frontend, err := zmq4.NewSocket(zmq4.ROUTER)
	if err != nil {
		log.Printf("can't create socket %s", err.Error())
		return err
	}
	defer nw.CloseSocket(frontend)

	err = frontend.Bind(fmt.Sprintf("tcp://*:%d", s.serverConfig.Port))
	if err != nil {
		log.Printf("can't bind to the socket %s", err.Error())
		return err
	}

	backend, err := zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		log.Printf("can't create socket %s", err.Error())
		return err
	}
	err = backend.Bind(internalAddress)

	if err != nil {
		log.Printf("can't bind to the socket %s", err.Error())
		return err
	}
	for i := 0; i < s.serverConfig.Workers; i++ {
		go s.startWorker(internalAddress, backend)
	}

	return zmq4.Proxy(frontend, backend, nil)
}

func (s *Server) Stop() error {
	*s.shouldStop <- true
	return <-*s.stopError
}

func (s *Server) startWorker(internalAddress string, backend *zmq4.Socket) {
	worker, err := zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		log.Printf("can't create socket %s", err.Error())
		return
	}
	defer nw.CloseSocket(worker)
	err = worker.Connect(internalAddress)
	for {
		s.maybeProcessMessage(worker)

		select {

		case shouldStop := <-*s.shouldStop:
			if shouldStop {
				log.Print("stopping the server")
				s.stop(backend)
				return
			}
		}
	}
}

func (s *Server) maybeProcessMessage(worker *zmq4.Socket) {
	request, err := worker.RecvMessage(0)
	if err != nil {
		log.Print(err)
		return
	}
	if len(request) != 2 {
		log.Print("invalid request ", request)
		return
	}

	client, content := s.parseMessage(request)

	respChan := make(chan []byte)
	errChan := make(chan error)
	go s.handleRequest(content, &respChan, &errChan)
	go s.sendResponse(worker, client, &respChan, &errChan)
}

func (s *Server) toByteArray(input []string) []byte {
	res := make([]byte, 0)

	if input == nil {
		return res
	}

	for i := 0; i < len(input); i++ {
		res = append(res, []byte((input)[i])...)
	}

	return res
}
func (s *Server) parseMessage(msg []string) (string, []byte) {
	index := 1

	if msg[1] == "" {
		index++

	}
	return string(s.toByteArray(msg[:2])), s.toByteArray(msg[2:])
}
func (s *Server) handleRequest(data []byte, respChannel *chan []byte, errChannel *chan error) {
	s.handler.Handle(data, respChannel, errChannel)
}

func (s *Server) sendResponse(router *zmq4.Socket, client string, respChannel *chan []byte, errChannel *chan error) bool {
	for {
		select {
		case resp := <-*respChannel:
			return s.respond(router, client, resp)
		case err := <-*errChannel:
			return s.respond(router, client, util.StrToByteSlice(err.Error()))
		}
	}
}

func (s *Server) respond(router *zmq4.Socket, client string, resp []byte) bool {
	for i := s.serverConfig.MaxSendResponseRetries; i > 0; i-- {

		_, err := router.SendMessage(client, resp)
		if err != nil {
			log.Printf("couldn't send the response, will retry %s", err.Error())
			time.Sleep(s.serverConfig.TimeToSleepBetweenSendResponseRetries * time.Second)
		} else {
			return true
		}
	}
	return false
}

func (s *Server) stop(backend *zmq4.Socket) {
	*s.stopError <- backend.Close()
}
