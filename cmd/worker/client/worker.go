package client

import (
	"github.com/skirkyn/dcw/cmd/dto"
)

type Worker struct {
	requestTransformer  func(dto.Request[int]) ([]byte, error)
	responseTransformer func([]byte) (dto.Response[[]string], error)
	client              Client
}
