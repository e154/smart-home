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

// Javascript Binding
//
// Map
//	.SetElementState(device, elementName, newState)
//	.GetElement(device, elementName) -> MapElementBind
//	.GetElements(device) -> []MapElementBind
//
type MapBind struct {
	Map *Map
}

func (e *MapBind) SetElementState(device *DeviceBind, elementName, newState string) {
	if device == nil {
		log.Error("device is nil")
		return
	}

	if device.model == nil {
		log.Error("device.Model is nil")
		return
	}

	e.Map.SetElementState(device.model, elementName, newState)
}

func (e *MapBind) GetElement(device *DeviceBind, elementName string) (element *MapElementBind) {
	if device == nil {
		log.Error("device is nil")
		return
	}

	if device.model == nil {
		log.Error("device.Model is nil")
		return
	}

	mapElement, err := e.Map.GetElement(device.model, elementName)
	if device.model == nil {
		log.Error(err.Error())
		return
	}

	element = &MapElementBind{mapElement}
	return
}

func (e *MapBind) GetElements(device *DeviceBind) (elements []*MapElementBind) {

	elements = make([]*MapElementBind, 0)

	if device == nil {
		return
	}

	mapElements := e.Map.GetElements(device.model)
	for _, mapElement := range mapElements {
		elements = append(elements, &MapElementBind{mapElement})
	}

	return
}
