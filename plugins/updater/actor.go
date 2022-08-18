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

package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/e154/smart-home/common/events"

	"github.com/Masterminds/semver"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/version"
	"go.uber.org/atomic"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	eventBus          bus.Bus
	checkLock         *sync.Mutex
	latestVersion     string
	latestDownloadUrl string
	latestVersionTime time.Time
	lastCheck         time.Time
	currentVersion    *semver.Version
}

// NewActor ...
func NewActor(entityManager entity_manager.EntityManager,
	eventBus bus.Bus) *Actor {

	var v = "v0.0.1"
	if version.VersionString != "?" {
		v = version.VersionString
	}
	currentVersion, err := semver.NewVersion(v)
	if err != nil {
		log.Error(err.Error())
	}

	return &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:          common.EntityId(fmt.Sprintf("%s.%s", EntityUpdater, Name)),
			Name:        Name,
			Description: "sun plugin",
			EntityType:  EntityUpdater,
			Value:       atomic.NewString(entity_manager.StateAwait),
			AttrMu:      &sync.RWMutex{},
			Attrs:       NewAttr(),
			Manager:     entityManager,
			States:      NewStates(),
			Actions:     NewActions(),
		},
		eventBus:       eventBus,
		checkLock:      &sync.Mutex{},
		currentVersion: currentVersion,
	}
}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {
	return e
}

func (e *Actor) setState(v string) {

	switch v {
	case "exist_update":
		state := e.States["exist_update"]
		e.State = &state
		e.Value.Store(entity_manager.StateOk)
		return
	case entity_manager.StateAwait, entity_manager.StateOk, entity_manager.StateInProcess:
		state := e.States["enabled"]
		e.State = &state
	case entity_manager.StateError:
		state := e.States["error"]
		e.State = &state
	default:
		state := e.States["disabled"]
		e.State = &state
	}

	e.Value.Store(v)
}

func (e *Actor) check() {

	e.checkLock.Lock()
	var err error
	defer func() {
		if err != nil {
			e.setState(entity_manager.StateError)
			return
		}
		e.checkLock.Unlock()
	}()

	e.setState(entity_manager.StateInProcess)

	var resp *http.Response
	if resp, err = http.Get(uri); err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return
	}

	data := GithubReleaseLatest{}
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	e.setState(entity_manager.StateOk)

	e.lastCheck = time.Now()
	e.latestVersion = data.TagName
	e.latestVersionTime = data.CreatedAt
	for _, asset := range data.Assets {
		e.latestDownloadUrl = asset.BrowserDownloadUrl
	}

	var releaseVersion *semver.Version
	if releaseVersion, err = semver.NewVersion(e.latestVersion); err == nil {
		// found update
		if compare := e.currentVersion.Compare(releaseVersion); compare < 0 {
			e.setState("exist_update")
		}
	}

	oldState := e.GetEventState(e)

	e.AttrMu.Lock()
	e.Attrs[AttrUpdaterLatestVersion].Value = e.latestVersion
	e.Attrs[AttrUpdaterLatestVersionTime].Value = e.latestVersionTime
	e.Attrs[AttrUpdaterLatestLatestDownloadUrl].Value = e.latestDownloadUrl
	e.Attrs[AttrUpdaterLatestCheck].Value = e.lastCheck
	e.AttrMu.Unlock()

	e.eventBus.Publish(bus.TopicEntities, events.EventStateChanged{
		PluginName: e.Id.PluginName(),
		EntityId:   e.Id,
		OldState:   oldState,
		NewState:   e.GetEventState(e),
	})
}
