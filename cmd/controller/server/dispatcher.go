package server

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/common/dto"
)

type Dispatcher struct {
	handlers           map[dto.Type]common.Function[dto.Request[any], []byte]
	requestTransformer dto.RequestTransformer[any]
}

func NewDispatcher(handlers map[dto.Type]common.Function[dto.Request[any], []byte], requestTransformer dto.RequestTransformer[any]) common.Function[[]byte, []byte] {
	return &Dispatcher{handlers, requestTransformer}
}
func (d *Dispatcher) Apply(in []byte) ([]byte, error) {
	req, err := d.requestTransformer.BytesToRequest(in)
	if err != nil {
		return nil, err
	}
	if handler, ok := d.handlers[req.Type]; ok {
		return handler.Apply(req)
	}
	return nil, errors.New("no handler")
}
