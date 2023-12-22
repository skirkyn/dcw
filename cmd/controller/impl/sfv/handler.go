package sfv

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/controller/handler"
	"github.com/skirkyn/dcw/cmd/controller/supplier"
	"github.com/skirkyn/dcw/cmd/dto"
	"github.com/skirkyn/dcw/cmd/dto/impl/sfv"
	"log"
)

type StringGeneratorHandler struct {
	supplier            supplier.Supplier[int, []string]
	requestTransformer  func([]byte) (dto.Request[int], error)
	responseTransformer func(dto.Response[[]string]) ([]byte, error)
}

func NewGeneratorHandler(supplier supplier.Supplier[int, []string],
	requestTransformer func([]byte) (dto.Request[int], error),
	responseTransformer func(dto.Response[[]string]) ([]byte, error)) handler.Handler {
	return &StringGeneratorHandler{supplier, requestTransformer, responseTransformer}
}

func (gh *StringGeneratorHandler) Handle(reqRaw []byte,
	respChannel chan []byte, errChannel chan error) {
	req, err := gh.requestTransformer(reqRaw)

	if err != nil && gh.handleError(err.Error(), respChannel, errChannel) != nil {
		return
	}

	result, err := gh.supplier.Next(req.Body())
	//return gp.toResponse(res, sfv)
	resp := sfv.Response{Data: result, Err: extractErrorTextAndDone(err)}

	bytes, err := gh.responseTransformer(resp)
	if err != nil && gh.handleError(err.Error(), respChannel, errChannel) != nil {
		return
	}

	respChannel <- bytes
}

func (gh *StringGeneratorHandler) handleError(message string, respChannel chan []byte, errChannel chan error) error {
	log.Printf("cant transform the request %s", message)
	//todo limit retries
	resp, err := gh.responseTransformer(sfv.Response{Err: sfv.NewError("error converting request", true)})

	if err != nil {
		errChannel <- err
	} else {
		respChannel <- resp
	}

	return err
}

func extractErrorTextAndDone(err error) dto.Error {

	if err != nil {
		return sfv.NewError(err.Error(), !errors.Is(err, PotentialResultsExhaustedError))
	}
	return nil
}
