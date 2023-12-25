package server

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/common"
)

type Dispatcher struct {
	handlers           map[common.Type]common.Function[common.Request[any], []byte]
	requestTransformer common.RequestTransformer[any]
}

func NewDispatcher(handlers map[common.Type]common.Function[common.Request[any], []byte], requestTransformer common.RequestTransformer[any]) common.Function[[]byte, []byte] {
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
