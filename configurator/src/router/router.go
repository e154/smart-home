package router

import (
	"github.com/astaxie/beego"
	"../controllers"
)

func Initialize() {
	beego.ErrorController(&controllers.ErrorController{})

	beego.Router("/", &controllers.DashboardController{}, "*:Index")
	beego.Router("/signin", &controllers.DashboardController{}, "*:Signin")
	beego.Router("/signout", &controllers.DashboardController{}, "*:Signout")
	beego.Router("/recovery", &controllers.DashboardController{}, "*:Recovery")
	beego.Router("/reset", &controllers.DashboardController{}, "*:Reset")
	beego.Router("/*", &controllers.DashboardController{}, "*:Index")
}