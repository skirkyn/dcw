package impl

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/dto"
	"github.com/skirkyn/dcw/cmd/worker/client"
	"github.com/skirkyn/dcw/cmd/worker/result"
	"github.com/skirkyn/dcw/cmd/worker/runner"
)

type Handler[In any] struct {
	client             client.Client
	requestTransformer dto.RequestTransformer[In]
	runner             runner.Runner
}

func NewHandler[In any](client client.Client, requestTransformer dto.RequestTransformer[In], runner runner.Runner) result.Handler[In] {
	return &Handler[In]{client, requestTransformer, runner}
}

func (h *Handler[In]) Handle(res dto.Request[In]) error {
	if res == nil {
		return errors.New("result is null")
	}
	req, err := h.requestTransformer.RequestToBytes(res)
	if err != nil {
		return err
	}
	_, err = h.client.Call(req)
	return err

}
