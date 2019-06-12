package telemetry_map

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/stream"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("telemetry_map")
)

type Map struct {
	devices *Devices
	stream  *stream.StreamService
}

func NewMap(stream *stream.StreamService,
	adaptors *adaptors.Adaptors) *Map {

	_map := &Map{
		devices: NewDevices(adaptors),
		stream:  stream,
	}

	return _map
}

func (m *Map) Run(core *core.Core) {
	m.devices.core = core

	m.stream.Subscribe("map.get.devices.states", m.devices.streamGetDevicesStates)
	m.stream.Subscribe("map.get.telemetry", m.streamTelemetry)
}

func (m *Map) Stop() {
	m.stream.UnSubscribe("map.get.devices.states")
	m.stream.UnSubscribe("map.get.telemetry")
}

func (m *Map) BroadcastOne(pack string, deviceId int64, elementName string) {

	var body interface{}
	var ok bool

	switch pack {
	case "devices":
		body, ok = m.devices.BroadcastOne(deviceId, elementName)
	}

	if ok {
		m.sendMsg(body)
	}
}

func (t *Map) Broadcast(pack string) {

	var body interface{}
	var ok bool

	switch pack {
	case "devices":
		body, ok = t.devices.Broadcast()
	}

	if (ok) {
		t.sendMsg(body)
	}
}

func (t *Map) sendMsg(body interface{}) {

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "broadcast",
		"value": map[string]interface{}{
			"type": "map.telemetry",
			"body": body,
		},
	})

	t.stream.Broadcast(msg)
}

// only on request: 'map.get.telemetry'
//
func (m *Map) streamTelemetry(client *stream.Client, value interface{}) {

}
