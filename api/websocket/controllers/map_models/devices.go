package map_models

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/stream"
	"sync"
)

type DeviceState struct {
	Id          int64          `json:"id"`
	Status      *m.DeviceState `json:"status"`
	Options     interface{}    `json:"options"`
	ElementName string         `json:"element_name"`
}

type Devices struct {
	sync.Mutex
	Total       int64                   `json:"total"`
	DeviceStats map[string]*DeviceState `json:"device_stats"`
	adaptors    *adaptors.Adaptors
	core        *core.Core
}

func NewDevices(adaptors *adaptors.Adaptors,
	core *core.Core) *Devices {
	return &Devices{
		DeviceStats: make(map[string]*DeviceState),
		adaptors:    adaptors,
		core:        core,
	}
}

func (d *Devices) Update() {

	if d.core == nil {
		return
	}

	mapElements := d.core.Map.GetAllElements()
	d.Total = int64(len(mapElements))

	d.Lock()
	d.DeviceStats = make(map[string]*DeviceState)
	for _, mapElement := range mapElements {
		deviceId := mapElement.Device.Id
		key := d.key(deviceId, mapElement.ElementName)
		d.DeviceStats[key] = &DeviceState{
			Id:          deviceId,
			Status:      mapElement.State,
			Options:     mapElement.Options,
			ElementName: mapElement.ElementName,
		}
	}
	d.Unlock()
}

func (d *Devices) Broadcast() (map[string]interface{}, bool) {

	d.Update()

	return map[string]interface{}{
		"devices": d,
	}, true
}

func (d *Devices) BroadcastOne(deviceId int64, elementName string) (map[string]interface{}, bool) {

	d.Update()

	d.Lock()
	key := d.key(deviceId, elementName)
	state, ok := d.DeviceStats[key]
	if !ok {
		d.Unlock()
		return nil, false
	}
	d.Unlock()

	return map[string]interface{}{
		"device": state,
	}, true
}

// only on request: 'map.get.devices.states'
//
func (d *Devices) GetDevicesStates(client *stream.Client, message stream.Message) {

	d.Update()

	msg := &stream.Message{
		Id:      message.Id,
		Forward: stream.Response,
		Payload: map[string]interface{}{
			"states": d,
		},
	}

	client.Send <- msg.Pack()
}

func (d *Devices) key(deviceId int64, elementName string) string {
	return fmt.Sprintf("%d_%s", deviceId, elementName)
}
