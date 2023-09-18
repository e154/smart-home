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

package onvif

import (
	"time"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	// Name ...
	Name = "onvif"
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"

	Version = "0.0.1"

	AttrAddress   = "address"
	AttrOnvifPort = "onvifPort"
	AttrRtspPort  = "rtspPort"
	AttrUserName  = "userName"
	AttrPassword  = "password"
	AttrConnected = "connected"
	AttrOffline   = "offline"

	AttrStreamUri   = "streamUri"
	AttrSnapshotUri = "snapshotUri"

	AttrMotion     = "motion"
	AttrMotionTime = "motionTime"

	ActionContinuousMove     = "continuousMove"
	ActionStopContinuousMove = "stopContinuousMove"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrStreamUri: {
			Name: AttrStreamUri,
			Type: common.AttributeEncrypted,
		},
		AttrSnapshotUri: {
			Name: AttrSnapshotUri,
			Type: common.AttributeEncrypted,
		},
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

func NewActions() map[string]supervisor.ActorAction {
	return map[string]supervisor.ActorAction{
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
