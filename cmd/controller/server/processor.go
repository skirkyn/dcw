package server

import (
	"github.com/skirkyn/dcw/cmd/controller/processor"
	"github.com/skirkyn/dcw/cmd/controller/work"
	"github.com/skirkyn/dcw/cmd/dto"
)

type Processor[In any, Out any] struct {
	supplier   work.Supplier[In, Out]
	toResponse func(Out, error) dto.Response[Out]
}

func NewProcessor[In any, Out any](supplier work.Supplier[In, Out]) processor.Processor[In, Out] {
	return &Processor[In, Out]{supplier: supplier}
}

func (gp *Processor[In, Out]) Process(request dto.Request[In]) dto.Response[Out] {
	res, err := gp.supplier.Supply(request.Body())
	return gp.toResponse(res, err)
}
