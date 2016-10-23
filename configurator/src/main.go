package main

import (
	"github.com/astaxie/beego"
	"./api"
)

func main() {
	api.Initialize()
	beego.Info("Starting....")
	beego.Run()
}
