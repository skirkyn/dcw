package hr

import (
	"fmt"
	"github.com/unknownfeature/dcw/cmd/common"
)

type FormattingBodySupplier[In any] struct {
	format string
}

func NewFormattingBodySupplier[In any](format string) common.Function[In, []byte] {
	return &FormattingBodySupplier[In]{format: format}
}

func (sf *FormattingBodySupplier[In]) Apply(in In) ([]byte, error) {
	return []byte(fmt.Sprintf(sf.format, in)), nil
}
