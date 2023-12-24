package result

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/common"
	"log"
)

type Handler[Result any] struct {
	requestTransformer common.RequestTransformer[string]
}

func NewHandler[Result any](requestTransformer common.RequestTransformer[string]) common.Function[[]byte, []byte] {
	return &Handler[Result]{requestTransformer}
}

func (h *Handler[Result]) Apply(in []byte) ([]byte, error) {
	if in == nil {
		log.Println("can't handle result nil")
		return nil, errors.New("nil result")
	}
	req, err := h.requestTransformer.BytesToRequest(in)

	if err != nil {
		log.Printf("error transforming request %s", string(in))
		return nil, errors.New("nil result")
	}
	log.Println("+++++++ !!!!!!!! found result !!!!!! +++++++++")

	log.Println(req.Body)

	return []byte{}, err
}
