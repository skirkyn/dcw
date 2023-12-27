package main

import (
	"github.com/unknownfeature/dcw/cmd/common/config"
	"github.com/unknownfeature/dcw/cmd/controller/server"
	"github.com/unknownfeature/dcw/cmd/controller/server/zmq"
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
	dispatcher, err := getDispatcher(commonConfig, controllerConfig)
	if err != nil {
		log.Fatal("can't create dispatcher", err)
	}
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
