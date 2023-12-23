package svf

import (
	"context"
	"github.com/skirkyn/dcw/cmd/util/bytz"
	"github.com/skirkyn/dcw/cmd/worker/work"
	"golang.org/x/sync/semaphore"
)

type WorkSupplier struct {
	supply    []byte
	semaphore *semaphore.Weighted
	context   context.Context
}

func NewSupplier(batchSize int, semaphore *semaphore.Weighted, context context.Context) work.Supplier {
	return &WorkSupplier{bytz.IntToByteSlice(batchSize), semaphore, context}
}

func (s *WorkSupplier) Supply() ([]byte, error) {
	err := s.semaphore.Acquire(s.context, 1)
	defer s.semaphore.Release(1)
	if err != nil {
		return nil, err
	}
	return s.supply, nil
}
