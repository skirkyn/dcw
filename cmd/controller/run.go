package main

import (
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/common/config"
	"github.com/skirkyn/dcw/cmd/common/dto"
	"github.com/skirkyn/dcw/cmd/controller/result"
	"github.com/skirkyn/dcw/cmd/controller/server"
	"github.com/skirkyn/dcw/cmd/controller/server/zmq"
	"github.com/skirkyn/dcw/cmd/controller/sfv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	controllerConfig, err := config.ReadControllerConfig()
	if err != nil {
		log.Fatal("can't read worker config", err)
	}

	commonConfig, err := config.ReadCommonConfig()

	if err != nil {
		log.Fatal("can't read common config", err)
	}

	workServerConfig := zmq.Config{
		Workers:                               controllerConfig.Workers,
		Port:                                  commonConfig.ControllerPort,
		MaxSendResponseRetries:                controllerConfig.MaxSendRetries,
		TimeToSleepBetweenSendResponseRetries: time.Duration(controllerConfig.MaxSendRetriesTtsSec),
	}
	// todo this also has to be extracted into a factory
	workSupplier, err := sfv.ForStandard(sfv.Decimals, controllerConfig.ResLength, sfv.Simple)
	if err != nil {
		log.Fatal("can't create supplier for the server", err)
	}
	workResTrans := dto.NewResponseTransformer[[]string]()
	workHandler := sfv.NewGeneratorHandler(workSupplier, workResTrans)
	resultRespTrans := dto.NewResponseTransformer[string]()

	resultHandler := result.NewHandler[string](resultRespTrans)
	handlers := map[dto.Type]common.Function[dto.Request[any], []byte]{dto.Work: workHandler, dto.Result: resultHandler}

	dispatcher := server.NewDispatcher(handlers, dto.NewRequestTransformer[any]())
	workServer, err := zmq.NewServer(dispatcher, workServerConfig)

	if err != nil {
		log.Fatal("can't create server", err)
	}
	wWg, err := workServer.Start()
	if err != nil {
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
