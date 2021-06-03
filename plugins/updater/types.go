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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"time"
)

type GithubRelease struct {
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

type GithubReleaseLatest struct {
	TagName   string          `json:"tag_name"`
	CreatedAt time.Time       `json:"created_at"`
	Assets    []GithubRelease `json:"assets"`
}

const (
	AttrUpdaterLatestVersion           = "latest_version"
	AttrUpdaterLatestVersionTime       = "latest_version_time"
	AttrUpdaterLatestLatestDownloadUrl = "latest_download_url"
	AttrUpdaterLatestCheck             = "last_check"
	Name                               = "updater"
	EntityUpdater                      = common.EntityType(Name)
)

func NewAttr() m.Attributes {
	return m.Attributes{
		AttrUpdaterLatestVersion: {
			Name: AttrUpdaterLatestVersion,
			Type: common.AttributeString,
		},
		AttrUpdaterLatestVersionTime: {
			Name: AttrUpdaterLatestVersionTime,
			Type: common.AttributeTime,
		},
		AttrUpdaterLatestLatestDownloadUrl: {
			Name: AttrUpdaterLatestLatestDownloadUrl,
			Type: common.AttributeString,
		},
		AttrUpdaterLatestCheck: {
			Name: AttrUpdaterLatestCheck,
			Type: common.AttributeTime,
		},
	}
}

func NewStates() map[string]entity_manager.ActorState {
	return map[string]entity_manager.ActorState{
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
	}
}

func NewActions() map[string]entity_manager.ActorAction {
	return map[string]entity_manager.ActorAction{
		"check": {
			Name:        "check",
			Description: "Check version",
		},
	}
}
