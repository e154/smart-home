// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

// Entity ...
type Entity struct {
	Total    int64                  `json:"total"`
	Entities map[string]EntityState `json:"entities"`
}

// EntityManager ...
type EntityManager struct {
	publisher  IPublisher
	total      metrics.Counter
	updateLock sync.Mutex
	entities   map[string]EntityState
}

// NewEntityManager ...
func NewEntityManager(publisher IPublisher) *EntityManager {
	return &EntityManager{
		total:     metrics.NewCounter(),
		publisher: publisher,
		entities:  make(map[string]EntityState),
	}
}

func (d *EntityManager) update(t interface{}) {
	var cursor interface{}
	cursor = "entity"
	switch v := t.(type) {
	case EntityAdd:
		d.total.Inc(v.Num)
	case EntityDelete:
		d.total.Dec(v.Num)
	case EntitySetState:
		if d.updateStatus(v) {
			return
		}
		cursor = EntityCursor{
			DeviceId:    v.DeviceId,
			ElementName: v.ElementName,
		}
	case EntitySetOption:
		if d.updateOption(v) {
			return
		}
		cursor = EntityCursor{
			DeviceId:    v.DeviceId,
			ElementName: v.ElementName,
		}
	default:
		return
	}

	d.broadcast(cursor)
}

func (d *EntityManager) updateOption(v EntitySetOption) (exist bool) {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	key := d.key(v.DeviceId, v.ElementName)
	if element, ok := d.entities[key]; ok {
		if fmt.Sprintf("%v", element.StateOptions) == fmt.Sprintf("%v", v.StateOptions) {
			exist = true
			return
		}

		d.entities[key] = EntityState{
			DeviceId:     v.DeviceId,
			ElementName:  v.ElementName,
			StateId:      v.StateId,
			StateOptions: v.StateOptions,
		}
		return
	}

	d.entities[key] = EntityState{
		DeviceId:     v.DeviceId,
		ElementName:  v.ElementName,
		StateId:      v.StateId,
		StateOptions: v.StateOptions,
	}

	return
}

func (d *EntityManager) updateStatus(v EntitySetState) (exist bool) {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	key := d.key(v.DeviceId, v.ElementName)
	if entity, ok := d.entities[key]; ok {
		if entity.StateId == v.StateId {
			exist = true
			return
		}

		d.entities[key] = EntityState{
			DeviceId:     v.DeviceId,
			ElementName:  v.ElementName,
			StateId:      v.StateId,
			StateOptions: nil,
		}
		return
	}

	d.entities[key] = EntityState{
		DeviceId:     v.DeviceId,
		ElementName:  v.ElementName,
		StateId:      v.StateId,
		StateOptions: nil,
	}

	return
}

func (d *EntityManager) broadcast(cursor interface{}) {
	go d.publisher.Broadcast(cursor)
}

// Snapshot ...
func (d *EntityManager) Snapshot() Entity {
	elements := make(map[string]EntityState)

	var total int64
	d.updateLock.Lock()
	for k, v := range d.entities {
		elements[k] = v
		total++
	}
	d.updateLock.Unlock()

	return Entity{
		Total:    total,
		Entities: elements,
	}
}

func (d *EntityManager) key(deviceId int64, elementName string) string {
	return fmt.Sprintf("%d_%s", deviceId, elementName)
}

// EntityAdd ...
type EntityAdd struct {
	Num int64
}

// EntityDelete ...
type EntityDelete struct {
	Num int64
}

// EntityState ...
type EntityState struct {
	DeviceId     int64       `json:"device_id"`
	ElementName  string      `json:"element_name"`
	StateId      int64       `json:"state_id"`
	StateOptions interface{} `json:"state_options"`
}

// EntitySetState ...
type EntitySetState struct {
	DeviceId    int64  `json:"device_id"`
	ElementName string `json:"element_name"`
	StateId     int64  `json:"state_id"`
}

// EntitySetOption ...
type EntitySetOption struct {
	DeviceId     int64       `json:"device_id"`
	ElementName  string      `json:"element_name"`
	StateId      int64       `json:"state_id"`
	StateOptions interface{} `json:"state_options"`
}

// EntityCursor ...
type EntityCursor struct {
	DeviceId    int64
	ElementName string
}
