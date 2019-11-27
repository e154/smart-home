package endpoint

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/notify"
	"github.com/e154/smart-home/system/scripts"
)

type CommonEndpoint struct {
	adaptors      *adaptors.Adaptors
	core          *core.Core
	accessList    *access_list.AccessListService
	scriptService *scripts.ScriptService
	gate          *gate_client.GateClient
	notify        *notify.Notify
	mqtt          *mqtt.Mqtt
}

func NewCommonEndpoint(adaptors *adaptors.Adaptors,
	core *core.Core,
	accessList *access_list.AccessListService,
	scriptService *scripts.ScriptService,
	gate *gate_client.GateClient,
	notify *notify.Notify,
	mqtt *mqtt.Mqtt) *CommonEndpoint {
	return &CommonEndpoint{
		adaptors:      adaptors,
		core:          core,
		accessList:    accessList,
		scriptService: scriptService,
		gate:          gate,
		notify:        notify,
		mqtt:          mqtt,
	}
}
