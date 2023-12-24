package sfv

import (
	"context"
	"github.com/skirkyn/dcw/cmd/common"
	"golang.org/x/sync/semaphore"
)

type WorkSupplier struct {
	batchSize          int
	requestTransformer common.RequestTransformer[int]
	semaphore          *semaphore.Weighted
	context            context.Context
}

func NewSupplier(batchSize int, requestTransformer common.RequestTransformer[int], semaphore *semaphore.Weighted, context context.Context) common.Supplier[[]byte] {

	return &WorkSupplier{batchSize, requestTransformer, semaphore, context}
}

func (s *WorkSupplier) Supply() ([]byte, error) {
	err := s.semaphore.Acquire(s.context, 1)
	defer s.semaphore.Release(1)
	if err != nil {
		return nil, err
	}
	req := common.Request[int]{Type: common.Work, Body: s.batchSize}
	bytes, err := s.requestTransformer.RequestToBytes(req)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
