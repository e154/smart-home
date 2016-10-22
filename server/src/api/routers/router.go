package routers

import (
	"../controllers"
	"github.com/astaxie/beego"
)

func Initialize() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSRouter("/node/:id([0-9]+)", &controllers.NodeController{}, "get:GetOne"),
			beego.NSRouter("/node", &controllers.NodeController{}, "get:GetAll"),
			beego.NSRouter("/node", &controllers.NodeController{}, "post:Post"),
			beego.NSRouter("/node/:id([0-9]+)", &controllers.NodeController{}, "put:Put"),
			beego.NSRouter("/node/:id([0-9]+)", &controllers.NodeController{}, "delete:Delete"),

		),
	)
	beego.AddNamespace(ns)
}
