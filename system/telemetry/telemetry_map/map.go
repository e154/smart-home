package telemetry_map

import (
	"encoding/json"
	"github.com/e154/smart-home/system/stream"
	"github.com/op/go-logging"
	"github.com/e154/smart-home/adaptors"
)

var (
	log = logging.MustGetLogger("telemetry_map")
)

func NewMap(stream *stream.StreamService,
	adaptors *adaptors.Adaptors, ) *Map {

	_map := &Map{
		Devices: &Devices{
			DeviceStats: make(map[int64]*DeviceState),
			adaptors:    adaptors,
		},
		stream: stream,
	}

	return _map
}

type Map struct {
	Devices *Devices
	stream  *stream.StreamService
}

func (m *Map) Run() {
	m.stream.Subscribe("map.get.devices.states", m.Devices.streamGetDevicesStates)
	m.stream.Subscribe("map.get.telemetry", m.streamTelemetry)
}

func (m *Map) Stop() {
	m.stream.UnSubscribe("map.get.devices.states")
	m.stream.UnSubscribe("map.get.telemetry")
}

func (m *Map) BroadcastOne(pack string, id int64) {

	var body interface{}
	var ok bool

	switch pack {
	case "devices":
		body, ok = m.Devices.BroadcastOne(id)
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
		body, ok = t.Devices.Broadcast()
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
