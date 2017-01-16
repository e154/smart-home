package router

import (
	"github.com/astaxie/beego"
	"../controllers"
)

func Initialize() {
	beego.ErrorController(&controllers.ErrorController{})

	beego.Router("/", &controllers.DashboardController{}, "get:Index")
	beego.Router("/login", &controllers.DashboardController{}, "*:Login")
	beego.Router("/logout", &controllers.DashboardController{}, "*:Logout")
	beego.Router("/recovery", &controllers.DashboardController{}, "get:Recovery")
	beego.Router("/reset", &controllers.DashboardController{}, "get:Reset")
	beego.Router("/reset", &controllers.DashboardController{}, "post:ResetPost")
	beego.Router("/*", &controllers.DashboardController{}, "*:Index")
}