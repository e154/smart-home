package telemetry

import (
	"encoding/json"
	"reflect"
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/stream"
	"github.com/e154/smart-home/api/core"
	"github.com/e154/smart-home/api/log"
)

type Devices struct {
	Total  int64				`json:"total"`
	Status map[int64]*models.DeviceState	`json:"status"`
}

func (d *Devices) Update() {
	var err error
	if d.Total, err = models.GetDevicesCount(); err != nil {
		log.Error(err.Error())
	}

	d.Status = core.CorePtr().Map.GetDevicesStates()
}

func (d *Devices) Broadcast() {

	d.Update()

	msg, _ := json.Marshal(map[string]interface{}{"type": "broadcast",
		"value": map[string]interface{}{"type": "telemetry", "body": map[string]interface{}{
			"devices": d,
		}}},
	)
	Hub.Broadcast(string(msg))
}

func (d *Devices) BroadcastOne(id int64) {

	d.Status = core.CorePtr().Map.GetDevicesStates()
	state, ok := d.Status[id]
	if !ok {
		return
	}

	msg, _ := json.Marshal(
		map[string]interface{}{"type": "broadcast",
			"value": map[string]interface{}{"type": "telemetry", "body": map[string]interface{}{
				"device": map[string]interface{}{"id": id, "state": state},
				"devices": d,
			}}},
	)

	Hub.Broadcast(string(msg))
}

// only on request: 'get.devices.states'
//
func streamGetDevicesStates(client *stream.Client, value interface{}) {

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	Telemetry.Devices.Update()

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "states": Telemetry.Devices.Status})
	client.Send(string(msg))
}