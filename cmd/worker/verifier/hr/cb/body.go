package cb

import (
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/worker/verifier/hr"
)

type BodySupplier struct {
	formatter common.Function[string, []byte]
}

func NewBodySupplier(requestTemplate string) common.Function[string, []byte] {
	return &BodySupplier{hr.NewFormattingBodySupplier[string](requestTemplate)}
}

func (sf *BodySupplier) Apply(in string) ([]byte, error) {
	return sf.formatter.Apply(in)
}
