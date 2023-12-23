package zmq

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"github.com/skirkyn/dcw/cmd/controller/handler"
	"github.com/skirkyn/dcw/cmd/controller/server"
	"github.com/skirkyn/dcw/cmd/util/bytz"
	"github.com/skirkyn/dcw/cmd/util/socket"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const internalAddress = "inproc://backend"

type Config struct {
	Workers                               int
	Port                                  int
	MaxSendResponseRetries                int
	TimeToSleepBetweenSendResponseRetries time.Duration
}

type Server struct {
	handler         handler.Handler
	stopped         *atomic.Bool
	pendingRequests *sync.WaitGroup
	serverConfig    Config
	backend         *zmq4.Socket
	frontend        *zmq4.Socket
}

func NewServer(handler handler.Handler,
	serverConfig Config) (server.Server, error) {
	backend, err := zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		log.Printf("can't create socket %s", err.Error())
		return nil, err
	}
	frontend, err := zmq4.NewSocket(zmq4.ROUTER)
	if err != nil {
		log.Printf("can't create socket %s", err.Error())
		return nil, err
	}
	return &Server{handler, &atomic.Bool{}, &sync.WaitGroup{}, serverConfig, backend, frontend}, nil
}

func (s *Server) Start() error {

	err := s.frontend.Bind(fmt.Sprintf("tcp://*:%d", s.serverConfig.Port))
	if err != nil {
		log.Printf("can't bind to the socket %s", err.Error())
		return err
	}

	if err != nil {
		log.Printf("can't create socket %s", err.Error())
		return err
	}
	err = s.backend.Bind(internalAddress)

	if err != nil {
		log.Printf("can't bind to the socket %s", err.Error())
		return err
	}
	for i := 0; i < s.serverConfig.Workers; i++ {
		go s.startWorker(internalAddress)
	}

	return zmq4.Proxy(s.frontend, s.backend, nil)
}

func (s *Server) Stop() error {
	s.stopped.Store(true)
	s.pendingRequests.Wait()

	err := s.frontend.Close()
	if err != nil {
		log.Printf("can't close socket %s", err.Error())
	}
	err = s.backend.Close()
	if err != nil {
		log.Printf("can't close socket %s", err.Error())

	}
	return err
}

func (s *Server) startWorker(internalAddress string) {
	worker, err := zmq4.NewSocket(zmq4.DEALER)
	lock := sync.Mutex{}
	if err != nil {
		log.Printf("can't create socket %s", err.Error())
		return
	}
	defer socket.CloseSocket(worker)
	err = worker.Connect(internalAddress)
	if err != nil {
		log.Printf("can't connect worker socket %s", err.Error())
		return
	}
	for {
		if s.stopped.Load() {
			return
		}
		go s.maybeProcessMessage(worker, &lock)
	}
}

func (s *Server) maybeProcessMessage(worker *zmq4.Socket, lock *sync.Mutex) {

	lock.Lock()
	request, err := worker.RecvMessage(zmq4.DONTWAIT)
	lock.Unlock()

	if err != nil {
		log.Print(err)
		return
	}

	if len(request) != 2 {
		log.Print("invalid request ", request)
		return
	}

	client, content := s.parseMessage(request)

	res := s.handler.Handle(content)

	s.respond(worker, lock, client, res)
}

func (s *Server) parseMessage(msg []string) (string, []byte) {
	index := 1

	if msg[1] == "" {
		index++

	}
	return string(bytz.SliceToByteSlice(msg[:2])), bytz.SliceToByteSlice(msg[2:])
}

func (s *Server) respond(router *zmq4.Socket, lock *sync.Mutex, client string, resp []byte) {
	for i := s.serverConfig.MaxSendResponseRetries; i > 0; i-- {

		lock.Lock()
		_, err := router.SendMessage(client, resp)
		lock.Unlock()

		if err == nil {
			return
		}

		log.Printf("couldn't send the response, will retry %s", err.Error())
		time.Sleep(s.serverConfig.TimeToSleepBetweenSendResponseRetries * time.Second)

	}
}
