package telemetry_map

import (
	"github.com/e154/smart-home/api/stream"
	"encoding/json"
)

var (
	Hub		stream.Hub
)

func NewMap() *Map {

	Hub	= stream.GetHub()

	_map := &Map{
		Devices: &Devices{
			DeviceStats: make(map[int64]*DeviceState),
		},
	}

	Hub.Subscribe("map.get.devices.states", _map.Devices.streamGetDevicesStates)
	Hub.Subscribe("map.get.telemetry", _map.streamTelemetry)

	return _map
}

type Map struct {
	Devices *Devices
}

func (m *Map) Run() {

}

func (m *Map) Stop() {

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

func (t *Map) sendMsg (body interface{}) {

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "broadcast",
		"value": map[string]interface{}{
			"type": "map.telemetry",
			"body": body,
		},
	})

	Hub.Broadcast(msg)
}

// only on request: 'map.get.telemetry'
//
func (m *Map) streamTelemetry (client *stream.Client, value interface{}) {


}