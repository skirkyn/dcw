package runner

import (
	"github.com/unknownfeature/dcw/cmd/common"
	"github.com/unknownfeature/dcw/cmd/common/dto"
	"github.com/unknownfeature/dcw/cmd/worker/client"
	"log"
	"sync"
	"sync/atomic"
)

type Runner interface {
	Start() *sync.WaitGroup
	Stop()
}

type Config struct {
	WorkersCount int
}
type DefaultRunner[Result any] struct {
	config          Config
	client          client.Client
	worker          common.Function[[]byte, *dto.Request[Result]]
	requestSupplier common.Supplier[[]byte]
	resultHandler   common.Consumer[dto.Request[Result]]
	stop            *atomic.Bool
}

func NewDefaultRunner[Result any](config Config, client client.Client, worker common.Function[[]byte, *dto.Request[Result]], requestSupplier common.Supplier[[]byte], resultHandler common.Consumer[dto.Request[Result]]) Runner {

	return &DefaultRunner[Result]{config, client, worker, requestSupplier, resultHandler, &atomic.Bool{}}
}

func (r *DefaultRunner[Result]) Start() *sync.WaitGroup {

	wg := sync.WaitGroup{}
	wg.Add(r.config.WorkersCount)

	for i := r.config.WorkersCount; i > 0; i-- {
		go r.runWorker(&wg)
	}
	return &wg
}

func (r *DefaultRunner[Result]) Stop() {
	r.stop.Store(true)
}

func (r *DefaultRunner[Result]) runWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if r.stop.Load() {
			return
		}
		log.Println("calling supply")
		supply, err := r.requestSupplier.Supply()
		if err != nil {
			log.Printf("error calling supply, will exit %s", err.Error())
		}
		r.doWork(supply)
	}
}

func (r *DefaultRunner[Result]) doWork(req []byte) {
	log.Println("calling server")
	resp, err := r.client.Call(req)
	if err != nil && err.Error() == "potential results exhausted" {
		log.Printf("error calling the server %s", err.Error())
		r.Stop()

		return
	}
	log.Println("applying result")
	res, err := r.worker.Apply(resp)
	if err != nil {
		log.Printf("error processing work from server %s", err.Error())
		return
	}
	if res == nil {
		return
	}

	// todo add retries
	err = r.resultHandler.Consume(*res)
	if err != nil {
		log.Printf("error handling result %s", err.Error())
	} else {
		r.Stop()
	}

}
