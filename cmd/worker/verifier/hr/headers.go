package hr

import (
	"github.com/unknownfeature/dcw/cmd/common"
)

type SimpleHeadersSupplier struct {
	defaultHeaders map[string]string
}

func NewSimpleHeadersSupplier(defaultHeaders map[string]string) common.Function[any, map[string]string] {
	return &SimpleHeadersSupplier{defaultHeaders: defaultHeaders}
}

func (sf *SimpleHeadersSupplier) Apply(in any) (map[string]string, error) {
	return sf.defaultHeaders, nil
}
