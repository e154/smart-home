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

package speedtest

import (
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

const (
	Name = "speedtest"

	AttrPoint = "location"
	AttrCity  = "city"

	AttrDLSpeed  = "DLSpeed"
	AttrULSpeed  = "ULSpeed"
	AttrLatency  = "Latency"
	AttrName     = "Name"
	AttrCountry  = "Country"
	AttrSponsor  = "Sponsor"
	AttrDistance = "Distance"
	AttrJitter   = "Jitter"

	StateInProcess = "in process"
	StateCompleted = "completed"

	ActionCheck = "check"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrPoint: {
			Name: AttrPoint,
			Type: common.AttributePoint,
		},
		AttrName: {
			Name: AttrName,
			Type: common.AttributeString,
		},
		AttrCountry: {
			Name: AttrCountry,
			Type: common.AttributeString,
		},
		AttrSponsor: {
			Name: AttrSponsor,
			Type: common.AttributeString,
		},
		AttrDistance: {
			Name: AttrDistance,
			Type: common.AttributeFloat,
		},
		AttrLatency: {
			Name: AttrLatency,
			Type: common.AttributeString,
		},
		AttrJitter: {
			Name: AttrJitter,
			Type: common.AttributeString,
		},
		AttrDLSpeed: {
			Name: AttrDLSpeed,
			Type: common.AttributeFloat,
		},
		AttrULSpeed: {
			Name: AttrULSpeed,
			Type: common.AttributeFloat,
		},
	}
}

// NewSettings ...
func NewSettings() m.Attributes {
	return m.Attributes{
		AttrCity: {
			Name: AttrCity,
			Type: common.AttributeString,
		},
		AttrPoint: {
			Name: AttrPoint,
			Type: common.AttributePoint,
		},
	}
}

// NewStates ...
func NewStates() (states map[string]plugins.ActorState) {

	states = map[string]plugins.ActorState{
		StateInProcess: {
			Name:        StateInProcess,
			Description: "test in process",
		},
		StateCompleted: {
			Name:        StateCompleted,
			Description: "test completed",
		},
	}

	return
}

func NewActions() map[string]plugins.ActorAction {
	return map[string]plugins.ActorAction{
		ActionCheck: {
			Name:        ActionCheck,
			Description: "run test",
		},
	}
}
