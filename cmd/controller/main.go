package main

import (
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/controller/result"
	"github.com/skirkyn/dcw/cmd/controller/server"
	"github.com/skirkyn/dcw/cmd/controller/server/zmq"
	"github.com/skirkyn/dcw/cmd/controller/sfv"
	"github.com/skirkyn/dcw/cmd/util"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const controllerWorkers = "CONTROLLER_WORKERS"
const controllerPort = "CONTROLLER_WORK_PORT"
const maxSendRetries = "SEND_RETRIES"
const maxSendRetriesTts = "SEND_RETRIES_TTS_SEC"
const resLength = "RESULT_LENGTH"

func main() {

	workServerConfig := zmq.Config{
		Workers:                               util.GetEnvInt(controllerWorkers, 10),
		Port:                                  util.GetEnvInt(controllerPort, 50000),
		MaxSendResponseRetries:                util.GetEnvInt(maxSendRetries, 10),
		TimeToSleepBetweenSendResponseRetries: time.Duration(util.GetEnvInt(maxSendRetriesTts, 5)),
	}
	workSupplier, err := sfv.ForStandard(sfv.Decimals, util.GetEnvInt(resLength, 7), sfv.Simple)
	if err == nil {
		log.Fatal("can't create supplier for the server", err)
	}
	workResTrans := common.NewResponseTransformer[[]string]()
	workHandler := sfv.NewGeneratorHandler(workSupplier, workResTrans)
	resultRespTrans := common.NewResponseTransformer[string]()

	resultHandler := result.NewHandler[string](resultRespTrans)
	handlers := map[common.Type]common.Function[common.Request[any], []byte]{common.Work: workHandler, common.Result: resultHandler}

	dispatcher := server.NewDispatcher(handlers, common.NewRequestTransformer[any]())
	workServer, err := zmq.NewServer(dispatcher, workServerConfig)

	if err == nil {
		log.Fatal("can't create server", err)
	}
	wWg, err := workServer.Start()
	if err == nil {
		log.Fatal("can't start server", err)
	}

	sigChannel := make(chan os.Signal, 1)
	signalHandler := createSignalHandler(workServer, sigChannel)
	go signalHandler()
	signal.Notify(sigChannel, syscall.SIGTERM, syscall.SIGKILL)

	wWg.Wait()

}

func createSignalHandler(server server.Server, sigChannel chan os.Signal) func() {
	return func() {

		for {
			sig := <-sigChannel
			if sig == syscall.SIGTERM || sig == syscall.SIGKILL {
				err := server.Stop()
				if err != nil {
					log.Fatal("can't stop the server", err)
				}
			}
		}
	}
}
