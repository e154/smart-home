package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/takama/daemon"
	"github.com/e154/smart-home/api"
	"github.com/astaxie/beego"
)

const (
	name        = "smart-home-server"
	description = "Smart Home Server"
)

var dependencies = []string{"mysql.service"}

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	api.Initialize()
	var port int
	port = beego.BConfig.Listen.HTTPPort

	for {
		select {
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ", port)

			if killSignal == os.Interrupt {
				return "Server was interruped by system signal", nil
			}
			return "Server was killed", nil
		}
	}
}

func ServiceInitialize() {
	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
