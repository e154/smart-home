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

package telegram

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	// Name ...
	Name = "telegram"

	AttrToken     = "token"
	AttrPin       = "pin"
	AttrChatID    = "chat_id"
	AttrBody      = "body"
	AttrPhotoUri  = "photo_uri"
	AttrPhotoPath = "photo_path"
	AttrFilePath  = "file_path"
	AttrFileUri   = "file_uri"
	AttrKeys      = "keys"

	Version = "0.0.1"
)

const (
	AttrConnected = "connected"
	AttrOffline   = "offline"
)

const (
	// FuncEntityAction ...
	FuncEntityAction = "telegramAction"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return nil
}

// NewMessageParams ...
func NewMessageParams() m.Attributes {
	return map[string]*m.Attribute{
		AttrChatID: {
			Name: AttrChatID,
			Type: common.AttributeInt,
		},
		AttrBody: {
			Name: AttrBody,
			Type: common.AttributeString,
		},
		AttrPhotoUri: {
			Name: AttrPhotoUri,
			Type: common.AttributeArray,
		},
		AttrPhotoPath: {
			Name: AttrPhotoPath,
			Type: common.AttributeArray,
		},
		AttrFileUri: {
			Name: AttrFileUri,
			Type: common.AttributeArray,
		},
		AttrFilePath: {
			Name: AttrFilePath,
			Type: common.AttributeArray,
		},
		AttrKeys: {
			Name: AttrKeys,
			Type: common.AttributeArray,
		},
	}
}

// NewSettings ...
func NewSettings() m.Attributes {
	return map[string]*m.Attribute{
		AttrToken: {
			Name: AttrToken,
			Type: common.AttributeEncrypted,
		},
		AttrPin: {
			Name: AttrPin,
			Type: common.AttributeEncrypted,
		},
	}
}

// NewStates ...
func NewStates() (states map[string]supervisor.ActorState) {

	states = map[string]supervisor.ActorState{
		AttrConnected: {
			Name:        AttrConnected,
			Description: "connected",
		},
		AttrOffline: {
			Name:        AttrOffline,
			Description: "offline",
		},
	}

	return
}

// Command ...
type Command struct {
	UserName, Text string
	ChatId         int64
}
