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
	"github.com/Masterminds/semver"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/version"
	"go.uber.org/atomic"
	"io"
	"net/http"
	"sync"
	"time"
)

type Actor struct {
	entity_manager.BaseActor
	checkLock         *sync.Mutex
	latestVersion     string
	latestDownloadUrl string
	latestVersionTime time.Time
	lastCheck         time.Time
	currentVersion    *semver.Version
}

func NewActor(entityManager entity_manager.EntityManager) *Actor {
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
			AttrMu:      &sync.Mutex{},
			Attrs:       NewAttr(),
			Manager:     entityManager,
			States: map[string]entity_manager.ActorState{
				"enabled": {
					Name:        "enabled",
					Description: "Enabled",
				},
				"disabled": {
					Name:        "disabled",
					Description: "Disabled",
				},
				"error": {
					Name:        "error",
					Description: "Error",
				},
				"exist_update": {
					Name:        "exist_update",
					Description: "Exist update",
				},
			},
			Actions: map[string]entity_manager.ActorAction{
				"check": {
					Name:        "check",
					Description: "Check version",
				},
			},
		},
		checkLock:      &sync.Mutex{},
		currentVersion: currentVersion,
	}
}

func (e *Actor) Spawn() entity_manager.PluginActor {
	return e
}

func (u *Actor) setState(v string) {

	switch v {
	case "exist_update":
		state := u.States["exist_update"]
		u.State = &state
		u.Value.Store(entity_manager.StateOk)
		return
	case entity_manager.StateAwait, entity_manager.StateOk, entity_manager.StateInProcess:
		state := u.States["enabled"]
		u.State = &state
	case entity_manager.StateError:
		state := u.States["error"]
		u.State = &state
	default:
		state := u.States["disabled"]
		u.State = &state
	}

	u.Value.Store(v)
}

func (u *Actor) check() {

	u.checkLock.Lock()
	var err error
	defer func() {
		if err != nil {
			u.setState(entity_manager.StateError)
			return
		}
		u.checkLock.Unlock()
	}()

	u.setState(entity_manager.StateInProcess)

	var resp *http.Response
	if resp, err = http.Get(uri); err != nil {
		return
	}
	defer resp.Body.Close()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return
	}

	data := GithubReleaseLatest{}
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	u.setState(entity_manager.StateOk)

	u.lastCheck = time.Now()
	u.latestVersion = data.TagName
	u.latestVersionTime = data.CreatedAt
	for _, asset := range data.Assets {
		u.latestDownloadUrl = asset.BrowserDownloadUrl
	}

	var releaseVersion *semver.Version
	if releaseVersion, err = semver.NewVersion(u.latestVersion); err == nil {
		// found update
		if compare := u.currentVersion.Compare(releaseVersion); compare < 0 {
			u.setState("exist_update")
		}
	}

	oldState := u.GetEventState(u)

	u.AttrMu.Lock()
	u.Attrs[AttrUpdaterLatestVersion].Value = u.latestVersion
	u.Attrs[AttrUpdaterLatestVersionTime].Value = u.latestVersionTime
	u.Attrs[AttrUpdaterLatestLatestDownloadUrl].Value = u.latestDownloadUrl
	u.Attrs[AttrUpdaterLatestCheck].Value = u.lastCheck
	u.AttrMu.Unlock()

	u.Send(entity_manager.MessageStateChanged{
		StorageSave: true,
		OldState:    oldState,
		NewState:    u.GetEventState(u),
	})
}
