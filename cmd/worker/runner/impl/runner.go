package impl

import (
	"github.com/skirkyn/dcw/cmd/worker/client"
	"github.com/skirkyn/dcw/cmd/worker/result"
	"github.com/skirkyn/dcw/cmd/worker/work"
	"github.com/skirkyn/dcw/cmd/worker/worker"
	"log"
	"sync/atomic"
)

type Config struct {
	WorkersCount int
}
type Runner[Result any] struct {
	config          Config
	client          client.Client
	worker          worker.Worker[Result]
	requestSupplier work.Supplier
	resultHandler   result.Handler[Result]
	stop            *atomic.Bool
}

func NewRunner[Result any](config Config, client client.Client, worker worker.Worker[Result], requestSupplier work.Supplier, resultHandler result.Handler[Result]) Runner[Result] {

	return Runner[Result]{config, client, worker, requestSupplier, resultHandler, &atomic.Bool{}}
}

func (r *Runner[Result]) Start() {

	for i := r.config.WorkersCount; i > 0; i-- {
		go r.runWorker()
	}

}

func (r *Runner[Result]) Stop() {
	r.stop.Store(true)
}

func (r *Runner[Result]) runWorker() {

	for {
		if r.stop.Load() {
			return
		}
		supply, err := r.requestSupplier.Supply()
		if err != nil {
			log.Printf("error calling supply, will exit %s", err.Error())
		}
		go r.doWork(supply)
	}
}

func (r *Runner[Result]) doWork(req []byte) {

	resp, err := r.client.Call(req)
	if err != nil {
		log.Printf("error calling the server %s", err.Error())
	}

	result := r.worker.Process(resp)
	if result != nil {
		err = r.resultHandler.Handle(result)
		if err != nil {
			log.Printf("error handling result %s", err.Error())

		}
	}

}
