package sfv

import (
	"github.com/skirkyn/dcw/cmd/controller/server"
	"github.com/skirkyn/dcw/cmd/dto"
	"log"
)

type StringGeneratorHandler struct {
	processor           server.Processor[int, []string]
	requestTransformer  dto.RequestTransformer[int]
	responseTransformer dto.ResponseTransformer[[]string]
}

func NewGeneratorHandler(processor server.Processor[int, []string],
	requestTransformer dto.RequestTransformer[int],
	responseTransformer dto.ResponseTransformer[[]string]) server.Handler {
	return &StringGeneratorHandler{processor, requestTransformer, responseTransformer}
}

func (gh *StringGeneratorHandler) Handle(reqRaw []byte,
	respChannel *chan []byte, errChannel *chan error) {
	req, err := gh.requestTransformer.Transform(reqRaw)

	if err != nil && gh.handleError(err.Error(), respChannel, errChannel) != nil {
		return
	}

	result := gh.processor.Process(req)
	resp, err := gh.responseTransformer.Transform(&result)

	if err != nil && gh.handleError(err.Error(), respChannel, errChannel) != nil {
		return
	}

	*respChannel <- resp
}

func (gh *StringGeneratorHandler) handleError(message string, respChannel *chan []byte, errChannel *chan error) error {
	log.Printf("cant transform the request %s", message)
	//todo limit retries
	errResp := dto.NewErrorResponse[[]string](NewError("error converting request", true))
	resp, err := gh.responseTransformer.Transform(&errResp)

	if err != nil {
		*errChannel <- err
	} else {
		*respChannel <- resp
	}

	return err
}
