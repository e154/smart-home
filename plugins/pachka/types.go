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

package pachka

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"time"
)

const (
	// Name ...
	Name = "pachka"

	AttrToken    = "token"
	AttrEntityID = "entity_id"
	AttrBody     = "body"

	Version = "0.0.1"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return nil
}

// NewMessageParams ...
func NewMessageParams() m.Attributes {
	return map[string]*m.Attribute{
		AttrEntityID: {
			Name: AttrEntityID,
			Type: common.AttributeInt,
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
		AttrToken: {
			Name: AttrToken,
			Type: common.AttributeEncrypted,
		},
	}
}

// Command ...
type Command struct {
	UserName, Text string
	ChatId         int64
}

type EntityType string

type RequestMessage struct {
	EntityType *EntityType `json:"entity_type,omitempty"`
	EntityId   int64       `json:"entity_id,omitempty"`
	Content    string      `json:"content"`
}

type ResponseMessage struct {
	Data struct {
		ID        int64 `json:"id"`
		ChatId    int64 `json:"chat_id"`
		Content   string
		UserId    int64
		CreatedAt time.Time
	} `json:"data"`
}

type ErrorItem struct {
	Key     string      `json:"key"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Payload interface{} `json:"payload"`
}

type ResponseError struct {
	Errors []*ErrorItem `json:"errors"`
}
