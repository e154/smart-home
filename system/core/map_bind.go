package core

import (
	m "github.com/e154/smart-home/models"
)

// Javascript Binding
//
// map
//	.setElementState(device, elementName, newState)
//	.getElement(device, elementName) -> MapElementBind
//	.getElements(device) -> []MapElementBind
//
type MapBind struct {
	Map *Map
}

func (e *MapBind) SetElementState(device *m.Device, elementName, newState string) {
	if device == nil {
		log.Error("device is nil")
		return
	}
	e.Map.SetElementState(device, elementName, newState)
}

func (e *MapBind) GetElement(device *m.Device, elementName string) (element *MapElementBind) {
	element = &MapElementBind{e.Map.GetElement(device, elementName)}
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
