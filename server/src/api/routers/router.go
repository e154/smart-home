package routers

import (
	"../controllers"
	"github.com/astaxie/beego"
	"../stream"
)

func Initialize() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSRouter("/ws)", &stream.StreamCotroller{}, "get:Get"),
			beego.NSRouter("/ws/*)", &stream.StreamCotroller{}, "get:Get"),

			beego.NSRouter("/node/:id([0-9]+)", &controllers.NodeController{}, "get:GetOne"),
			beego.NSRouter("/node", &controllers.NodeController{}, "get:GetAll"),
			beego.NSRouter("/node", &controllers.NodeController{}, "post:Post"),
			beego.NSRouter("/node/:id([0-9]+)", &controllers.NodeController{}, "put:Put"),
			beego.NSRouter("/node/:id([0-9]+)", &controllers.NodeController{}, "delete:Delete"),

			beego.NSRouter("/device/:id([0-9]+)", &controllers.DeviceController{}, "get:GetOne"),
			beego.NSRouter("/device", &controllers.DeviceController{}, "get:GetAll"),
			beego.NSRouter("/device", &controllers.DeviceController{}, "post:Post"),
			beego.NSRouter("/device/:id([0-9]+)", &controllers.DeviceController{}, "put:Put"),
			beego.NSRouter("/device/:id([0-9]+)", &controllers.DeviceController{}, "delete:Delete"),

			beego.NSRouter("/workflow/:id([0-9]+)", &controllers.WorkflowController{}, "get:GetOne"),
			beego.NSRouter("/workflow", &controllers.WorkflowController{}, "get:GetAll"),
			beego.NSRouter("/workflow", &controllers.WorkflowController{}, "post:Post"),
			beego.NSRouter("/workflow/:id([0-9]+)", &controllers.WorkflowController{}, "put:Put"),
			beego.NSRouter("/workflow/:id([0-9]+)", &controllers.WorkflowController{}, "delete:Delete"),

			beego.NSRouter("/flow/:id([0-9]+)", &controllers.FlowController{}, "get:GetOne"),
			beego.NSRouter("/flow/:id([0-9]+)/full", &controllers.FlowController{}, "get:GetOneFull"),
			beego.NSRouter("/flow/:id([0-9]+)/redactor", &controllers.FlowController{}, "get:GetOneRedactor"),
			beego.NSRouter("/flow/:id([0-9]+)/redactor", &controllers.FlowController{}, "put:UpdateRedactor"),
			beego.NSRouter("/flow", &controllers.FlowController{}, "get:GetAll"),
			beego.NSRouter("/flow", &controllers.FlowController{}, "post:Post"),
			beego.NSRouter("/flow/:id([0-9]+)", &controllers.FlowController{}, "put:Put"),
			beego.NSRouter("/flow/:id([0-9]+)", &controllers.FlowController{}, "delete:Delete"),
		),
	)
	beego.AddNamespace(ns)
}
