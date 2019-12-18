package core

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
	element = &MapElementBind{e.Map.GetElement(device.model, elementName)}
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
