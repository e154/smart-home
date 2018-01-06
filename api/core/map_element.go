package core

import (
	"github.com/e154/smart-home/api/models"
	"fmt"
)

type MapElement struct {
	Map		*Map
	Options interface{}
	Device  *models.Device
	State   *models.DeviceState
}

func (e *MapElement) SetState(systemName string) {

	for _, state := range e.Device.States {
		if state.SystemName != systemName {
			continue
		}

		if e.State != nil && e.State.SystemName == state.SystemName {
			return
		}

		e.State = state

		e.Map.telemetry.BroadcastOne("devices", e.Device.Id)
	}
}

func (e *MapElement) GetState() interface{} {

	return e.State
}

func (e *MapElement) SetOptions(options interface{}) {

	if fmt.Sprint(e.Options) == fmt.Sprint(options) {
		return
	}

	e.Options = options

	e.Map.telemetry.BroadcastOne("devices", e.Device.Id)
}

func (e *MapElement) GetOptions() interface{} {

	return e.Options
}

