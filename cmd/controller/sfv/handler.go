package sfv

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/common"
)

type StringGeneratorHandler struct {
	workSupplier        common.Function[int, []string]
	responseTransformer common.ResponseTransformer[[]string]
}

func NewGeneratorHandler(supplier common.Function[int, []string],
	responseTransformer common.ResponseTransformer[[]string]) common.Function[common.Request[any], []byte] {
	return &StringGeneratorHandler{supplier, responseTransformer}
}

func (gh *StringGeneratorHandler) Apply(req common.Request[any]) ([]byte, error) {

	result, err := gh.workSupplier.Apply(req.Body.(int))
	resp := common.Response[[]string]{Done: !errors.Is(err, PotentialResultsExhaustedError), Body: result}
	bytes, err := gh.responseTransformer.ResponseToBytes(resp)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
