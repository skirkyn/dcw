package main

import (
	"github.com/unknownfeature/dcw/cmd/common/config"
	"log"
)

func Run() {

	commonConfig, err := config.ReadCommonConfig()

	if err != nil {
		log.Fatal("can't read common config", err)
	}

	// todo it should be a factory based on config but it is what it is
	rnr, err := getRunner(commonConfig)
	if err != nil {
		log.Fatal("can't create runner", err)
	}
	wg := rnr.Start()
	wg.Wait()
}
