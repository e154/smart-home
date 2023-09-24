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
	m "github.com/e154/smart-home/models"
	"sync"
	"time"

	"github.com/Masterminds/semver"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/web"
	"github.com/e154/smart-home/system/supervisor"
	"github.com/e154/smart-home/version"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	checkLock         *sync.Mutex
	latestVersion     string
	latestDownloadUrl string
	latestVersionTime time.Time
	lastCheck         time.Time
	currentVersion    *semver.Version
}

// NewActor ...
func NewActor(entity *m.Entity, service supervisor.Service) *Actor {

	var v = "v0.0.1"
	if version.VersionString != "?" {
		v = version.VersionString
	}
	currentVersion, err := semver.NewVersion(v)
	if err != nil {
		log.Error(err.Error())
	}

	actor := &Actor{
		BaseActor:      supervisor.NewBaseActor(entity, service),
		checkLock:      &sync.Mutex{},
		currentVersion: currentVersion,
	}

	if actor.Actions == nil {
		actor.Actions = NewActions()
	}

	if actor.States == nil {
		actor.States = NewStates()
	}

	return actor
}

func (e *Actor) Destroy() {

}

func (e *Actor) Spawn() {
	e.check()

}

func (e *Actor) setState(v string) {

	switch v {
	case "exist_update":
		state := e.States["exist_update"]
		e.State = &state
		e.Value.Store(supervisor.StateOk)
		return
	case supervisor.StateAwait, supervisor.StateOk, supervisor.StateInProcess:
		state := e.States["enabled"]
		e.State = &state
	case supervisor.StateError:
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
			e.setState(supervisor.StateError)
			return
		}
		e.checkLock.Unlock()
	}()

	e.setState(supervisor.StateInProcess)

	var body []byte
	if _, body, err = e.Service.Crawler().Probe(web.Request{Method: "GET", Url: uri, Timeout: 5 * time.Second}); err != nil {
		return
	}

	data := GithubReleaseLatest{}
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	e.setState(supervisor.StateOk)

	e.lastCheck = time.Now()
	e.latestVersion = data.TagName
	e.latestVersionTime = data.CreatedAt
	for _, asset := range data.Assets {
		e.latestDownloadUrl = asset.BrowserDownloadUrl
	}

	var releaseVersion *semver.Version
	if releaseVersion, err = semver.NewVersion(e.latestVersion); err == nil {
		// found update
		if e.currentVersion != nil {
			if compare := e.currentVersion.Compare(releaseVersion); compare < 0 {
				e.setState("exist_update")
			}
		}
	}

	oldState := e.GetEventState()

	e.AttrMu.Lock()
	e.Attrs[AttrUpdaterLatestVersion].Value = e.latestVersion
	e.Attrs[AttrUpdaterLatestVersionTime].Value = e.latestVersionTime
	e.Attrs[AttrUpdaterLatestLatestDownloadUrl].Value = e.latestDownloadUrl
	e.Attrs[AttrUpdaterLatestCheck].Value = e.lastCheck
	e.AttrMu.Unlock()

	go e.SaveState(events.EventStateChanged{
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(),
		StorageSave: true,
	})
}
