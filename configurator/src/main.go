package main

import (
	"github.com/astaxie/beego"
	"./router"
)

func main() {
	Initialize()
	beego.Run()
}

func Initialize() {
	router.Initialize()

	beego.Info("Starting....")
	beego.Info("AppPath:", beego.AppPath)
	if(beego.AppConfig.String("runmode") == "dev") {
		beego.Info("Develment mode enabled")

		beego.SetStaticPath("/static", "../build/private")
		beego.SetStaticPath("/attach", "../../data")
	} else {
		beego.Info("Product mode enabled")

		beego.SetStaticPath("/static", "build/private")
		beego.SetStaticPath("/attach", "../data")
	}
}