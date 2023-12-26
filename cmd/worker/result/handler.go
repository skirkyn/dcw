package result

import (
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/common/dto"
	"github.com/skirkyn/dcw/cmd/worker/client"
)

type Handler[In any] struct {
	client             client.Client
	requestTransformer dto.RequestTransformer[In]
}

func NewHandler[In any](client client.Client, requestTransformer dto.RequestTransformer[In]) common.Consumer[dto.Request[In]] {
	return &Handler[In]{client, requestTransformer}
}

func (h *Handler[In]) Consume(res dto.Request[In]) error {

	req, err := h.requestTransformer.RequestToBytes(res)
	if err != nil {
		return err
	}
	_, err = h.client.Call(req)
	return err

}
