package sfv

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/controller/handler"
	"github.com/skirkyn/dcw/cmd/controller/work"
	"github.com/skirkyn/dcw/cmd/dto"
	"github.com/skirkyn/dcw/cmd/dto/sfv"
	"log"
)

type StringGeneratorHandler struct {
	supplier            work.Supplier[int, []string]
	requestTransformer  func([]byte) (dto.Request[int], error)
	responseTransformer func(dto.Response[[]string]) ([]byte, error)
}

func NewGeneratorHandler(supplier work.Supplier[int, []string],
	requestTransformer func([]byte) (dto.Request[int], error),
	responseTransformer func(dto.Response[[]string]) ([]byte, error)) handler.Handler {
	return &StringGeneratorHandler{supplier, requestTransformer, responseTransformer}
}

func (gh *StringGeneratorHandler) Handle(reqRaw []byte) []byte {
	req, err := gh.requestTransformer(reqRaw)

	if err != nil {
		return gh.handleError(err.Error())
	}

	result, err := gh.supplier.Supply(req.Body())
	//return gp.toResponse(res, sfv)
	resp := sfv.Response[[]string]{Data: result, Err: extractErrorTextAndDone(err)}

	bytes, err := gh.responseTransformer(resp)
	if err != nil {
		return gh.handleError(err.Error())
	}

	return bytes
}

func (gh *StringGeneratorHandler) handleError(message string) []byte {
	log.Printf("cant transform the request %s", message)
	//todo limit retries
	resp, err := gh.responseTransformer(sfv.Response[[]string]{Err: sfv.NewError("error converting request", true)})
	if err != nil {
		// we cant even transform it so something bizarre is going on
		return []byte(err.Error())
	}

	return resp
}

func extractErrorTextAndDone(err error) dto.Error {

	if err != nil {
		return sfv.NewError(err.Error(), !errors.Is(err, PotentialResultsExhaustedError))
	}
	return nil
}
