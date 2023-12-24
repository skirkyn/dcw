package hr

import (
	"errors"
	"fmt"
	"github.com/skirkyn/dcw/cmd/common"
)

type FormattingBodySupplier[In any] struct {
	format string
}

func NewFormattingBodySupplier[In any](format string) common.Function[In, []byte] {
	return &FormattingBodySupplier[In]{format: format}
}

func (sf *FormattingBodySupplier[In]) Apply(in In) ([]byte, error) {
	if in == nil {
		return nil, errors.New("input can't be nil")
	}
	return []byte(fmt.Sprintf(sf.format, in)), nil
}
