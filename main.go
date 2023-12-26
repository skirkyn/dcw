package main

import (
	"github.com/skirkyn/dcw/cmd/controller"
	"github.com/skirkyn/dcw/cmd/worker"
	"log"
	"os"
)

func main() {

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("mode should be specified")
	}

	if args[0] == "worker" {
		worker.Run()
	} else if args[0] == "controller" {
		controller.Run()
	}
}
