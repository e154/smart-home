package controllers

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/stream"
)

type ControllerCommon struct {
	adaptors *adaptors.Adaptors
	stream   *stream.StreamService
	endpoint *endpoint.Endpoint
}

func NewControllerCommon(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	endpoint *endpoint.Endpoint) *ControllerCommon {
	return &ControllerCommon{
		adaptors: adaptors,
		endpoint: endpoint,
		stream:   stream,
	}
}

func (c *ControllerCommon) Err(client *stream.Client, message stream.Message, err error) {
	msg, _ := json.Marshal(map[string]interface{}{"id": message.Id, "error": err.Error()})
	client.Send <- msg
}
