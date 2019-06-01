package endpoint

import (
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
)

type CommonEndpoint struct {
	adaptors      *adaptors.Adaptors
	core          *core.Core
	accessList    *access_list.AccessListService
	scriptService *scripts.ScriptService
}

func NewCommonEndpoint(adaptors *adaptors.Adaptors,
	core *core.Core,
	accessList *access_list.AccessListService,
	scriptService *scripts.ScriptService, ) *CommonEndpoint {
	return &CommonEndpoint{
		adaptors:      adaptors,
		core:          core,
		accessList:    accessList,
		scriptService: scriptService,
	}
}
