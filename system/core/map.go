// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package core

import (
	"fmt"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/telemetry"
	"strings"
	"sync"
)

type Map struct {
	telemetry telemetry.ITelemetry
	elements  sync.Map
}

func (b *Map) SetElementState(device *m.Device, elementName, systemName string) {

	if device == nil || elementName == "" || systemName == "" {
		return
	}

	hashKey := b.key(device, elementName)

	if v, ok := b.elements.Load(hashKey); ok {
		v.(*MapElement).SetState(systemName)
	} else {
		for _, state := range device.States {
			if state.SystemName != systemName {
				continue
			}
			b.NewMapElement(device, elementName, state)
			b.telemetry.BroadcastOne(telemetry.Device{Id: device.Id, ElementName: elementName})
		}
	}
}

func (b *Map) GetDevicesStates() (states map[int64]*m.DashboardDeviceStatus) {
	states = make(map[int64]*m.DashboardDeviceStatus)

	b.elements.Range(func(key, value interface{}) bool {
		element := value.(*MapElement)
		if element.State != nil {
			states[element.Device.Id] = &m.DashboardDeviceStatus{
				Id:          element.State.Id,
				Description: element.State.Description,
				SystemName:  element.State.SystemName,
				DeviceId:    element.State.DeviceId,
			}
		}
		return true
	})

	return
}

func (b *Map) GetElement(device *m.Device, elementName string) (element *MapElement) {

	if device == nil || elementName == "" {
		return
	}

	hashKey := b.key(device, elementName)

	if v, ok := b.elements.Load(hashKey); ok {
		element = v.(*MapElement)
	} else {
		element = b.NewMapElement(device, elementName, nil)
	}
	return
}

func (b *Map) GetAllElements() (elements []*MapElement) {

	b.elements.Range(func(key, value interface{}) bool {
		element := value.(*MapElement)
		elements = append(elements, element)
		return true
	})
	return
}

func (b *Map) GetElements(device *m.Device) (elements []*MapElement) {

	if device == nil {
		return nil
	}

	elements = make([]*MapElement, 0)

	partKeyName := fmt.Sprintf("device(%d)_elementName", device.Id)

	b.elements.Range(func(key, value interface{}) bool {
		if strings.Contains(key.(string), partKeyName) {
			element := value.(*MapElement)
			elements = append(elements, element)
		}

		return true
	})

	return
}

func (b *Map) NewMapElement(device *m.Device, elementName string, state *m.DeviceState) *MapElement {

	element := &MapElement{
		Map:         b,
		Device:      device,
		State:       state,
		Options:     nil,
		ElementName: elementName,
	}

	hashKey := b.key(device, elementName)

	b.elements.Store(hashKey, element)

	return element
}

func (b *Map) key(device *m.Device, elementName string) string {
	return fmt.Sprintf("device(%d)_elementName(%s)", device.Id, elementName)
}
