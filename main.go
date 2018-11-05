package main

import (
	"os"
	"github.com/e154/smart-home/api/server_v1"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("main")
)

func main() {

	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			log.Info(GetHumanVersion())
			return
		}
	}

	start()
}

func start() {

	container := BuildContainer()
	container.Invoke(func(server *server.Server) {
		server.Start()
	})
}
