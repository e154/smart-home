package controllers

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/telemetry"
)

type ControllerCommon struct {
	adaptors  *adaptors.Adaptors
	stream    *stream.StreamService
	endpoint  *endpoint.Endpoint
	core      *core.Core
	scripts   *scripts.ScriptService
	telemetry *telemetry.Telemetry
	gate      *gate_client.GateClient
}

func NewControllerCommon(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	endpoint *endpoint.Endpoint,
	scripts *scripts.ScriptService,
	core *core.Core,
	telemetry *telemetry.Telemetry,
	gate *gate_client.GateClient) *ControllerCommon {
	return &ControllerCommon{
		adaptors:  adaptors,
		endpoint:  endpoint,
		stream:    stream,
		core:      core,
		scripts:   scripts,
		telemetry: telemetry,
		gate:      gate,
	}
}

func (c *ControllerCommon) Err(client *stream.Client, message stream.Message, err error) {
	msg, _ := json.Marshal(map[string]interface{}{"id": message.Id, "error": err.Error()})
	client.Send <- msg
}
