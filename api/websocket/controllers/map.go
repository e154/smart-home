package controllers

import (
	"encoding/json"
	mapModels "github.com/e154/smart-home/api/websocket/controllers/map_models"
	"github.com/e154/smart-home/system/stream"
)

type ControllerMap struct {
	*ControllerCommon
	devices *mapModels.Devices
}

func NewControllerMap(common *ControllerCommon) *ControllerMap {
	return &ControllerMap{
		ControllerCommon: common,
		devices:          mapModels.NewDevices(common.adaptors, common.core),
	}
}

func (c *ControllerMap) Start() {
	c.telemetry.Subscribe("map", c)
	c.stream.Subscribe("map.get.devices.states", c.devices.GetDevicesStates)
	c.stream.Subscribe("map.get.telemetry", c.streamTelemetry)
}

func (c *ControllerMap) Stop() {
	c.telemetry.UnSubscribe("map")
	c.stream.UnSubscribe("map.get.devices.states")
	c.stream.UnSubscribe("map.get.telemetry")
}

func (c *ControllerMap) BroadcastOne(pack string, deviceId int64, elementName string) {
	var body interface{}
	var ok bool

	switch pack {
	case "devices":
		body, ok = c.devices.BroadcastOne(deviceId, elementName)
	}

	if ok {
		c.sendMsg(body)
	}
}

func (c *ControllerMap) Broadcast(pack string) {

	var body interface{}
	var ok bool

	switch pack {
	case "devices":
		body, ok = c.devices.Broadcast()
	}

	if ok {
		c.sendMsg(body)
	}
}

func (t *ControllerMap) sendMsg(body interface{}) {

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "broadcast",
		"value": map[string]interface{}{
			"type": "map.telemetry",
			"body": body,
		},
	})

	t.stream.Broadcast(msg)
}

func (m *ControllerMap) streamTelemetry(client *stream.Client, value interface{}) {

}
