package telemetry_map

import (
	"encoding/json"
	"reflect"
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/stream"
	"github.com/e154/smart-home/api/core"
	"github.com/e154/smart-home/api/log"
)

type Devices struct {
	Total  int64							`json:"total"`
	Status map[int64]*models.DeviceState	`json:"status"`
}

func (d *Devices) Update() {
	var err error
	if d.Total, err = models.GetDevicesCount(); err != nil {
		log.Error(err.Error())
	}

	d.Status = core.CorePtr().Map.GetDevicesStates()
}

func (d *Devices) Broadcast() (interface{}, bool) {

	d.Update()

	return map[string]interface{}{
		"devices": d,
	}, true
}

func (d *Devices) BroadcastOne(id int64) (interface{}, bool) {

	d.Status = core.CorePtr().Map.GetDevicesStates()
	state, ok := d.Status[id]
	if !ok {
		return nil, false
	}

	return map[string]interface{}{
		"device": map[string]interface{}{"id": id, "state": state},
		"devices": d,
	}, true
}

// only on request: 'dashboard.get.devices.states'
//
func (d *Devices) streamGetDevicesStates(client *stream.Client, value interface{}) {

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	d.Update()

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "states": d.Status})

	client.Send <- msg
}