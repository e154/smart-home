// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package memory_app

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	eventBus   bus.Bus
	updateLock *sync.Mutex
}

// NewActor ...
func NewActor(visor supervisor.Supervisor,
	eventBus bus.Bus,
	entity *m.Entity) *Actor {

	actor := &Actor{
		BaseActor: supervisor.BaseActor{
			Id:                common.EntityId(fmt.Sprintf("%s.%s", EntityMemory, Name)),
			Name:              Name,
			EntityType:        EntityMemory,
			UnitOfMeasurement: "GHz",
			AttrMu:            &sync.RWMutex{},
			Attrs:             NewAttr(),
			Supervisor:        visor,
		},
		eventBus:   eventBus,
		updateLock: &sync.Mutex{},
	}

	if entity != nil {
		actor.Metric = entity.Metrics
	}

	return actor
}

// Spawn ...
func (e *Actor) Spawn() supervisor.PluginActor {
	return e
}

func (u *Actor) selfUpdate() {

	u.updateLock.Lock()
	defer u.updateLock.Unlock()

	oldState := u.GetEventState(u)
	u.Now(oldState)

	var s runtime.MemStats
	runtime.ReadMemStats(&s)

	u.AttrMu.Lock()
	u.Attrs[AttrAlloc].Value = s.Alloc
	u.Attrs[AttrHeapAlloc].Value = s.HeapAlloc
	u.Attrs[AttrTotalAlloc].Value = s.TotalAlloc
	u.Attrs[AttrSys].Value = s.Sys
	u.Attrs[AttrNumGC].Value = s.NumGC
	u.Attrs[AttrLastGC].Value = time.Unix(0, int64(s.LastGC))
	u.AttrMu.Unlock()

	//u.SetMetric(u.Id, "memory_app", map[string]float32{
	//	"total_alloc": float32(s.TotalAlloc),
	//})

	u.eventBus.Publish(bus.TopicEntities, events.EventStateChanged{
		StorageSave: false,
		PluginName:  u.Id.PluginName(),
		EntityId:    u.Id,
		OldState:    oldState,
		NewState:    u.GetEventState(u),
	})
}
