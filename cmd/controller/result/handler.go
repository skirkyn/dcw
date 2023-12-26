package result

import (
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/common/dto"
	"log"
)

type Handler[Result any] struct {
	responseTransformer dto.ResponseTransformer[string]
}

func NewHandler[Result any](responseTransformer dto.ResponseTransformer[string]) common.Function[dto.Request[any], []byte] {
	return &Handler[Result]{responseTransformer}
}

func (h *Handler[Result]) Apply(req dto.Request[any]) ([]byte, error) {

	log.Println("+++++++ !!!!!!!! found result !!!!!! +++++++++")

	log.Println(req.Body)

	// maybe email the result

	return h.responseTransformer.ResponseToBytes(dto.Response[string]{Done: true, Body: ""})
}
