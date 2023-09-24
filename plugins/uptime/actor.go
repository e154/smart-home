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
	"time"

	"go.uber.org/atomic"

	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	appStarted time.Time
	total      *atomic.Uint64
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {
	actor := &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		appStarted: time.Now(),
		total:      atomic.NewUint64(0),
	}
	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}
	return actor
}

func (e *Actor) Destroy() {

}

func (e *Actor) Spawn() {
	e.update()
}

func (e *Actor) update() {

	oldState := e.GetEventState()

	e.Now(oldState)

	total, err := GetUptime()
	if err != nil {
		return
	}

	e.total.Store(total)

	e.AttrMu.Lock()
	e.Attrs[AttrUptimeTotal].Value = e.total.Load()
	e.Attrs[AttrUptimeAppStarted].Value = e.appStarted
	e.AttrMu.Unlock()

	go e.SaveState(events.EventStateChanged{
		PluginName: e.Id.PluginName(),
		EntityId:   e.Id,
		OldState:   oldState,
		NewState:   e.GetEventState(),
	})
}
