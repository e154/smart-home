package main

import (
	"github.com/astaxie/beego"
	"./api"
	"./static"
	"time"
)

func main() {
	api.Initialize()
	static.Initialize()
	beego.Info("Starting....")
	go beego.Run()

	for ;; {
		time.Sleep(time.Second * 1)
	}
}
