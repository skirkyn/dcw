package worker

import (
	"github.com/skirkyn/dcw/cmd/worker/client"
	"github.com/skirkyn/dcw/cmd/worker/worker"
	"log"
)

func runWorker[Out any, Res any](timesToRetry int, requestSupplier chan []byte, client client.Client, worker worker.Worker[Out], stop chan bool, res chan Res) error {

	for {
		select {
		case shouldStop := <-stop:
			if shouldStop {
				return nil
			}
		default:
			req := <-requestSupplier
			for i := timesToRetry; i > 0; i-- {
				resp, err := client.Call(req)
				if err != nil {
					log.Printf("error calling the server %s", err.Error())

				} else {
					go processBatch(resp, worker, res)
				}

			}
		}
	}
}

func processBatch[Out any, Res any](client client.Client, batch []byte, worker worker.Worker[Out], res chan Res) {
	result, err := worker.Process(batch)
	if err != nil {

	}
}
