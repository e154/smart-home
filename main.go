package main

import (
	"os"
	"github.com/op/go-logging"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/sirupsen/logrus"
	l "github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/system/initial"
)

var (
	log = logging.MustGetLogger("main")
)

func main() {

	args := os.Args[1:]
	for _, arg := range args {
		switch arg {
		case "-v", "--version":
			log.Info(GetHumanVersion())
			return
		case "-r":
			container := BuildContainer()
			container.Invoke(func(server *server.Server,
				graceful *graceful_service.GracefulService,
				lx *logrus.Logger,
				initialService *initial.InitialService) {

				l.Initialize(lx)
				initialService.Reset()
			})
		}
	}

	start()
}

func start() {

	container := BuildContainer()
	container.Invoke(func(server *server.Server,
		core *core.Core,
		graceful *graceful_service.GracefulService,
		lx *logrus.Logger,
		initialService *initial.InitialService) {

		l.Initialize(lx)
		go server.Start()
		go core.Run()

		graceful.Wait()
	})
}
