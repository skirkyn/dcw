package sfv

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/common"
)

type StringGeneratorHandler struct {
	workSupplier        common.Function[int, []string]
	requestTransformer  common.RequestTransformer[int]
	responseTransformer common.ResponseTransformer[[]string]
}

func NewGeneratorHandler(supplier common.Function[int, []string],
	requestTransformer common.RequestTransformer[int],
	responseTransformer common.ResponseTransformer[[]string]) common.Function[[]byte, []byte] {
	return &StringGeneratorHandler{supplier, requestTransformer, responseTransformer}
}

func (gh *StringGeneratorHandler) Apply(reqRaw []byte) ([]byte, error) {
	req, err := gh.requestTransformer.BytesToRequest(reqRaw)

	if err != nil {
		return nil, err
	}

	result, err := gh.workSupplier.Apply(req.Body)
	resp := common.Response[[]string]{Done: !errors.Is(err, PotentialResultsExhaustedError), Body: result}
	bytes, err := gh.responseTransformer.ResponseToBytes(resp)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
