package main

import (
	"github.com/astaxie/beego"
	"./api"
	"time"
)

func main() {
	api.Initialize()
	beego.Info("Starting....")
	go beego.Run()

	for ;; {
		time.Sleep(time.Second * 1)
	}
}
