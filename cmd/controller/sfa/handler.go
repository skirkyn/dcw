package sfa

import (
	"errors"
	"github.com/unknownfeature/dcw/cmd/common"
	"github.com/unknownfeature/dcw/cmd/common/dto"
	"log"
)

type StringGeneratorHandler struct {
	workSupplier        common.Function[int, []string]
	responseTransformer dto.ResponseTransformer[[]string]
}

func NewGeneratorHandler(supplier common.Function[int, []string],
	responseTransformer dto.ResponseTransformer[[]string]) common.Function[dto.Request[any], []byte] {
	return &StringGeneratorHandler{supplier, responseTransformer}
}

func (gh *StringGeneratorHandler) Apply(req dto.Request[any]) ([]byte, error) {

	result, err := gh.workSupplier.Apply(int(req.Body.(float64)))
	resp := dto.Response[[]string]{Done: errors.Is(err, PotentialResultsExhaustedError), Body: result}
	bytes, err := gh.responseTransformer.ResponseToBytes(resp)
	if err != nil {
		return nil, err
	}
	if len(bytes) != 0 {
		log.Printf("sending another batch with %s as last", string(result[len(result)-1]))
	}
	return bytes, nil
}
