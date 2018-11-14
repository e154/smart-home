package controllers

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
)

type ControllersV1 struct {
	Index    *ControllerIndex
	Node     *ControllerNode
	Swagger  *ControllerSwagger
	Script   *ControllerScript
	Workflow *ControllerWorkflow
	Device   *ControllerDevice
}

func NewControllersV1(adaptors *adaptors.Adaptors,
	core *core.Core,
	scriptService *scripts.ScriptService) *ControllersV1 {
	common := NewControllerCommon(adaptors, core)
	return &ControllersV1{
		Index:    NewControllerIndex(common),
		Node:     NewControllerNode(common),
		Swagger:  NewControllerSwagger(common),
		Script:   NewControllerScript(common, scriptService),
		Workflow: NewControllerWorkflow(common),
		Device:   NewControllerDevice(common),
	}
}
