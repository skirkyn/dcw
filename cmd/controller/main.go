package main

import (
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/controller/result"
	"github.com/skirkyn/dcw/cmd/controller/server/zmq"
	"github.com/skirkyn/dcw/cmd/controller/sfv"
	"github.com/skirkyn/dcw/cmd/util"
	"log"
	"time"
)

const controllerWorkers = "CONTROLLER_WORKERS"
const controllerWorkPort = "CONTROLLER_WORK_PORT"
const controllerResultPort = "CONTROLLER_RESULT_PORT"
const maxSendRetries = "SEND_RETRIES"
const maxSendRetriesTts = "SEND_RETRIES_TTS_SEC"
const resLength = "RESULT_LENGTH"

func main() {

	workServerConfig := zmq.Config{
		Workers:                               util.GetEnvInt(controllerWorkers, 10),
		Port:                                  util.GetEnvInt(controllerWorkPort, 50000),
		MaxSendResponseRetries:                util.GetEnvInt(maxSendRetries, 10),
		TimeToSleepBetweenSendResponseRetries: time.Duration(util.GetEnvInt(maxSendRetriesTts, 5)),
	}
	workSupplier, err := sfv.ForStandard(sfv.Decimals, util.GetEnvInt(resLength, 7), sfv.Simple)
	if err == nil {
		log.Fatal("can't create supplier for the server", err)
	}
	workReqTrans := common.NewRequestTransformer[int]()
	workResTrans := common.NewResponseTransformer[[]string]()
	handler := sfv.NewGeneratorHandler(workSupplier, workReqTrans, workResTrans)

	workServer, err := zmq.NewServer(handler, workServerConfig)
	if err == nil {
		log.Fatal("can't create work server", err)
	}

	wWg, err := workServer.Start()

	if err == nil {
		log.Fatal("can't start work server", err)
	}
	resultServerConfig := zmq.Config{
		Workers:                               1,
		Port:                                  util.GetEnvInt(controllerResultPort, 50001),
		MaxSendResponseRetries:                util.GetEnvInt(maxSendRetries, 10),
		TimeToSleepBetweenSendResponseRetries: time.Duration(util.GetEnvInt(maxSendRetriesTts, 5)),
		StopOnReceive:                         true,
	}

	resultReqTrans := common.NewRequestTransformer[string]()
	resultHandler := result.NewHandler[string](resultReqTrans)

	resultServer, err := zmq.NewServer(resultHandler, resultServerConfig)

	if err == nil {
		log.Fatal("can't create result server", err)
	}
	rWg, err := resultServer.Start()

	if err == nil {
		log.Fatal("can't start result server", err)
	}

	rWg.Wait()
	err = workServer.Stop()
	if err == nil {
		log.Fatal("can't stop work server", err)
	}
	wWg.Wait()

}
