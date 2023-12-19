package server

import (
	"github.com/skirkyn/dcw/cmd/dto"
)

type Server[Req dto.Request[any], Resp dto.Response[any]] interface {
	Start(port int, host string) error
	Stop()
}

type Handler[Req dto.Request[any], Resp dto.Response[any]] interface {
	Handle(Req, *chan Resp) error
}
