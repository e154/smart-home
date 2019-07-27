package endpoint

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/scripts"
)

type CommonEndpoint struct {
	adaptors      *adaptors.Adaptors
	core          *core.Core
	accessList    *access_list.AccessListService
	scriptService *scripts.ScriptService
	gate          *gate_client.GateClient
}

func NewCommonEndpoint(adaptors *adaptors.Adaptors,
	core *core.Core,
	accessList *access_list.AccessListService,
	scriptService *scripts.ScriptService,
	gate *gate_client.GateClient, ) *CommonEndpoint {
	return &CommonEndpoint{
		adaptors:      adaptors,
		core:          core,
		accessList:    accessList,
		scriptService: scriptService,
		gate:          gate,
	}
}
