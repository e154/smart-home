package core

type MapBind struct {
	_map		*Map
}

func (m *MapBind) SetDeviceState(devId int64, _state string) {

	m._map.SetDeviceState(devId, _state)
}

func (m *MapBind) GetDeviceState(id int64) interface{} {

	state, ok := m._map.GetDevicesStates()[id]
	if ok {
		return state.SystemName
	} else {
		return nil
	}
}
