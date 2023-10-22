// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package version

import (
	"runtime"
	"sync"

	m "github.com/e154/smart-home/models"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/system/supervisor"
	"github.com/e154/smart-home/version"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	updateLock *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	actor := &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		updateLock: &sync.Mutex{},
	}
	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}
	return actor
}

func (e *Actor) Destroy() {

}

func (e *Actor) Spawn() {

}

func (e *Actor) selfUpdate() {

	e.updateLock.Lock()
	defer e.updateLock.Unlock()

	oldState := e.GetEventState()
	e.Now(oldState)

	var s runtime.MemStats
	runtime.ReadMemStats(&s)

	e.AttrMu.Lock()
	e.Attrs[AttrVersion].Value = version.VersionString
	e.Attrs[AttrRevision].Value = version.RevisionString
	e.Attrs[AttrRevisionURL].Value = version.RevisionURLString
	e.Attrs[AttrGenerated].Value = version.GeneratedString
	e.Attrs[AttrDevelopers].Value = version.DevelopersString
	e.Attrs[AttrBuildNum].Value = version.BuildNumString
	e.Attrs[AttrDockerImage].Value = version.DockerImageString
	e.Attrs[AttrGoVersion].Value = version.GoVersion
	e.AttrMu.Unlock()

	go e.SaveState(events.EventStateChanged{
		StorageSave: false,
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(),
	})
}
