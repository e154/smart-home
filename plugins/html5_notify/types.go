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

package html5_notify

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	// Name ...
	Name = "html5_notify"

	AttrTitle              = "title"
	AttrUserIDS            = "userIDS"
	AttrActions            = "actions"
	AttrBadge              = "badge"
	AttrBody               = "body"
	AttrData               = "data"
	AttrDir                = "dir"
	AttrIcon               = "icon"
	AttrImage              = "image"
	AttrLang               = "lang"
	AttrRenotify           = "renotify"
	AttrRequireInteraction = "requireInteraction"
	AttrSilent             = "silent"
	AttrTag                = "tag"
	AttrTimestamp          = "timestamp"

	Version = "0.0.1"
)

// NewMessageParams ...
func NewMessageParams() m.Attributes {
	return map[string]*m.Attribute{
		AttrTitle: {
			Name: AttrTitle,
			Type: common.AttributeString,
		},
		AttrUserIDS: {
			Name: AttrUserIDS,
			Type: common.AttributeString,
		},
		AttrActions: {
			Name: AttrActions,
			Type: common.AttributeArray,
		},
		AttrBadge: {
			Name: AttrBadge,
			Type: common.AttributeString,
		},
		AttrBody: {
			Name: AttrBody,
			Type: common.AttributeString,
		},
		AttrData: {
			Name: AttrData,
			Type: common.AttributeString,
		},
		AttrDir: {
			Name: AttrDir,
			Type: common.AttributeString,
		},
		AttrIcon: {
			Name: AttrIcon,
			Type: common.AttributeString,
		},
		AttrImage: {
			Name: AttrImage,
			Type: common.AttributeString,
		},
		AttrLang: {
			Name: AttrLang,
			Type: common.AttributeString,
		},
		AttrRenotify: {
			Name: AttrRenotify,
			Type: common.AttributeBool,
		},
		AttrRequireInteraction: {
			Name: AttrRequireInteraction,
			Type: common.AttributeBool,
		},
		AttrSilent: {
			Name: AttrSilent,
			Type: common.AttributeBool,
		},
		AttrTag: {
			Name: AttrTag,
			Type: common.AttributeString,
		},
		AttrTimestamp: {
			Name: AttrTimestamp,
			Type: common.AttributeInt,
		},
	}
}

// NewSettings ...
func NewSettings() map[string]*m.Attribute {
	return map[string]*m.Attribute{}
}

type NotificationAction struct {
	Action string `json:"action"`
	Icon   string `json:"icon"`
	Title  string `json:"title"`
}

type NotificationOptions struct {
	Actions            []NotificationAction `json:"actions,omitempty"`
	Badge              string               `json:"badge,omitempty"`
	Body               string               `json:"body,omitempty"`
	Data               any                  `json:"data,omitempty"`
	Dir                string               `json:"dir,omitempty"`
	Icon               string               `json:"icon,omitempty"`
	Image              string               `json:"image,omitempty"`
	Lang               string               `json:"lang,omitempty"`
	Renotify           bool                 `json:"renotify,omitempty"`
	RequireInteraction bool                 `json:"requireInteraction,omitempty"`
	Silent             bool                 `json:"silent,omitempty"`
	Tag                string               `json:"tag,omitempty"`
	Timestamp          int64                `json:"timestamp,omitempty"`
}

type Notification struct {
	Title   string               `json:"title"`
	Options *NotificationOptions `json:"options,omitempty"`
}
