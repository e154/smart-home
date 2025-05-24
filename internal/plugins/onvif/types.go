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

package onvif

import (
	"time"

	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

const (
	// Name ...
	Name = "onvif"
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"

	AttrAddress              = "address"
	AttrOnvifPort            = "onvifPort"
	AttrRtspPort             = "rtspPort"
	AttrUserName             = "userName"
	AttrPassword             = "password"
	AttrConnected            = "connected"
	AttrOffline              = "offline"
	AttrRequireAuthorization = "requireAuthorization"

	AttrMotion     = "motion"
	AttrMotionTime = "motionTime"

	ActionContinuousMove     = "continuousMove"
	ActionStopContinuousMove = "stopContinuousMove"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrMotion: {
			Name: AttrMotion,
			Type: common.AttributeBool,
		},
	}
}

// NewSettings ...
func NewSettings() map[string]*m.Attribute {
	return map[string]*m.Attribute{
		AttrAddress: {
			Name:  AttrAddress,
			Type:  common.AttributeString,
			Value: "192.168.0.1",
		},
		AttrOnvifPort: {
			Name:  AttrOnvifPort,
			Type:  common.AttributeInt,
			Value: "80",
		},
		AttrRtspPort: {
			Name:  AttrRtspPort,
			Type:  common.AttributeInt,
			Value: "554",
		},
		AttrUserName: {
			Name:  AttrUserName,
			Type:  common.AttributeString,
			Value: "admin",
		},
		AttrPassword: {
			Name: AttrPassword,
			Type: common.AttributeEncrypted,
		},
		AttrRequireAuthorization: {
			Name:  AttrRequireAuthorization,
			Type:  common.AttributeBool,
			Value: true,
		},
	}
}

// NewStates ...
func NewStates() (states map[string]plugins.ActorState) {

	states = map[string]plugins.ActorState{
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

func NewActions() map[string]plugins.ActorAction {
	return map[string]plugins.ActorAction{
		ActionContinuousMove: {
			Name:        ActionContinuousMove,
			Description: "camera control",
		},
		ActionStopContinuousMove: {
			Name:        ActionStopContinuousMove,
			Description: "camera control",
		},
	}
}

type MotionAlarm struct {
	State bool
	Time  time.Time
}

type ConnectionStatus struct {
	Connected bool
}

type StreamList struct {
	List        []string
	SnapshotUri *string
}
