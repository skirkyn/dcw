package pt

import (
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/worker/verifier/hr"
)

var defaultRequestTemplate = `
{
  "recaptcha_token": "",
  "proof_token": "%s"
}`

type BodySupplier struct {
	formatter common.Function[string, []byte]
}

func NewBodySupplier() common.Function[map[string]string, []byte] {
	return &BodySupplier{hr.NewFormattingBodySupplier[string](defaultRequestTemplate)}
}

func (sf *BodySupplier) Apply(in map[string]string) ([]byte, error) {

	return sf.formatter.Apply(in["proof_token"])
}
