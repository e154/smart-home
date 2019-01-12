package server

import (
	"github.com/e154/smart-home/system/swaggo/gin-swagger/swaggerFiles"
	"github.com/gin-gonic/gin"
)

func (s *Server) setControllers() {

	r := s.engine
	basePath := r.Group("/api")

	v1 := basePath.Group("/v1")
	v1.GET("/", s.ControllersV1.Index.Index)
	v1.GET("/swagger", func(context *gin.Context) {
		context.Redirect(302, "/api/v1/swagger/index.html")
	})
	v1.GET("/swagger/*any", s.ControllersV1.Swagger.WrapHandler(swaggerFiles.Handler))

	// ws
	v1.GET("/ws", s.af.Auth, s.streamService.Ws)
	v1.GET("/ws/*any", s.af.Auth, s.streamService.Ws)

	// auth
	v1.POST("/signin", s.ControllersV1.Auth.SignIn)
	v1.POST("/signout", s.af.Auth, s.ControllersV1.Auth.SignOut)
	v1.POST("/recovery", s.af.Auth, s.ControllersV1.Auth.Recovery)
	v1.POST("/reset", s.af.Auth, s.ControllersV1.Auth.Reset)
	v1.GET("/access_list", s.af.Auth, s.ControllersV1.Auth.AccessList)

	// nodes
	v1.POST("/node", s.af.Auth, s.ControllersV1.Node.AddNode)
	v1.GET("/node/:id", s.af.Auth, s.ControllersV1.Node.GetNodeById)
	v1.PUT("/node/:id", s.af.Auth, s.ControllersV1.Node.UpdateNode)
	v1.DELETE("/node/:id", s.af.Auth, s.ControllersV1.Node.DeleteNodeById)
	v1.GET("/nodes", s.af.Auth, s.ControllersV1.Node.GetNodeList)
	v1.GET("/nodes/search", s.af.Auth, s.ControllersV1.Node.Search)

	// scripts
	v1.POST("/script", s.af.Auth, s.ControllersV1.Script.AddScript)
	v1.GET("/script/:id", s.af.Auth, s.ControllersV1.Script.GetScriptById)
	v1.PUT("/script/:id", s.af.Auth, s.ControllersV1.Script.UpdateScript)
	v1.DELETE("/script/:id", s.af.Auth, s.ControllersV1.Script.DeleteScriptById)
	v1.GET("/scripts", s.af.Auth, s.ControllersV1.Script.GetScriptList)
	v1.POST("/script/:id/exec", s.af.Auth, s.ControllersV1.Script.Exec)
	v1.GET("/scripts/search", s.af.Auth, s.ControllersV1.Script.Search)

	// workflow
	v1.POST("/workflow", s.af.Auth, s.ControllersV1.Workflow.AddWorkflow)
	v1.GET("/workflow/:id", s.af.Auth, s.ControllersV1.Workflow.GetWorkflowById)
	v1.PUT("/workflow/:id", s.af.Auth, s.ControllersV1.Workflow.UpdateWorkflow)
	v1.DELETE("/workflow/:id", s.af.Auth, s.ControllersV1.Workflow.DeleteWorkflowById)
	v1.GET("/workflow", s.af.Auth, s.ControllersV1.Workflow.GetWorkflowList)

	// device
	v1.POST("/device", s.af.Auth, s.ControllersV1.Device.AddDevice)
	v1.GET("/device/:id", s.af.Auth, s.ControllersV1.Device.GetDeviceById)
	v1.PUT("/device/:id", s.af.Auth, s.ControllersV1.Device.UpdateDevice)
	v1.DELETE("/device/:id", s.af.Auth, s.ControllersV1.Device.DeleteDeviceById)
	v1.GET("/devices", s.af.Auth, s.ControllersV1.Device.GetDeviceList)
	v1.GET("/devices/search", s.af.Auth, s.ControllersV1.Device.Search)

	// role
	v1.POST("/role", s.af.Auth, s.ControllersV1.Role.AddRole)
	v1.GET("/role/:name", s.af.Auth, s.ControllersV1.Role.GetRoleByName)
	v1.GET("/role/:name/access_list", s.af.Auth, s.ControllersV1.Role.GetAccessList)
	v1.PUT("/role/:name/access_list", s.af.Auth, s.ControllersV1.Role.UpdateAccessList)
	v1.PUT("/role/:name", s.af.Auth, s.ControllersV1.Role.UpdateRole)
	v1.DELETE("/role/:name", s.af.Auth, s.ControllersV1.Role.DeleteRoleByName)
	v1.GET("/roles", s.af.Auth, s.ControllersV1.Role.GetRoleList)
	v1.GET("/roles/search", s.af.Auth, s.ControllersV1.Role.Search)

	// user
	v1.POST("/user", s.af.Auth, s.ControllersV1.User.AddUser)
	v1.GET("/user/:id", s.af.Auth, s.ControllersV1.User.GetUserById)
	v1.PUT("/user/:id", s.af.Auth, s.ControllersV1.User.UpdateUser)
	v1.DELETE("/user/:id", s.af.Auth, s.ControllersV1.User.DeleteUserById)
	v1.PUT("/user/:id/update_status", s.af.Auth, s.ControllersV1.User.UpdateStatus)
	v1.GET("/users", s.af.Auth, s.ControllersV1.User.GetUserList)
}
