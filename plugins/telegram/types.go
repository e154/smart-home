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

package telegram

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	Name      = "telegram"
	AttrToken = "token"
	AttrName  = "name"
	AttrBody  = "body"
)

const (
	// StatusDelivered ...
	StatusDelivered = "delivered"
)

func NewAttr() m.Attributes {
	return nil
}

func NewMessageParams() m.Attributes {
	return map[string]*m.Attribute{
		AttrName: {
			Name: AttrName,
			Type: common.AttributeString,
		},
		AttrBody: {
			Name: AttrBody,
			Type: common.AttributeString,
		},
	}
}

func NewSettings() m.Attributes {
	return map[string]*m.Attribute{
		AttrToken: {
			Name: AttrToken,
			Type: common.AttributeString,
		},
	}
}

// Command ...
type Command struct {
	UserName, Text string
	ChatId         int64
}

const banner = `
Smart home system

Version:
%s

command:
%s
`
