package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/skirkyn/dcw/cmd/common/config"
	"github.com/skirkyn/dcw/cmd/common/dto"
	"github.com/skirkyn/dcw/cmd/worker/client/zmq"
	"github.com/skirkyn/dcw/cmd/worker/result"
	"github.com/skirkyn/dcw/cmd/worker/runner"
	"github.com/skirkyn/dcw/cmd/worker/sfa"
	"github.com/skirkyn/dcw/cmd/worker/verifier/hr"
	"github.com/skirkyn/dcw/cmd/worker/verifier/hr/cb"
	"github.com/skirkyn/dcw/cmd/worker/verifier/hr/cb/vi"
	"golang.org/x/sync/semaphore"
	"log"
	"net/http"
	"runtime"
	"time"
)

func Run() {
	workerConfig, err := config.ReadWorkerConfig()
	if err != nil {
		log.Fatal("can't read worker config", err)
	}

	commonConfig, err := config.ReadCommonConfig()

	if err != nil {
		log.Fatal("can't read common config", err)
	}

	// todo it should be a factory based on config but it is what it is
	rnr := getCbRunner(commonConfig, workerConfig)
	wg := rnr.Start()
	wg.Wait()
}
func getCbRunner(commonConfig config.CommonConfig, workerConfig config.WorkerConfig[config.HttpRequestVerifier[config.CbCustomConfig], config.CbCustomConfig]) runner.Runner {

	clientConfig := zmq.Config{

		Host:                                 commonConfig.ControllerHost,
		Port:                                 commonConfig.ControllerPort,
		Id:                                   uuid.NewString(),
		MaxConnectAttempts:                   workerConfig.ConnAttempts,
		TimeToSleepBetweenConnectionAttempts: time.Duration(workerConfig.ConnAttemptsTtsSec),
	}

	// todo authorization
	client, err := zmq.NewZMQClient(clientConfig)

	if err != nil {
		log.Fatal("can't create workClient", err)
	}

	runnerConfig := runner.Config{WorkersCount: runtime.NumCPU() * workerConfig.WorkersFactor}
	sem := semaphore.NewWeighted(int64(runnerConfig.WorkersCount * workerConfig.WorkersSemaphoreWeight))
	ctx := context.Background()
	workReqTrans := dto.NewRequestTransformer[int]()
	workRespTrans := dto.NewResponseTransformer[[]string]()

	resultRequestTrans := dto.NewRequestTransformer[string]()
	headersSupplier := hr.NewSimpleHeadersSupplier(workerConfig.VerifierConfig.Headers)
	method := workerConfig.VerifierConfig.Method
	verifyRequestSupplier := hr.NewRequestSupplier[string](method, headersSupplier, workerConfig.VerifierConfig.Method, cb.NewBodySupplier(workerConfig.VerifierConfig.Body))
	proofTokenRequestSupplier := hr.NewRequestSupplier[map[string]string](method, headersSupplier, workerConfig.VerifierConfig.CustomConfig.VerifyIdentityUrl, vi.NewBodySupplier(workerConfig.VerifierConfig.CustomConfig.VerifyIdentityBody))

	successPredicate := cb.NewResponseHandler(http.DefaultClient, proofTokenRequestSupplier)
	verifier := hr.NewVerifier[string](http.DefaultClient, verifyRequestSupplier, successPredicate)
	worker := sfa.NewWorker(sem, ctx, workRespTrans, verifier)
	workReqSupplier := sfa.NewSupplier(workerConfig.BatchSize, workReqTrans, sem, ctx)
	resHandler := result.NewHandler[string](client, resultRequestTrans)

	return runner.NewDefaultRunner[string](runnerConfig, client, worker, workReqSupplier, resHandler)
}
