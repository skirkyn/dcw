package result

import (
	"github.com/unknownfeature/dcw/cmd/common"
	"github.com/unknownfeature/dcw/cmd/common/dto"
	"log"
)

type Handler[Result any] struct {
	responseTransformer dto.ResponseTransformer[string]
	resultConsumer      common.Consumer[bool]
}

func NewHandler[Result any](responseTransformer dto.ResponseTransformer[string], resultConsumer common.Consumer[bool]) common.Function[dto.Request[any], []byte] {
	return &Handler[Result]{responseTransformer, resultConsumer}
}

func (h *Handler[Result]) Apply(req dto.Request[any]) ([]byte, error) {

	log.Println("+++++++ !!!!!!!! found result !!!!!! +++++++++ \n", req.Body)

	err := h.resultConsumer.Consume(true)
	if err != nil {
		log.Println("error consuming result ", err)
	}

	// maybe email the result

	return h.responseTransformer.ResponseToBytes(dto.Response[string]{Done: true, Body: ""})
}
