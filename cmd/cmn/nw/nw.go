package nw

import (
	"github.com/pebbe/zmq4"
	"log"
)

func CloseSocket(socket *zmq4.Socket) {
	err := socket.Close()
	if err != nil {
		log.Printf("can't close the socket %s", err.Error())
	}
}
