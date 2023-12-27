package zmq

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"github.com/unknownfeature/dcw/cmd/common"
	"github.com/unknownfeature/dcw/cmd/controller/server"
	"github.com/unknownfeature/dcw/cmd/util"
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
	handler         common.Function[[]byte, []byte]
	stopped         *atomic.Bool
	pendingRequests *sync.WaitGroup
	serverConfig    Config
	backend         *zmq4.Socket
	frontend        *zmq4.Socket
}

func NewServer(handler common.Function[[]byte, []byte],
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

func (s *Server) Start() (*sync.WaitGroup, error) {

	err := s.frontend.Bind(fmt.Sprintf("tcp://*:%d", s.serverConfig.Port))

	if err != nil {
		log.Printf("can't create socket %s", err.Error())
		return nil, err
	}
	err = s.backend.Bind(internalAddress)

	if err != nil {
		log.Printf("can't bind to the socket %s", err.Error())
		return nil, err
	}
	wg := sync.WaitGroup{}
	wg.Add(s.serverConfig.Workers)
	for i := 0; i < s.serverConfig.Workers; i++ {
		go s.startWorker(internalAddress, &wg)
	}

	return &wg, zmq4.Proxy(s.frontend, s.backend, nil)
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

func (s *Server) startWorker(internalAddress string, wg *sync.WaitGroup) {
	worker, err := zmq4.NewSocket(zmq4.DEALER)
	lock := sync.Mutex{}
	if err != nil {
		log.Printf("can't create socket %s", err.Error())
		return
	}
	defer util.CloseSocket(worker)
	defer wg.Done()
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
	request, err := worker.RecvMessage(0)

	if err != nil {
		return
	}

	if len(request) != 2 {
		log.Print("invalid request ", request)
		lock.Unlock()
		return
	}

	client, content := s.parseMessage(request)

	res, err := s.handler.Apply(content)
	if err != nil {
		s.respond(worker, client, res)
		lock.Unlock()
		log.Printf("error handling request %s", err.Error())
	}
	s.respond(worker, client, res)
	lock.Unlock()
}

func (s *Server) parseMessage(msg []string) (string, []byte) {
	return string(util.SliceToByteSlice(msg[:1])), util.SliceToByteSlice(msg[1:])
}

func (s *Server) respond(router *zmq4.Socket, client string, resp []byte) {
	for i := s.serverConfig.MaxSendResponseRetries; i > 0; i-- {
		_, err := router.SendMessage(client, resp)

		if err == nil {
			return
		}

		log.Printf("couldn't send the response, will retry %s", err.Error())
		time.Sleep(s.serverConfig.TimeToSleepBetweenSendResponseRetries * time.Second)

	}
}
