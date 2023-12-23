package worker

import (
	"context"
	"github.com/skirkyn/dcw/cmd/dto"
	"github.com/skirkyn/dcw/cmd/worker/verifier"
	"golang.org/x/sync/semaphore"
	"log"
)

type Worker struct {
	semaphore   *semaphore.Weighted
	context     context.Context
	transformer dto.ResponseTransformer[[]string]
	verifier    verifier.Verifier[string]
}

func (w *Worker) Process(work []byte) []byte {
	resp, err := w.transformer.BytesToResponse(work)
	if err != nil {
		log.Printf("can't process response %s because of %s", string(work), err.Error())
	}
	if resp.Error() != nil {
		if !resp.Error().Retry() {
			return make([]byte, 0)
		}
	}

	input := resp.Body()

	for i := 0; i < len(input); i++ {
		if w.verifier.Verify(input[i]) {
			return []byte(input[i])
		}
	}

	return nil
}
