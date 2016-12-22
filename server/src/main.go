package main

import (
	"github.com/astaxie/beego"
	"./api"
	"time"
	"./api/log"
)

func main() {
	api.Initialize()
	log.Info("Starting....")
	go beego.Run()

	for ;; {
		time.Sleep(time.Second * 1)
	}
}
