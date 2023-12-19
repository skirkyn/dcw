package zmq

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/controller/generators"
	"github.com/skirkyn/dcw/cmd/controller/generators/gerrorrs"
	"github.com/skirkyn/dcw/cmd/dto"
	"github.com/skirkyn/dcw/cmd/dto/bf"
)

type StringGeneratorHandler struct {
	Generator *generators.Generator[string]
}

func NewGeneratorHandler(generator *generators.Generator[string]) StringGeneratorHandler {
	return StringGeneratorHandler{generator}
}

func (gh *StringGeneratorHandler) Handle(req dto.Request[int], respChannel *chan dto.Response[[]string]) {
	if req.Body() > 0 {
		res, err := (*gh.Generator).Next(req.Body())
		done, errStr := gh.extractErrorTextAndDone(err)
		*respChannel <- bf.Response[[]string]{Data: res, Done: done, Err: errStr}
	}
}

func (gh *StringGeneratorHandler) extractErrorTextAndDone(err error) (bool, string) {

	if err != nil {
		return errors.Is(err, gerrorrs.PotentialResultsExhaustedError), err.Error()
	}
	return false, ""
}
