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

	// nodes
	v1.POST("/node", s.ControllersV1.Node.AddNode)
	v1.GET("/node/:id", s.ControllersV1.Node.GetNodeById)
	v1.PUT("/node/:id", s.ControllersV1.Node.UpdateNode)
	v1.DELETE("/node/:id", s.ControllersV1.Node.DeleteNodeById)
	v1.GET("/node", s.ControllersV1.Node.GetNodeList)

	// scripts
	v1.POST("/script", s.ControllersV1.Script.AddScript)
	v1.GET("/script/:id", s.ControllersV1.Script.GetScriptById)
	v1.PUT("/script/:id", s.ControllersV1.Script.UpdateScript)
	v1.DELETE("/script/:id", s.ControllersV1.Script.DeleteScriptById)
	v1.GET("/script", s.ControllersV1.Script.GetScriptList)
	v1.POST("/script/:id/exec", s.ControllersV1.Script.Exec)

	// workflow
	v1.POST("/workflow", s.ControllersV1.Workflow.AddWorkflow)
	v1.GET("/workflow/:id", s.ControllersV1.Workflow.GetWorkflowById)
	v1.PUT("/workflow/:id", s.ControllersV1.Workflow.UpdateWorkflow)
	v1.DELETE("/workflow/:id", s.ControllersV1.Workflow.DeleteWorkflowById)
	v1.GET("/workflow", s.ControllersV1.Workflow.GetWorkflowList)

	// device
	v1.POST("/device", s.ControllersV1.Device.AddDevice)
	v1.GET("/device/:id", s.ControllersV1.Device.GetDeviceById)
	v1.PUT("/device/:id", s.ControllersV1.Device.UpdateDevice)
	v1.DELETE("/device/:id", s.ControllersV1.Device.DeleteDeviceById)
	v1.GET("/device", s.ControllersV1.Device.GetDeviceList)

	// role
	v1.POST("/role", s.ControllersV1.Role.AddRole)
	v1.GET("/role/:name", s.ControllersV1.Role.GetRoleByName)
	v1.GET("/role/:name/access_list", s.ControllersV1.Role.GetAccessList)
	v1.PUT("/role/:name/access_list", s.ControllersV1.Role.UpdateAccessList)
	v1.PUT("/role/:name", s.ControllersV1.Role.UpdateRole)
	v1.DELETE("/role/:name", s.ControllersV1.Role.DeleteRoleByName)
	v1.GET("/roles", s.ControllersV1.Role.GetRoleList)
	v1.GET("/roles/search", s.ControllersV1.Role.Search)

	// user
	v1.POST("/user", s.ControllersV1.User.AddUser)
	v1.GET("/user/:id", s.ControllersV1.User.GetUserById)
}
