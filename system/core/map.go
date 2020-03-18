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
	"github.com/e154/smart-home/system/metrics"
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

func (b *Map) SetElementState(elementName, stateSystemName string) {

	if elementName == "" || stateSystemName == "" {
		return
	}

	hashKey := b.key(elementName)

	if v, ok := b.elements.Load(hashKey); ok {
		v.(*MapElement).SetState(stateSystemName)
	} else {
		mapElement, err := b.NewMapElement(elementName, common.String(stateSystemName))
		if err != nil {
			log.Error(err.Error())
			return
		}

		if mapElement.State == nil {
			return
		}
		b.metric.Update(metrics.MapElementSetState{
			StateId:     mapElement.State.Id,
			DeviceId:    mapElement.State.DeviceId,
			ElementName: elementName,
		})
	}
}

func (b *Map) GetElement(elementName string) (element *MapElement, err error) {

	if elementName == "" {
		err = fmt.Errorf("bad parameters")
		return
	}

	hashKey := b.key(elementName)

	if v, ok := b.elements.Load(hashKey); ok {
		element = v.(*MapElement)
	} else {
		element, err = b.NewMapElement(elementName, nil)
	}

	return
}

func (b *Map) NewMapElement(elementName string, stateSystemName *string) (*MapElement, error) {

	element, err := NewMapElement(elementName, stateSystemName, b, b.adaptors)
	if err != nil {
		return nil, err
	}

	hashKey := b.key(elementName)

	b.elements.Store(hashKey, element)

	b.metric.Update(metrics.MapElementAdd{Num: 1})

	return element, nil
}

func (b *Map) key(elementName string) string {
	return fmt.Sprintf("%s", elementName)
}
