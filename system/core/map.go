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
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/metrics"
	"strings"
	"sync"
)

type Map struct {
	metric   *metrics.MetricManager
	elements sync.Map
	adaptors *adaptors.Adaptors
}

func NewMap(metric *metrics.MetricManager,
	adaptors *adaptors.Adaptors) *Map {
	return &Map{
		metric:   metric,
		adaptors: adaptors,
	}
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
			if _, err := b.NewMapElement(device, elementName, state); err != nil {
				log.Error(err.Error())
				return
			}

			b.metric.Update(metrics.MapElementSetState{
				StateId:     state.Id,
				DeviceId:    state.DeviceId,
				ElementName: elementName,
			})
		}
	}
}

func (b *Map) GetElement(device *m.Device, elementName string) (element *MapElement, err error) {

	if device == nil || elementName == "" {
		err = fmt.Errorf("bad parameters")
		return
	}

	hashKey := b.key(device, elementName)

	if v, ok := b.elements.Load(hashKey); ok {
		element = v.(*MapElement)
	} else {
		element, err = b.NewMapElement(device, elementName, nil)
	}
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

func (b *Map) NewMapElement(device *m.Device, elementName string, state *m.DeviceState) (*MapElement, error) {

	element, err := NewMapElement(device, elementName, state, b, b.adaptors)
	if err != nil {
		return nil, err
	}

	hashKey := b.key(device, elementName)

	b.elements.Store(hashKey, element)

	b.metric.Update(metrics.MapElementAdd{Num: 1})

	return element, nil
}

func (b *Map) key(device *m.Device, elementName string) string {
	return fmt.Sprintf("device(%d)_elementName(%s)", device.Id, elementName)
}
