package dashboard_models

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/stream"
	"github.com/op/go-logging"
	"reflect"
)

var (
	log = logging.MustGetLogger("models")
)

type Devices struct {
	Total    int64                    `json:"total"`
	Status   map[int64]*m.DeviceState `json:"status"`
	adaptors *adaptors.Adaptors       `json:"-"`
	core     *core.Core               `json:"-"`
}

func NewDevices(adaptors *adaptors.Adaptors,
	core *core.Core) *Devices {
	return &Devices{
		Status:   make(map[int64]*m.DeviceState),
		adaptors: adaptors,
		core:     core,
	}
}

func (d *Devices) Update() {

	var err error
	if _, d.Total, err = d.adaptors.Device.List(999, 0, "", ""); err != nil {
		log.Error(err.Error())
	}

	d.Status = d.core.Map.GetDevicesStates()
}

func (d *Devices) Broadcast() (map[string]interface{}, bool) {

	d.Update()

	return map[string]interface{}{
		"devices": d,
	}, true
}

func (d *Devices) BroadcastOne(deviceId int64, elementName string) (map[string]interface{}, bool) {

	d.Status = d.core.Map.GetDevicesStates()
	state, ok := d.Status[deviceId]
	if !ok {
		return nil, false
	}

	return map[string]interface{}{
		"device":  map[string]interface{}{"id": deviceId, "state": state},
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
