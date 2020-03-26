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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/metrics"
	"reflect"
	"sync"
	"time"
)

type MapElement struct {
	Map         *Map
	Options     interface{}
	State       *m.DeviceState
	adaptors    *adaptors.Adaptors
	elementLock *sync.Mutex
	mapElement  *m.MapElement
}

func NewMapElement(elementName string,
	systemName *string,
	_map *Map,
	adaptors *adaptors.Adaptors) (*MapElement, error) {

	mapElement, err := adaptors.MapElement.GetByName(elementName)
	if err != nil {
		return nil, err
	}

	var state *m.DeviceState

	if systemName != nil {
		for _, _state := range mapElement.Prototype.Device.States {
			if _state.SystemName != common.StringValue(systemName) {
				continue
			}
			state = _state
		}
	}

	return &MapElement{
		Map:         _map,
		State:       state,
		Options:     nil,
		mapElement:  mapElement,
		adaptors:    adaptors,
		elementLock: &sync.Mutex{},
	}, nil
}

func (e *MapElement) SetState(systemName string) {
	e.elementLock.Lock()
	defer e.elementLock.Unlock()

	if e.State != nil && e.State.SystemName == systemName {
		return
	}

	for _, state := range e.mapElement.Prototype.States {
		if state.DeviceState.SystemName != systemName {
			continue
		}

		if e.State != nil && e.State.Id == state.DeviceStateId {
			return
		}

		e.State = state.DeviceState

		go e.updateDeviceHistory(state.DeviceState)

		go e.Map.metric.Update(metrics.MapElementSetState{
			DeviceId:    e.State.DeviceId,
			ElementName: e.mapElement.Name,
			StateId:     e.State.Id,
		})
	}
}

func (e *MapElement) GetState() interface{} {
	e.elementLock.Lock()
	defer e.elementLock.Unlock()

	return e.State
}

func (e *MapElement) SetOptions(options interface{}) {
	e.elementLock.Lock()
	defer e.elementLock.Unlock()

	if fmt.Sprint(e.Options) == fmt.Sprint(options) {
		return
	}

	e.Options = options

	if e.State == nil {
		return
	}

	go e.updateDeviceHistory(options)

	e.Map.metric.Update(metrics.MapElementSetOption{
		StateId:      e.State.Id,
		DeviceId:     e.mapElement.Prototype.DeviceId,
		ElementName:  e.mapElement.Name,
		StateOptions: options,
	})
}

func (e *MapElement) GetOptions() interface{} {
	e.elementLock.Lock()
	defer e.elementLock.Unlock()

	return e.Options
}

func (e *MapElement) updateDeviceHistory(t interface{}) {
	e.elementLock.Lock()
	defer e.elementLock.Unlock()

	switch e.mapElement.PrototypeType {
	case common.PrototypeTypeDevice:
	default:
		return
	}

	var historyType = common.MapDeviceHistoryState
	var description string
	switch v := t.(type) {
	case *m.DeviceState:
		description = v.Description
	case m.DeviceState:
		description = v.Description
	case string:
		description = v
		historyType = common.MapDeviceHistoryOption
	case interface{}:
		description = fmt.Sprintf("%v", v)
		historyType = common.MapDeviceHistoryOption
	default:
		log.Warnf("unknown object type %v", reflect.TypeOf(t).String())
	}

	var logLevel = common.LogLevelInfo

	_, err := e.adaptors.MapDeviceHistory.Add(m.MapDeviceHistory{
		MapElementId: e.mapElement.Id,
		MapDeviceId:  e.mapElement.PrototypeId,
		LogLevel:     common.LogLevelInfo,
		Type:         historyType,
		Description:  description,
	})
	if err != nil {
		log.Error(err.Error())
	}

	e.Map.metric.Update(metrics.HistoryItem{
		DeviceName:        e.mapElement.Name,
		DeviceDescription: e.mapElement.Description,
		Type:              string(historyType),
		LogLevel:          string(logLevel),
		Description:       description,
		CreatedAt:         time.Now(),
	})
}

func (e *MapElement) CustomHistory(logLevel, t, desc string) {
	e.elementLock.Lock()
	defer e.elementLock.Unlock()

	e.adaptors.MapDeviceHistory.Add(m.MapDeviceHistory{
		MapElementId: e.mapElement.Id,
		MapDeviceId:  e.mapElement.PrototypeId,
		LogLevel:     common.LogLevel(logLevel),
		Type:         common.MapDeviceHistoryType(t),
		Description:  desc,
	})

	e.Map.metric.Update(metrics.HistoryItem{
		DeviceName:        e.mapElement.Name,
		DeviceDescription: e.mapElement.Description,
		Type:              string(t),
		LogLevel:          string(logLevel),
		Description:       desc,
		CreatedAt:         time.Now(),
	})
}
