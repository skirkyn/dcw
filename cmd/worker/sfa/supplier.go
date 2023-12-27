package sfa

import (
	"context"
	"github.com/unknownfeature/dcw/cmd/common"
	"github.com/unknownfeature/dcw/cmd/common/dto"
	"golang.org/x/sync/semaphore"
)

type WorkSupplier struct {
	batchSize          int
	requestTransformer dto.RequestTransformer[int]
	semaphore          *semaphore.Weighted
	context            context.Context
}

func NewSupplier(batchSize int, requestTransformer dto.RequestTransformer[int], semaphore *semaphore.Weighted, context context.Context) common.Supplier[[]byte] {

	return &WorkSupplier{batchSize, requestTransformer, semaphore, context}
}

func (s *WorkSupplier) Supply() ([]byte, error) {
	err := s.semaphore.Acquire(s.context, 1)
	if err != nil {
		return nil, err
	}
	req := dto.Request[int]{Type: dto.Work, Body: s.batchSize}
	bytes, err := s.requestTransformer.RequestToBytes(req)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
