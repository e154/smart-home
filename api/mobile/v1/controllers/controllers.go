package controllers

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
)

type MobileControllersV1 struct {
	Auth     *ControllerAuth
	Workflow *ControllerWorkflow
	Map      *ControllerMap
}

func NewMobileControllersV1(adaptors *adaptors.Adaptors,
	core *core.Core,
	scriptService *scripts.ScriptService,
	accessList *access_list.AccessListService,
	command *endpoint.Endpoint) *MobileControllersV1 {
	common := NewControllerCommon(adaptors, core, accessList, command)
	return &MobileControllersV1{
		Auth:     NewControllerAuth(common),
		Workflow: NewControllerWorkflow(common),
		Map:      NewControllerMap(common),
	}
}
