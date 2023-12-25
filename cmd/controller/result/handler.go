package result

import (
	"github.com/skirkyn/dcw/cmd/common"
	"log"
)

type Handler[Result any] struct {
	responseTransformer common.ResponseTransformer[string]
}

func NewHandler[Result any](responseTransformer common.ResponseTransformer[string]) common.Function[common.Request[any], []byte] {
	return &Handler[Result]{responseTransformer}
}

func (h *Handler[Result]) Apply(req common.Request[any]) ([]byte, error) {

	log.Println("+++++++ !!!!!!!! found result !!!!!! +++++++++")

	log.Println(req.Body)

	// maybe email the result

	return h.responseTransformer.ResponseToBytes(common.Response[string]{Done: true, Body: ""})
}
