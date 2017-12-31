package telemetry_map

import (
	"encoding/json"
	"reflect"
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/stream"
	"github.com/e154/smart-home/api/core"
	"github.com/e154/smart-home/api/log"
	"sync"
)

type DeviceState struct {
	Id      int64                              `json:"id"`
	Status  *models.DeviceState                `json:"status"`
	Options interface{}                        `json:"options"`
}

type Devices struct {
	sync.Mutex
	Total       int64                          `json:"total"`
	DeviceStats map[int64]*DeviceState         `json:"device_stats"`
}

func (d *Devices) Update() {
	var err error
	if d.Total, err = models.GetDevicesCount(); err != nil {
		log.Error(err.Error())
	}

	mapElemets := core.CorePtr().Map.GetAllElements()

	d.Lock()
	defer d.Unlock()
	d.DeviceStats = make(map[int64]*DeviceState)
	for id, mapElement := range mapElemets {
		d.DeviceStats[id] = &DeviceState{
			Id: id,
			Status: mapElement.State,
			Options: mapElement.Options,
		}
	}
}

func (d *Devices) Broadcast() (interface{}, bool) {

	d.Update()

	return map[string]interface{}{
		"devices": d,
	}, true
}

func (d *Devices) BroadcastOne(id int64) (interface{}, bool) {

	d.Update()

	d.Lock()
	state, ok := d.DeviceStats[id]
	if !ok {
		return nil, false
	}
	d.Unlock()

	return map[string]interface{}{
		"device": state,
	}, true
}

// only on request: 'map.get.devices.states'
//
func (d *Devices) streamGetDevicesStates(client *stream.Client, value interface{}) {

	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	d.Update()

	msg, _ := json.Marshal(map[string]interface{}{
		"id": v["id"],
		"states": d,
	})

	client.Send <- msg
}