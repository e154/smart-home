package main

import (
	"time"
	"github.com/astaxie/beego"
	"./api"
)

func main() {
	api.Initialize()

	beego.Info("Starting....")
	go beego.Run()

	for ;; {
		time.Sleep(time.Second * 1)
	}
}
