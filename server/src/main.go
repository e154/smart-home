package main

import (
	"github.com/astaxie/beego"
	"./api"
	"./static"
	"time"
	"./api/log"
)

func main() {
	api.Initialize()
	static.Initialize()
	log.Info("Starting....")
	go beego.Run()

	for ;; {
		time.Sleep(time.Second * 1)
	}
}
