package main

import (
	"github.com/astaxie/beego"
	"./controllers"
)

func main() {
	Initialize()
	beego.Info("Starting....")
	beego.Run()
}

func Initialize() {
	beego.Router("/", &controllers.DashboardController{}, "get:Index")
	beego.Router("/*", &controllers.DashboardController{}, "get:Index")

	beego.Info("AppPath:", beego.AppPath)
	if(beego.AppConfig.String("runmode") == "dev") {
		beego.Info("Develment mode enabled")

		beego.SetStaticPath("/static", "static_source")
		beego.SetStaticPath("/attach", "../../data")
	} else {
		beego.Info("Product mode enabled")

		beego.SetStaticPath("/static", "static_source")
		beego.SetStaticPath("/attach", "../data")
	}
}