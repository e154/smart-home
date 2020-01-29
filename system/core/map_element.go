package core

import (
	"fmt"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/telemetry"
	"sync"
)

type MapElement struct {
	sync.Mutex
	Map         *Map
	Options     interface{}
	Device      *m.Device
	State       *m.DeviceState
	ElementName string
}

func (e *MapElement) SetState(systemName string) {
	e.Lock()
	defer e.Unlock()

	for _, state := range e.Device.States {
		if state.SystemName != systemName {
			continue
		}

		if e.State != nil && e.State.SystemName == state.SystemName {
			return
		}

		e.State = state

		e.Map.telemetry.BroadcastOne(telemetry.Device{Id: e.Device.Id, ElementName: e.ElementName})
	}
}

func (e *MapElement) GetState() interface{} {
	e.Lock()
	defer e.Unlock()

	return e.State
}

func (e *MapElement) SetOptions(options interface{}) {
	e.Lock()
	defer e.Unlock()

	if fmt.Sprint(e.Options) == fmt.Sprint(options) {
		return
	}

	e.Options = options

	e.Map.telemetry.BroadcastOne(telemetry.Device{Id: e.Device.Id, ElementName: e.ElementName})
}

func (e *MapElement) GetOptions() interface{} {
	e.Lock()
	defer e.Unlock()

	return e.Options
}
