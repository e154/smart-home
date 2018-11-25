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
	"github.com/e154/smart-home/system/backup"
	"fmt"
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
		case "-backup":
			container := BuildContainer()
			container.Invoke(func(
				backup *backup.Backup) {

				if err := backup.New(); err != nil {
					fmt.Println(err.Error())
				}
			})
			return
		case "-restore":
			if len(os.Args) <3 {
				fmt.Println("need backup name")
				return
			}
			container := BuildContainer()
			container.Invoke(func(
				backup *backup.Backup) {

				if err := backup.Restore(os.Args[2]); err != nil {
					fmt.Println(err.Error())
				}
			})
			return
		case "-reset":
			container := BuildContainer()
			container.Invoke(func(
				initialService *initial.InitialService) {

				initialService.Reset()
			})
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
		lx *logrus.Logger,
		initialService *initial.InitialService) {

		l.Initialize(lx)
		go server.Start()
		go core.Run()

		graceful.Wait()
	})
}
