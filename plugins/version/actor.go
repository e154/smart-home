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
	"fmt"
	"runtime"
	"sync"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/supervisor"
	"github.com/e154/smart-home/version"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	eventBus   bus.Bus
	updateLock *sync.Mutex
}

// NewActor ...
func NewActor(visor supervisor.Supervisor,
	eventBus bus.Bus) *Actor {

	actor := &Actor{
		BaseActor: supervisor.BaseActor{
			Id:                common.EntityId(fmt.Sprintf("%s.%s", EntityVersion, Name)),
			Name:              Name,
			EntityType:        EntityVersion,
			UnitOfMeasurement: "",
			AttrMu:            &sync.RWMutex{},
			Attrs:             NewAttr(),
			Supervisor:        visor,
		},
		eventBus:   eventBus,
		updateLock: &sync.Mutex{},
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
	u.Attrs[AttrVersion].Value = version.VersionString
	u.Attrs[AttrRevision].Value = version.RevisionString
	u.Attrs[AttrRevisionURL].Value = version.RevisionURLString
	u.Attrs[AttrGenerated].Value = version.GeneratedString
	u.Attrs[AttrDevelopers].Value = version.DevelopersString
	u.Attrs[AttrBuildNum].Value = version.BuildNumString
	u.Attrs[AttrDockerImage].Value = version.DockerImageString
	u.Attrs[AttrGoVersion].Value = version.GoVersion
	u.AttrMu.Unlock()

	u.eventBus.Publish(bus.TopicEntities, events.EventStateChanged{
		StorageSave: false,
		PluginName:  u.Id.PluginName(),
		EntityId:    u.Id,
		OldState:    oldState,
		NewState:    u.GetEventState(u),
	})
}
