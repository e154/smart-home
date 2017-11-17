package core

import (
	"github.com/e154/smart-home/api/models"
)

type MapBind struct {
	Map *Map
}

func (e *MapBind) SetElementState(device *models.Device, state string) {

	e.Map.SetElementState(device, state)
}

func (e *MapBind) GetElement(device *models.Device) *MapElement {
	return e.Map.GetElement(device)
}
