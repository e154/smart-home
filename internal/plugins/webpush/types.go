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

package webpush

import (
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
)

const (
	// Name ...
	Name = "webpush"

	AttrTitle          = "title"
	AttrUserIDS        = "userIDS"
	AttrBody           = "body"
	AttrPublicKey      = "public_key"
	AttrPrivateKey     = "private_key"
	TopicPluginWebpush = "system/plugins/webpush"
)

// NewMessageParams ...
func NewMessageParams() m.Attributes {
	return map[string]*m.Attribute{
		AttrUserIDS: {
			Name: AttrUserIDS,
			Type: common.AttributeString,
		},
		AttrTitle: {
			Name: AttrTitle,
			Type: common.AttributeString,
		},
		AttrBody: {
			Name: AttrBody,
			Type: common.AttributeString,
		},
	}
}

// NewSettings ...
func NewSettings() m.Attributes {
	return map[string]*m.Attribute{
		AttrPublicKey: {
			Name: AttrPublicKey,
			Type: common.AttributeEncrypted,
		},
		AttrPrivateKey: {
			Name: AttrPrivateKey,
			Type: common.AttributeEncrypted,
		},
	}
}

type Notification struct {
	Title string `json:"title"`
}
