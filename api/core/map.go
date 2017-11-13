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

// ------------------------------------------------
// Device states
// ------------------------------------------------
func (b *Map) SetDeviceState(devId int64, _state string) {

	//TODO add cache
	states, err := models.GetAllDeviceStateByDevice(devId)
	if err != nil {
		return
	}

	b.Lock()
	defer b.Unlock()

	for _, state := range states {
		if state.SystemName == _state {
			//..
			if old_state, ok := b.deviceStates[devId]; ok {
				if old_state.Id == state.Id {
					return
				}
			}

			b.deviceStates[devId] = state
			//..
			break
		}
	}


	b.telemetry.BroadcastOne("devices", devId)
}

func (b *Map) GetDevicesStates() map[int64]*models.DeviceState {
	b.Lock()
	defer b.Unlock()
	return b.deviceStates
}
