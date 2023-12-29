package sfa

import (
	"errors"
	"github.com/unknownfeature/dcw/cmd/common"
	"github.com/unknownfeature/dcw/cmd/common/dto"
	"log"
	"sync/atomic"
)

type StringGeneratorHandler struct {
	workSupplier        common.Function[int, []string]
	responseTransformer dto.ResponseTransformer[[]string]
	resultFound         *atomic.Bool
}

// todo refactor this
func NewGeneratorHandler(supplier common.Function[int, []string],
	responseTransformer dto.ResponseTransformer[[]string]) *StringGeneratorHandler {
	return &StringGeneratorHandler{supplier, responseTransformer, &atomic.Bool{}}
}

func (gh *StringGeneratorHandler) Apply(req dto.Request[any]) ([]byte, error) {

	if gh.resultFound.Load() {
		resp := dto.Response[[]string]{Done: true, Body: []string{}}
		bytes, _ := gh.responseTransformer.ResponseToBytes(resp)
		return bytes, PotentialResultsExhaustedError
	}
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

func (gh *StringGeneratorHandler) Consume(resFound bool) error {

	gh.resultFound.Store(resFound)
	return nil

}
