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

package zone

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
	"sync"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	eventBus event_bus.EventBus
	entities []entity_manager.PluginActor
	stateMu  *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	scriptService scripts.ScriptService,
	adaptors *adaptors.Adaptors,
	eventBus event_bus.EventBus) *Actor {

	actor := &Actor{
		BaseActor: entity_manager.NewBaseActor(entity, scriptService, adaptors),
		eventBus:  eventBus,
		stateMu:   &sync.Mutex{},
	}
	actor.Setts = entity.Settings
	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}
	return actor
}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {
	return e
}

// SetState ...
func (e *Actor) SetState(params entity_manager.EntityStateParams) error {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()

	oldState := e.GetEventState(e)

	now := e.Now(oldState)

	var changed bool
	var err error
	e.SettingsMu.Lock()
	if changed, err = e.Setts.Deserialize(params.SettingsValue); !changed {
		if err != nil {
			log.Warn(err.Error())
		}

		if oldState.LastUpdated != nil {
			delta := now.Sub(*oldState.LastUpdated).Milliseconds()
			if delta < 200 {
				e.SettingsMu.Unlock()
				return nil
			}
		}
	}
	e.SettingsMu.Unlock()

	e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventStateChanged{
		Type:     e.Id.Type(),
		EntityId: e.Id,
		OldState: oldState,
		NewState: e.GetEventState(e),
	})

	return nil
}
