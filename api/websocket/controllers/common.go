package controllers

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
)

type ControllerCommon struct {
	adaptors *adaptors.Adaptors
	stream   *stream.StreamService
	endpoint *endpoint.Endpoint
	core     *core.Core
	scripts  *scripts.ScriptService
}

func NewControllerCommon(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	endpoint *endpoint.Endpoint,
	scripts *scripts.ScriptService,
	core *core.Core) *ControllerCommon {
	return &ControllerCommon{
		adaptors: adaptors,
		endpoint: endpoint,
		stream:   stream,
		core:     core,
		scripts:  scripts,
	}
}

func (c *ControllerCommon) Err(client *stream.Client, message stream.Message, err error) {
	msg, _ := json.Marshal(map[string]interface{}{"id": message.Id, "error": err.Error()})
	client.Send <- msg
}
