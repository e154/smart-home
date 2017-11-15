package core

import (
	"github.com/e154/smart-home/api/models"
	"sync"
)

//TODO refactor map system
type Map struct {
	telemetry		Telemetry
	sync.Mutex
	deviceStates	map[int64]*models.DeviceState
}

func (b *Map) SetDeviceState(dev *models.Device, _state string) {

	b.Lock()
	defer b.Unlock()

	for _, state := range dev.States {
		if state.SystemName == _state {
			//..
			if old_state, ok := b.deviceStates[dev.Id]; ok {
				if old_state.Id == state.Id {
					return
				}
			}

			b.deviceStates[dev.Id] = state
			//..
			break
		}
	}


	b.telemetry.BroadcastOne("devices", dev.Id)
}

func (b *Map) GetDevicesStates() map[int64]*models.DeviceState {
	b.Lock()
	defer b.Unlock()
	return b.deviceStates
}
