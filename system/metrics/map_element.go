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

package metrics

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"sync"
)

// MapElement ...
type MapElement struct {
	Total    int64                      `json:"total"`
	Elements map[string]MapElementState `json:"devices"`
}

// MapElementManager ...
type MapElementManager struct {
	publisher  IPublisher
	total      metrics.Counter
	updateLock sync.Mutex
	elements   map[string]MapElementState
}

// NewMapElementManager ...
func NewMapElementManager(publisher IPublisher) *MapElementManager {
	return &MapElementManager{
		total:     metrics.NewCounter(),
		publisher: publisher,
		elements:  make(map[string]MapElementState),
	}
}

func (d *MapElementManager) update(t interface{}) {
	var cursor interface{}
	cursor = "map_element"
	switch v := t.(type) {
	case MapElementAdd:
		d.total.Inc(v.Num)
	case MapElementDelete:
		d.total.Dec(v.Num)
	case MapElementSetState:
		if d.updateStatus(v) {
			return
		}
		cursor = MapElementCursor{
			DeviceId:    v.DeviceId,
			ElementName: v.ElementName,
		}
	case MapElementSetOption:
		if d.updateOption(v) {
			return
		}
		cursor = MapElementCursor{
			DeviceId:    v.DeviceId,
			ElementName: v.ElementName,
		}
	default:
		return
	}

	d.broadcast(cursor)
}

func (d *MapElementManager) updateOption(v MapElementSetOption) (exist bool) {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	key := d.key(v.DeviceId, v.ElementName)
	if element, ok := d.elements[key]; ok {
		if fmt.Sprintf("%v", element.StateOptions) == fmt.Sprintf("%v", v.StateOptions) {
			exist = true
			return
		}

		d.elements[key] = MapElementState{
			DeviceId:     v.DeviceId,
			ElementName:  v.ElementName,
			StateId:      v.StateId,
			StateOptions: v.StateOptions,
		}
		return
	}

	d.elements[key] = MapElementState{
		DeviceId:     v.DeviceId,
		ElementName:  v.ElementName,
		StateId:      v.StateId,
		StateOptions: v.StateOptions,
	}

	return
}

func (d *MapElementManager) updateStatus(v MapElementSetState) (exist bool) {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	key := d.key(v.DeviceId, v.ElementName)
	if element, ok := d.elements[key]; ok {
		if element.StateId == v.StateId {
			exist = true
			return
		}

		d.elements[key] = MapElementState{
			DeviceId:     v.DeviceId,
			ElementName:  v.ElementName,
			StateId:      v.StateId,
			StateOptions: nil,
		}
		return
	}

	d.elements[key] = MapElementState{
		DeviceId:     v.DeviceId,
		ElementName:  v.ElementName,
		StateId:      v.StateId,
		StateOptions: nil,
	}

	return
}

func (d *MapElementManager) broadcast(cursor interface{}) {
	go d.publisher.Broadcast(cursor)
}

// Snapshot ...
func (d *MapElementManager) Snapshot() MapElement {
	elements := make(map[string]MapElementState)

	var total int64
	d.updateLock.Lock()
	for k, v := range d.elements {
		elements[k] = v
		total++
	}
	d.updateLock.Unlock()

	return MapElement{
		Total:    total,
		Elements: elements,
	}
}

func (d *MapElementManager) key(deviceId int64, elementName string) string {
	return fmt.Sprintf("%d_%s", deviceId, elementName)
}

// MapElementAdd ...
type MapElementAdd struct {
	Num int64
}

// MapElementDelete ...
type MapElementDelete struct {
	Num int64
}

// MapElementState ...
type MapElementState struct {
	DeviceId     int64       `json:"device_id"`
	ElementName  string      `json:"element_name"`
	StateId      int64       `json:"state_id"`
	StateOptions interface{} `json:"state_options"`
}

// MapElementSetState ...
type MapElementSetState struct {
	DeviceId    int64  `json:"device_id"`
	ElementName string `json:"element_name"`
	StateId     int64  `json:"state_id"`
}

// MapElementSetOption ...
type MapElementSetOption struct {
	DeviceId     int64       `json:"device_id"`
	ElementName  string      `json:"element_name"`
	StateId      int64       `json:"state_id"`
	StateOptions interface{} `json:"state_options"`
}

// MapElementCursor ...
type MapElementCursor struct {
	DeviceId    int64
	ElementName string
}
