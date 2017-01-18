package router

import (
	"github.com/astaxie/beego"
	"../controllers"
)

func Initialize() {
	beego.ErrorController(&controllers.ErrorController{})

	beego.Router("/", &controllers.DashboardController{}, "*:Index")
	beego.Router("/signin", &controllers.DashboardController{}, "*:Signin")
	beego.Router("/signout", &controllers.DashboardController{}, "get:Signout")
	beego.Router("/recovery", &controllers.DashboardController{}, "post:Recovery")
	beego.Router("/reset", &controllers.DashboardController{}, "post:Reset")
	beego.Router("/*", &controllers.DashboardController{}, "*:Index")
}