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

package node

import (
	"encoding/json"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"time"
)

const (
	Name            = "node"
	EntityNode      = common.EntityType("node")
	TopicPluginNode = "plugin.node"
)

const (
	AttrThread    = "thread"
	AttrRps       = "rps"
	AttrMin       = "min"
	AttrMax       = "max"
	AttrStartedAt = "started_at"
	AttrNodeLogin = "node_login"
	AttrNodePass  = "node_pass"
)

func NewAttr() m.Attributes {
	return m.Attributes{
		AttrThread: {
			Name: AttrThread,
			Type: common.AttributeInt,
		},
		AttrRps: {
			Name: AttrRps,
			Type: common.AttributeInt,
		},
		AttrMin: {
			Name: AttrMin,
			Type: common.AttributeInt,
		},
		AttrMax: {
			Name: AttrMax,
			Type: common.AttributeInt,
		},
		AttrStartedAt: {
			Name: AttrStartedAt,
			Type: common.AttributeTime,
		},
	}
}

func NewSettings() m.Attributes {
	return m.Attributes{
		AttrNodeLogin: {
			Name: AttrNodeLogin,
			Type: common.AttributeString,
		},
		AttrNodePass: {
			Name: AttrNodePass,
			Type: common.AttributeString,
		},
	}
}

func NewStates() (states map[string]entity_manager.ActorState) {

	states = map[string]entity_manager.ActorState{
		"wait": {
			Name:        "wait",
			Description: "Wait",
		},
		"connected": {
			Name:        "connected",
			Description: "Connected",
		},
		"error": {
			Name:        "error",
			Description: "Error",
		},
	}

	return
}

// MqttMessage ...
type MqttMessage struct {
	Dup      bool
	Qos      uint8
	Retained bool
	Topic    string
	PacketID uint16
	Payload  []byte
}

// DeviceType ...
type DeviceType string

// MessageRequest ...
type MessageRequest struct {
	EntityId   common.EntityId `json:"entity_id"`
	DeviceType DeviceType      `json:"device_type"`
	Properties json.RawMessage `json:"properties"`
	Command    json.RawMessage `json:"command"`
}

// MessageResponse ...
type MessageResponse struct {
	EntityId   common.EntityId `json:"entity_id"`
	DeviceType DeviceType      `json:"device_type"`
	Properties json.RawMessage `json:"properties"`
	Response   json.RawMessage `json:"response"`
	Status     string          `json:"status"`
}

type MessageStatus struct {
	Status    string    `json:"status"`
	Thread    int       `json:"thread"`
	Rps       int64     `json:"rps"`
	Min       int64     `json:"min"`
	Max       int64     `json:"max"`
	StartedAt time.Time `json:"started_at"`
}
