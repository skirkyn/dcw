package sfv

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/controller/generators"
	"github.com/skirkyn/dcw/cmd/controller/generators/gerr"
	"github.com/skirkyn/dcw/cmd/dto"
)

type Processor[In any, Out any] struct {
	generator   generators.Generator[In, Out]
	emptyResult Out
}

func (gp *Processor[In, Out]) Process(val In) (Out, dto.Error) {
	res, err := gp.generator.Next(val)
	if err != nil {
		return gp.emptyResult, gp.extractErrorTextAndDone(err)
	}
	return res, nil
}

func (gp *Processor[In, Out]) extractErrorTextAndDone(err error) dto.Error {

	if err != nil {
		return NewError(err.Error(), !errors.Is(err, gerr.PotentialResultsExhaustedError))
	}
	return nil
}
