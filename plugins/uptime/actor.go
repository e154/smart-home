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

package uptime

import (
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/shirou/gopsutil/host"
	"go.uber.org/atomic"
	"sync"
	"time"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	appStarted time.Time
	total      *atomic.Uint64
	eventBus   event_bus.EventBus
}

// NewActor ...
func NewActor(entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus) *Actor {
	return &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:                common.EntityId(fmt.Sprintf("%s.%s", EntitySensor, Name)),
			Name:              Name,
			EntityType:        EntitySensor,
			UnitOfMeasurement: "days",
			AttrMu:            &sync.RWMutex{},
			Attrs:             NewAttr(),
			Manager:           entityManager,
		},
		eventBus:   eventBus,
		appStarted: time.Now(),
		total:      atomic.NewUint64(0),
	}
}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {
	return e
}

func (e *Actor) update() {

	oldState := e.GetEventState(e)

	e.Now(oldState)

	total, err := host.Uptime()
	if err != nil {
		return
	}

	e.total.Store(total)

	e.AttrMu.Lock()
	e.Attrs[AttrUptimeTotal].Value = e.total.Load()
	e.Attrs[AttrUptimeAppStarted].Value = e.appStarted
	e.AttrMu.Unlock()

	e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventStateChanged{
		Type:     e.Id.Type(),
		EntityId: e.Id,
		OldState: oldState,
		NewState: e.GetEventState(e),
	})
}
