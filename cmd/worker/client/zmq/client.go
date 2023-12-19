package zmq

import (
	"fmt"
	"github.com/skirkyn/dcw/cmd/dto"
	"github.com/zeromq/goczmq"
	"log"
	"time"
)

type ZMQClient[Req dto.Request[Req], Resp dto.Response[Resp]] struct {
	requestTransformer  dto.RequestTransformer[Req]
	responseTransformer dto.RequestTransformer[Req]
	socket              *goczmq.Sock
}

func NewZMQClient(port int, host string, maxAttemptsToReconnect int) error {
	var socket goczmq.Sock
	for i := 0; i < maxAttemptsToReconnect; i++ {
		socket, err := goczmq.NewDealer(fmt.Sprintf("tcp://%s:%d", host, port))
		if err != nil {
			log.Printf("can't connect to the server %s:%d %s", host, port, err)
			if i == maxAttemptsToReconnect-1 {
				return err
			}
		} else {
			break
		}
		time.Sleep(5 * time.Second)
	}

}
