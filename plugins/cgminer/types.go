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

package cgminer

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	// Name ...
	Name = "cgminer"
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"
)

const (
	// SettingHost ...
	SettingHost = "host"
	// SettingPort ...
	SettingPort = "port"
	// SettingTimeout ...
	SettingTimeout = "timeout"
	// SettingUser ...
	SettingUser = "user"
	// SettingPass ...
	SettingPass = "pass"
	// SettingManufacturer ...
	SettingManufacturer = "manufacturer"
	// SettingModel ...
	SettingModel = "model"

	// StateEnabled ...
	StateEnabled = "ENABLED"
	// StateDisabled ...
	StateDisabled = "DISABLED"
	// StateError ...
	StateError = "ERROR"

	// ActionEnabled ...
	ActionEnabled = "ENABLE"
	// ActionDisable ...
	ActionDisable = "DISABLE"
	// ActionCheck ...
	ActionCheck = "CHECK"

	Version = "0.0.1"
)

// store entity status in this struct
func NewAttr() m.Attributes {
	return nil
}

// entity settings
func NewSettings() m.Attributes {
	return m.Attributes{
		SettingHost: {
			Name: SettingHost,
			Type: common.AttributeString,
		},
		SettingManufacturer: {
			Name: SettingManufacturer,
			Type: common.AttributeString,
		},
		SettingModel: {
			Name: SettingModel,
			Type: common.AttributeString,
		},
		SettingPort: {
			Name: SettingPort,
			Type: common.AttributeInt,
		},
		SettingTimeout: {
			Name: SettingTimeout,
			Type: common.AttributeInt,
		},
		SettingUser: {
			Name: SettingUser,
			Type: common.AttributeString,
		},
		SettingPass: {
			Name: SettingPass,
			Type: common.AttributeString,
		},
	}
}

// state list entity
func NewStates() (states map[string]supervisor.ActorState) {

	states = map[string]supervisor.ActorState{
		StateEnabled: {
			Name:        StateEnabled,
			Description: "Enabled",
		},
		StateDisabled: {
			Name:        StateDisabled,
			Description: "Disabled",
		},
		StateError: {
			Name:        StateError,
			Description: "Error",
		},
	}

	return
}

// entity action list
func NewActions() map[string]supervisor.ActorAction {
	return map[string]supervisor.ActorAction{
		ActionEnabled: {
			Name:        "Enable",
			Description: "enable action",
		},
		ActionDisable: {
			Name:        "Disable",
			Description: "disable action",
		},
		ActionCheck: {
			Name:        "Check",
			Description: "check action",
		},
	}
}

// IMiner ...
type IMiner interface {
	Stats() (data []byte, err error)
	Devs() (data []byte, err error)
	Summary() (data []byte, err error)
	Pools() (data []byte, err error)
	AddPool(url string) (err error)
	Version() (data []byte, err error)
	Enable(poolId int64) error
	Disable(poolId int64) error
	Delete(poolId int64) error
	SwitchPool(poolId int64) error
	Restart() error
	Quit() error
	Bind() interface{}
}
