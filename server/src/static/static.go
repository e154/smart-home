package static

import (
	"github.com/astaxie/beego"
	"../controllers"
)

func Initialize() {
	beego.Router("/", &controllers.DashboardController{}, "get:Index")
	beego.Router("/*", &controllers.DashboardController{}, "get:Index")

	beego.Info("AppPath:", beego.AppPath)
	staticDir := beego.AppConfig.String("staticDir")

	if(beego.AppConfig.String("runmode") == "dev") {
		beego.SetStaticPath("/static", staticDir + "/static_source")
		beego.SetStaticPath("/attach", staticDir + "/../data")
	} else {
		beego.SetStaticPath("/static", staticDir + "/static_source")
		beego.SetStaticPath("/attach", staticDir + "/../data")
	}
}