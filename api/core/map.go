package core

import (
	"github.com/e154/smart-home/api/models"
	"sync"
)

type Map struct {
	telemetry Telemetry
	elements  sync.Map
}

func (b *Map) SetElementState(device *models.Device, systemName string) {

	if device == nil {
		return
	}

	if v, ok := b.elements.Load(device.Id); ok {
		v.(*MapElement).SetState(systemName)
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

	b.elements.Range(func(key, value interface{}) bool {
		element := value.(*MapElement)
		states[element.Device.Id] = element.State
		return true
	})

	return
}

func (b *Map) GetAllElements() map[int64]*MapElement {

	elements := make(map[int64]*MapElement)

	b.elements.Range(func(key, value interface{}) bool {
		element := value.(*MapElement)
		elements[element.Device.Id] = element
		return true
	})

	return elements
}

func (b *Map) GetElement(device *models.Device) *MapElement {

	if device == nil {
		return nil
	}

	if v, ok := b.elements.Load(device.Id); ok {
		return v.(*MapElement)
	}

	return b.NewMapElement(device, nil)

	return nil
}

func (b *Map) NewMapElement(device *models.Device, state *models.DeviceState) *MapElement {

	element := &MapElement{
		Map:     b,
		Device:  device,
		State:   state,
		Options: nil,
	}

	b.elements.Store(device.Id, element)

	return element
}
