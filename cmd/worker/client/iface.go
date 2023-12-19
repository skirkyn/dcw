package client

import "github.com/skirkyn/dcw/cmd/dto"

type Client[Req dto.Request[Req], Resp dto.Response[Resp]] interface {
	Send(Req) Resp
}
