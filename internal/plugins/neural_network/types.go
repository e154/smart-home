// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package neural_network

import (
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

const (
	Name                = "neural_network"
	EntityNeuralNetwork = string("neural_network")
)

const (
	SettingParam1 = "param1"
	SettingParam2 = "param2"

	AttrPhase = "phase"

	StateEnabled  = "enabled"
	StateDisabled = "disabled"

	ActionEnabled = "enable"
	ActionDisable = "disable"
)

// store entity status in this struct
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrPhase: {
			Name: AttrPhase,
			Type: common.AttributeString,
		},
	}
}

// entity settings
func NewSettings() m.Attributes {
	return m.Attributes{
		SettingParam1: {
			Name: SettingParam1,
			Type: common.AttributeString,
		},
		SettingParam2: {
			Name: SettingParam2,
			Type: common.AttributeString,
		},
	}
}

// state list entity
func NewStates() (states map[string]plugins.ActorState) {

	states = map[string]plugins.ActorState{
		StateEnabled: {
			Name:        StateEnabled,
			Description: "Enabled",
		},
		StateDisabled: {
			Name:        StateDisabled,
			Description: "disabled",
		},
	}

	return
}

// entity action list
func NewActions() map[string]plugins.ActorAction {
	return map[string]plugins.ActorAction{
		ActionEnabled: {
			Name:        ActionEnabled,
			Description: "enable",
		},
		ActionDisable: {
			Name:        ActionDisable,
			Description: "disable",
		},
	}
}
