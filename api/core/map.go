package core

import (
	"github.com/e154/smart-home/api/models"
	"sync"
)

//TODO refactor map system
type Map struct {
	telemetry Telemetry
	sync.Mutex
	elements  map[int64]*MapElement
}

func (b *Map) SetElementState(device *models.Device, systemName string) {

	if device == nil {
		return
	}

	b.Lock()
	defer b.Unlock()

	if _, ok := b.elements[device.Id]; ok {
		b.elements[device.Id].SetState(systemName)
	} else {
		for _, state := range device.States {
			if state.SystemName != systemName {
				continue
			}

			b.NewMapElement(device, state)

			b.telemetry.BroadcastOne("devices", device.Id)
		}
	}
}

func (b *Map) GetDevicesStates() (states map[int64]*models.DeviceState) {
	states = make(map[int64]*models.DeviceState)

	for _, element := range b.elements {
		states[element.Device.Id] = element.State
	}

	return
}

func (b *Map) GetAllElements() map[int64]*MapElement {
	b.Lock()
	defer b.Unlock()

	return b.elements
}

func (b *Map) GetElement(device *models.Device) *MapElement {

	if device == nil {
		return nil
	}

	if element, ok := b.elements[device.Id]; ok {
		return element
	}

	return b.NewMapElement(device, nil)

	return nil
}

func (b *Map) NewMapElement(device *models.Device, state *models.DeviceState) *MapElement {

	b.Lock()
	b.elements[device.Id] = &MapElement{
		Map: b,
		Device: device,
		State: state,
		Options: nil,
	}
	b.Unlock()

	return b.elements[device.Id]
}