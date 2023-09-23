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

package version

import (
	m "github.com/e154/smart-home/models"
	"runtime"
	"sync"

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
		BaseActor: supervisor.NewBaseActor(entity, service),
		updateLock: &sync.Mutex{},
	}
	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}
	return actor
}

func (u *Actor) Destroy() {

}

func (u *Actor) Spawn() {

}

func (u *Actor) selfUpdate() {

	u.updateLock.Lock()
	defer u.updateLock.Unlock()

	oldState := u.GetEventState()
	u.Now(oldState)

	var s runtime.MemStats
	runtime.ReadMemStats(&s)

	u.AttrMu.Lock()
	u.Attrs[AttrVersion].Value = version.VersionString
	u.Attrs[AttrRevision].Value = version.RevisionString
	u.Attrs[AttrRevisionURL].Value = version.RevisionURLString
	u.Attrs[AttrGenerated].Value = version.GeneratedString
	u.Attrs[AttrDevelopers].Value = version.DevelopersString
	u.Attrs[AttrBuildNum].Value = version.BuildNumString
	u.Attrs[AttrDockerImage].Value = version.DockerImageString
	u.Attrs[AttrGoVersion].Value = version.GoVersion
	u.AttrMu.Unlock()

	u.Service.EventBus().Publish("system/entities/"+u.Id.String(), events.EventStateChanged{
		StorageSave: false,
		PluginName:  u.Id.PluginName(),
		EntityId:    u.Id,
		OldState:    oldState,
		NewState:    u.GetEventState(),
	})
}
