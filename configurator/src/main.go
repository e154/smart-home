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

	data_dir := beego.AppConfig.String("data_dir")

	if(beego.AppConfig.String("runmode") == "dev") {
		beego.Info("Develment mode enabled")

		beego.SetStaticPath("/static", "../build/private")
		beego.SetStaticPath("/_static", "../build/public")
		beego.SetStaticPath("/attach", data_dir)
	} else {
		beego.Info("Product mode enabled")

		beego.SetStaticPath("/static", "build/private")
		beego.SetStaticPath("/_static", "build/public")
		beego.SetStaticPath("/attach", data_dir)
	}
}