package controllers

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/access_list"
)

type ControllersV1 struct {
	Index            *ControllerIndex
	Node             *ControllerNode
	Swagger          *ControllerSwagger
	Script           *ControllerScript
	Workflow         *ControllerWorkflow
	Device           *ControllerDevice
	Role             *ControllerRole
	User             *ControllerUser
	Auth             *ControllerAuth
	DeviceAction     *ControllerDeviceAction
	DeviceState      *ControllerDeviceState
	Map              *ControllerMap
	MapLayer         *ControllerMapLayer
	Image            *ControllerImage
	WorkflowScenario *ControllerWorkflowScenario
	Flow             *ControllerFlow
	Log              *ControllerLog
}

func NewControllersV1(adaptors *adaptors.Adaptors,
	core *core.Core,
	scriptService *scripts.ScriptService,
	accessList *access_list.AccessListService) *ControllersV1 {
	common := NewControllerCommon(adaptors, core, accessList)
	return &ControllersV1{
		Index:            NewControllerIndex(common),
		Node:             NewControllerNode(common),
		Swagger:          NewControllerSwagger(common),
		Script:           NewControllerScript(common, scriptService),
		Workflow:         NewControllerWorkflow(common),
		Device:           NewControllerDevice(common),
		Role:             NewControllerRole(common),
		User:             NewControllerUser(common),
		Auth:             NewControllerAuth(common),
		DeviceAction:     NewControllerDeviceAction(common),
		DeviceState:      NewControllerDeviceState(common),
		Map:              NewControllerMap(common),
		MapLayer:         NewControllerMapLayer(common),
		Image:            NewControllerImage(common),
		WorkflowScenario: NewControllerWorkflowScenario(common),
		Flow:             NewControllerFlow(common),
		Log:              NewControllerLog(common),
	}
}
