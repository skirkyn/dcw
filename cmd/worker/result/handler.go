package result

import (
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/worker/client"
)

type Handler[In any] struct {
	client             client.Client
	requestTransformer common.RequestTransformer[In]
}

func NewHandler[In any](client client.Client, requestTransformer common.RequestTransformer[In]) common.Consumer[common.Request[In]] {
	return &Handler[In]{client, requestTransformer}
}

func (h *Handler[In]) Consume(res common.Request[In]) error {

	req, err := h.requestTransformer.RequestToBytes(res)
	if err != nil {
		return err
	}
	_, err = h.client.Call(req)
	return err

}
