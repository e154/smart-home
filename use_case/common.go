package use_case

import (
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
)

type CommonCommand struct {
	adaptors      *adaptors.Adaptors
	core          *core.Core
	accessList    *access_list.AccessListService
	scriptService *scripts.ScriptService
}

func NewCommonCommand(adaptors *adaptors.Adaptors,
	core *core.Core,
	accessList *access_list.AccessListService,
	scriptService *scripts.ScriptService, ) *CommonCommand {
	return &CommonCommand{
		adaptors:      adaptors,
		core:          core,
		accessList:    accessList,
		scriptService: scriptService,
	}
}
