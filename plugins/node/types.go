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

package node

import (
	"encoding/json"
	"time"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	// Name ...
	Name = "node"
	// TopicPluginNode ...
	TopicPluginNode = "system/plugins/node"

	Version = "0.0.1"
)

const (
	// AttrThread ...
	AttrThread = "thread"
	// AttrRps ...
	AttrRps = "rps"
	// AttrMin ...
	AttrMin = "min"
	// AttrMax ...
	AttrMax = "max"
	// AttrLatency ...
	AttrLatency = "latency"
	// AttrStartedAt ...
	AttrStartedAt = "started_at"
	// AttrNodeLogin ...
	AttrNodeLogin = "node_login"
	// AttrNodePass ...
	AttrNodePass = "node_pass"
)

// NewAttr ...
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
		AttrLatency: {
			Name: AttrLatency,
			Type: common.AttributeInt,
		},
		AttrStartedAt: {
			Name: AttrStartedAt,
			Type: common.AttributeTime,
		},
	}
}

// NewSettings ...
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

// NewStates ...
func NewStates() (states map[string]supervisor.ActorState) {

	states = map[string]supervisor.ActorState{
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

// MessageStatus ...
type MessageStatus struct {
	Status    string    `json:"status"`
	Thread    int       `json:"thread"`
	Rps       int64     `json:"rps"`
	Min       int64     `json:"min"`
	Max       int64     `json:"max"`
	Laatency  int64     `json:"latency"`
	StartedAt time.Time `json:"started_at"`
}
