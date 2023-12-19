package zmq

import (
	"fmt"
	"github.com/skirkyn/dcw/cmd/controller/server"
	"github.com/skirkyn/dcw/cmd/dto"
	"github.com/zeromq/goczmq"
	"log"
)

type ZMQServer[Req dto.Request[any], Resp dto.Response[any]] struct {
	handler             server.Handler[Req, Resp]
	requestTransformer  dto.RequestTransformer[Req]
	responseTransformer dto.ResponseTransformer[Resp]
	clients             *map[string]*chan Resp
	stop                *chan bool
}

func NewZMQServer[Req dto.Request[any], Resp dto.Response[any]](handler server.Handler[Req, Resp],
	requestTransformer dto.RequestTransformer[Req],
	responseTransformer dto.ResponseTransformer[Resp]) server.Server[Req, Resp] {
	clients := make(map[string]*chan Resp)
	stop := make(chan bool)
	return &ZMQServer[Req, Resp]{handler: handler, requestTransformer: requestTransformer, responseTransformer: responseTransformer, clients: &clients, stop: &stop}
}

func (s *ZMQServer[Req, Resp]) Start(port int, host string) error {
	router, err := goczmq.NewRouter(fmt.Sprintf("tcp://*:%d", port))
	if err != nil {
		log.Printf("can't start the server %s", err.Error())
		return err
	}
	dealer, err := goczmq.NewDealer(fmt.Sprintf("tcp://%s:%d", host, port))
	if err != nil {
		log.Printf("can't start the server %s", err.Error())
		return err
	}
	go s.startInternal(dealer, router)
	return nil
}

func (s *ZMQServer[Req, Resp]) Stop() {
	*s.stop <- true
}

func (s *ZMQServer[Req, Resp]) startInternal(dealer *goczmq.Sock, router *goczmq.Sock) {
	defer router.Destroy()
	for {
		s.maybeProcessMessage(dealer, router)

		select {

		case shouldStop := <-*s.stop:
			if shouldStop {
				log.Print("stopping the server")
				return
			}
		}
	}
}

func (s *ZMQServer[Req, Resp]) maybeProcessMessage(dealer *goczmq.Sock, router *goczmq.Sock) {
	request, err := dealer.RecvMessage()
	if err != nil {
		log.Print(err)
		return
	}
	if len(request) != 2 {
		log.Print("invalid request ", request)
		return
	}

	client := string(request[0])
	channel, ok := (*s.clients)[client]

	if !ok {
		*(*s.clients)[client] = make(chan Resp)
		channel = (*s.clients)[client]
	}

	go s.handleRequest(request[1], channel)
	go s.sendResponse(router, client)
}

func (s *ZMQServer[Req, Resp]) handleRequest(data []byte, respChannel *chan Resp) {
	transformed, err := s.requestTransformer.Transform(data)
	if err != nil {
		*respChannel <- dto.NewErrorResponse(err.Error())
		return
	}
	err = s.handler.Handle(transformed, respChannel)
}

func (s *ZMQServer[Req, Resp]) sendResponse(router *goczmq.Sock, client string) {
	channel, ok := (*s.clients)[client]

	if !ok {
		log.Printf("can't send response to the client, client is lost %s", client)
		return
	}

	resp := <-*channel
	transformed, err := s.responseTransformer.Transform(resp)
	err = router.SendFrame([]byte(client), goczmq.FlagMore)
	if err != nil {
		log.Printf("can't send response to the client %s", client)
		return
	}
	err = router.SendFrame(transformed, goczmq.FlagNone)
	if err != nil {
		log.Print(err)
	}
}
