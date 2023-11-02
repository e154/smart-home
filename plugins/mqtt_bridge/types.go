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

package mqtt_bridge

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	// Name ...
	Name = "mqtt_bridge"
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"

	Version = "0.0.1"
)

type Direction string

const (
	DirectionIn   = Direction("in")
	DirectionOut  = Direction("out")
	DirectionBoth = Direction("both")
)

const (
	AttrConnected = "connected"
	AttrOffline   = "offline"

	AttrKeepAlive      = "keepAlive"
	AttrPingTimeout    = "pingTimeout"
	AttrBroker         = "broker"
	AttrClientID       = "clientID"
	AttrConnectTimeout = "connectTimeout"
	AttrCleanSession   = "cleanSession"
	AttrUsername       = "username"
	AttrPassword       = "password"
	AttrQos            = "qos"
	AttrDirection      = "direction"
	AttrTopics         = "topics"
)

// NewSettings ...
func NewSettings() m.Attributes {
	return m.Attributes{
		AttrKeepAlive: {
			Name:  AttrKeepAlive,
			Type:  common.AttributeInt,
			Value: 15,
		},
		AttrPingTimeout: {
			Name:  AttrPingTimeout,
			Type:  common.AttributeInt,
			Value: 10,
		},
		AttrBroker: {
			Name:  AttrBroker,
			Type:  common.AttributeString,
			Value: "tcp://server:port",
		},
		AttrClientID: {
			Name: AttrClientID,
			Type: common.AttributeString,
		},
		AttrConnectTimeout: {
			Name:  AttrConnectTimeout,
			Type:  common.AttributeInt,
			Value: 30,
		},
		AttrCleanSession: {
			Name:  AttrCleanSession,
			Type:  common.AttributeBool,
			Value: false,
		},
		AttrUsername: {
			Name: AttrUsername,
			Type: common.AttributeString,
		},
		AttrPassword: {
			Name: AttrPassword,
			Type: common.AttributeEncrypted,
		},
		AttrQos: {
			Name:  AttrQos,
			Type:  common.AttributeInt,
			Value: 0,
		},
		AttrDirection: {
			Name:  AttrDirection,
			Type:  common.AttributeString,
			Value: DirectionIn,
		},
		AttrTopics: {
			Name:  AttrTopics,
			Type:  common.AttributeString,
			Value: "owntracks/#",
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
