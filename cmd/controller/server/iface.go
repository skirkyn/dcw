package server

import (
	"github.com/skirkyn/dcw/cmd/dto"
)

type Server interface {
	Start() error
	Stop() error
}

type Handler interface {
	Handle([]byte, *chan []byte, *chan error)
}

type Processor[In any, Out any] interface {
	Process(In) dto.Response[Out]
}
