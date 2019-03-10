package server

import (
	"github.com/e154/smart-home/system/swaggo/gin-swagger/swaggerFiles"
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/common"
)

func (s *Server) setControllers() {

	r := s.engine

	r.Static("/upload", common.StoragePath())

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
	v1.POST("/script/:id/exec_src", s.af.Auth, s.ControllersV1.Script.ExecSrc)
	v1.GET("/scripts/search", s.af.Auth, s.ControllersV1.Script.Search)

	// workflow
	v1.POST("/workflow", s.af.Auth, s.ControllersV1.Workflow.AddWorkflow)
	v1.GET("/workflow/:id", s.af.Auth, s.ControllersV1.Workflow.GetWorkflowById)
	v1.PUT("/workflow/:id", s.af.Auth, s.ControllersV1.Workflow.UpdateWorkflow)
	v1.DELETE("/workflow/:id", s.af.Auth, s.ControllersV1.Workflow.DeleteWorkflowById)
	v1.GET("/workflows", s.af.Auth, s.ControllersV1.Workflow.GetWorkflowList)
	v1.GET("/workflows/search", s.af.Auth, s.ControllersV1.Workflow.Search)
	v1.PUT("/workflow/:id/update_scenario", s.af.Auth, s.ControllersV1.Workflow.UpdateScenario)

	// workflow scenario
	v1.GET("/workflow/:id/scenario/:scenario_id", s.af.Auth, s.ControllersV1.WorkflowScenario.GetWorkflowScenarioById)
	v1.GET("/workflow/:id/scenarios", s.af.Auth, s.ControllersV1.WorkflowScenario.GetWorkflowScenarioList)
	v1.POST("/workflow/:id/scenario", s.af.Auth, s.ControllersV1.WorkflowScenario.AddWorkflowScenario)
	v1.PUT("/workflow/:id/scenario/:scenario_id", s.af.Auth, s.ControllersV1.WorkflowScenario.UpdateWorkflowScenario)
	v1.DELETE("/workflow/:id/scenario/:scenario_id", s.af.Auth, s.ControllersV1.WorkflowScenario.DeleteWorkflowScenarioById)
	v1.GET("/workflow/:id/scenarios/search", s.af.Auth, s.ControllersV1.WorkflowScenario.Search)

	// device
	v1.POST("/device", s.af.Auth, s.ControllersV1.Device.AddDevice)
	v1.GET("/device/:id", s.af.Auth, s.ControllersV1.Device.GetDeviceById)
	v1.PUT("/device/:id", s.af.Auth, s.ControllersV1.Device.UpdateDevice)
	v1.DELETE("/device/:id", s.af.Auth, s.ControllersV1.Device.DeleteDeviceById)
	v1.GET("/devices", s.af.Auth, s.ControllersV1.Device.GetDeviceList)
	v1.GET("/devices/search", s.af.Auth, s.ControllersV1.Device.Search)

	// device actions
	v1.POST("/device_action", s.af.Auth, s.ControllersV1.DeviceAction.AddAction)
	v1.GET("/device_action/:id", s.af.Auth, s.ControllersV1.DeviceAction.GetById)
	v1.PUT("/device_action/:id", s.af.Auth, s.ControllersV1.DeviceAction.Update)
	v1.DELETE("/device_action/:id", s.af.Auth, s.ControllersV1.DeviceAction.DeleteById)
	v1.GET("/device_actions/:id", s.af.Auth, s.ControllersV1.DeviceAction.GetDeviceActionList)
	v1.GET("/device_action1/search", s.af.Auth, s.ControllersV1.DeviceAction.Search)

	// device states
	v1.POST("/device_state", s.af.Auth, s.ControllersV1.DeviceState.AddState)
	v1.GET("/device_state/:id", s.af.Auth, s.ControllersV1.DeviceState.GetById)
	v1.PUT("/device_state/:id", s.af.Auth, s.ControllersV1.DeviceState.Update)
	v1.DELETE("/device_state/:id", s.af.Auth, s.ControllersV1.DeviceState.DeleteById)
	v1.GET("/device_states/:id", s.af.Auth, s.ControllersV1.DeviceState.GetDeviceStateList)

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

	// maps
	v1.POST("/map", s.af.Auth, s.ControllersV1.Map.Add)
	v1.GET("/map/:id", s.af.Auth, s.ControllersV1.Map.GetById)
	v1.GET("/map/:id/full", s.af.Auth, s.ControllersV1.Map.GetFullMap)
	v1.PUT("/map/:id", s.af.Auth, s.ControllersV1.Map.Update)
	v1.DELETE("/map/:id", s.af.Auth, s.ControllersV1.Map.Delete)
	v1.GET("/maps", s.af.Auth, s.ControllersV1.Map.GetList)
	v1.GET("/maps/search", s.af.Auth, s.ControllersV1.Map.Search)

	// map_layer
	v1.POST("/map_layer", s.af.Auth, s.ControllersV1.MapLayer.Add)
	v1.GET("/map_layer/:id", s.af.Auth, s.ControllersV1.MapLayer.GetById)
	v1.PUT("/map_layer/:id", s.af.Auth, s.ControllersV1.MapLayer.Update)
	v1.DELETE("/map_layer/:id", s.af.Auth, s.ControllersV1.MapLayer.Delete)
	v1.GET("/map_layers", s.af.Auth, s.ControllersV1.MapLayer.GetList)
	v1.PUT("/map_layers/sort", s.af.Auth, s.ControllersV1.MapLayer.Sort)

	// map element
	v1.POST("/map_element", s.af.Auth, s.ControllersV1.MapElement.Add)
	v1.GET("/map_element/:id", s.af.Auth, s.ControllersV1.MapElement.GetById)
	v1.PUT("/map_element/:id", s.af.Auth, s.ControllersV1.MapElement.Update)
	v1.DELETE("/map_element/:id", s.af.Auth, s.ControllersV1.MapElement.Delete)
	v1.GET("/map_elements", s.af.Auth, s.ControllersV1.MapElement.GetList)
	v1.PUT("/map_elements/sort", s.af.Auth, s.ControllersV1.MapElement.Sort)

	// images
	v1.POST("/image", s.af.Auth, s.ControllersV1.Image.Add)
	v1.GET("/image/:id", s.af.Auth, s.ControllersV1.Image.GetById)
	v1.GET("/images", s.af.Auth, s.ControllersV1.Image.GetList)
	v1.POST("/image/upload", s.af.Auth, s.ControllersV1.Image.Upload)
	v1.PUT("/image/:id", s.af.Auth, s.ControllersV1.Image.UpdateImage)
	v1.DELETE("/image/:id", s.af.Auth, s.ControllersV1.Image.Delete)

	// flow
	v1.POST("/flow", s.af.Auth, s.ControllersV1.Flow.Add)
	v1.GET("/flow/:id", s.af.Auth, s.ControllersV1.Flow.GetById)
	v1.GET("/flows", s.af.Auth, s.ControllersV1.Flow.GetList)
	v1.GET("/flow/:id/redactor", s.af.Auth, s.ControllersV1.Flow.GetRedactor)
	v1.PUT("/flow/:id/redactor", s.af.Auth, s.ControllersV1.Flow.UpdateRedactor)
	v1.PUT("/flow/:id", s.af.Auth, s.ControllersV1.Flow.Update)
	v1.DELETE("/flow/:id", s.af.Auth, s.ControllersV1.Flow.Delete)
	v1.GET("/flows/search", s.af.Auth, s.ControllersV1.Flow.Search)

	// logs
	v1.POST("/log", s.af.Auth, s.ControllersV1.Log.Add)
	v1.GET("/log/:id", s.af.Auth, s.ControllersV1.Log.GetById)
	v1.DELETE("/log/:id", s.af.Auth, s.ControllersV1.Log.Delete)
	v1.GET("/logs", s.af.Auth, s.ControllersV1.Log.GetList)
	v1.GET("/logs/search", s.af.Auth, s.ControllersV1.Log.Search)

}
