package core

import "github.com/e154/smart-home/api/models"

type MapBind struct {
	_map		*Map
}

func (m *MapBind) SetDeviceState(dev *models.Device, _state string) {

	m._map.SetDeviceState(dev, _state)
}

func (m *MapBind) GetDeviceState(dev *models.Device) interface{} {

	state, ok := m._map.GetDevicesStates()[dev.Id]
	if ok {
		return state.SystemName
	} else {
		return nil
	}
}
