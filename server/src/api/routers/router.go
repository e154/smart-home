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
			beego.NSRouter("/device/group", &controllers.DeviceController{}, "get:GetGroup"),
			beego.NSRouter("/device/:id([0-9]+)/actions", &controllers.DeviceController{}, "get:GetActions"),

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
			beego.NSRouter("/flow/:id([0-9]+)/workers", &controllers.FlowController{}, "get:GetWorkers"),

			beego.NSRouter("/device_action/:id([0-9]+)", &controllers.DeviceActionController{}, "get:GetOne"),
			beego.NSRouter("/device_action", &controllers.DeviceActionController{}, "get:GetAll"),
			beego.NSRouter("/device_action", &controllers.DeviceActionController{}, "post:Post"),
			beego.NSRouter("/device_action/:id([0-9]+)", &controllers.DeviceActionController{}, "put:Put"),
			beego.NSRouter("/device_action/:id([0-9]+)", &controllers.DeviceActionController{}, "delete:Delete"),
			beego.NSRouter("/device_action/search", &controllers.DeviceActionController{}, "get:Search"),

			beego.NSRouter("/worker/:id([0-9]+)", &controllers.WorkerController{}, "get:GetOne"),
			beego.NSRouter("/worker", &controllers.WorkerController{}, "get:GetAll"),
			beego.NSRouter("/worker", &controllers.WorkerController{}, "post:Post"),
			beego.NSRouter("/worker/:id([0-9]+)", &controllers.WorkerController{}, "put:Put"),
			beego.NSRouter("/worker/:id([0-9]+)", &controllers.WorkerController{}, "delete:Delete"),

			beego.NSRouter("/script/:id([0-9]+)", &controllers.ScriptController{}, "get:GetOne"),
			beego.NSRouter("/script", &controllers.ScriptController{}, "get:GetAll"),
			beego.NSRouter("/script", &controllers.ScriptController{}, "post:Post"),
			beego.NSRouter("/script/:id([0-9]+)", &controllers.ScriptController{}, "put:Put"),
			beego.NSRouter("/script/:id([0-9]+)", &controllers.ScriptController{}, "delete:Delete"),
			beego.NSRouter("/script/:id([0-9]+)/exec", &controllers.ScriptController{}, "post:Exec"),
		),
	)
	beego.AddNamespace(ns)
}
