package main

import (
	"os"
	"github.com/op/go-logging"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/sirupsen/logrus"
	l "github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/api/server"
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
	container.Invoke(func(server *server.Server,
		core *core.Core,
		graceful *graceful_service.GracefulService,
		lx *logrus.Logger) {

		l.Initialize(lx)
		go server.Start()
		go core.Run()

		graceful.Wait()
	})
}
