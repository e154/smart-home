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

package moon

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
)

const (
	// Name ...
	Name = "moon"
	// EntityMoon ...
	EntityMoon = string("moon")
)

const (
	// StateAboveHorizon ...
	StateAboveHorizon = "aboveHorizon"
	// StateBelowHorizon ...
	StateBelowHorizon = "belowHorizon"
)

const (
	// StateNewMoon ...
	StateNewMoon = "new_moon"
	// StateWaxingCrescent ...
	StateWaxingCrescent = "waxing_crescent"
	// StateFirstQuarter ...
	StateFirstQuarter = "first_quarter"
	// StateWaxingGibbous ...
	StateWaxingGibbous = "waxing_gibbous"
	// StateFullMoon ...
	StateFullMoon = "full_moon"
	// StateWaningGibbous ...
	StateWaningGibbous = "waning_gibbous"
	// StateThirdQuarter ...
	StateThirdQuarter = "third_quarter"
	// StateWaningCrescent ...
	StateWaningCrescent = "waning_crescent"
)

const (
	// AttrHorizonState ...
	AttrHorizonState = "horizonState"
	// AttrPhase ...
	AttrPhase = "phase"
	// AttrAzimuth ...
	AttrAzimuth = "azimuth"
	// AttrElevation ...
	AttrElevation = "elevation"
	// AttrLat ...
	AttrLat = "lat"
	// AttrLon ...
	AttrLon = "lon"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{

		AttrElevation: {
			Name: AttrElevation,
			Type: common.AttributeFloat,
		},
		AttrAzimuth: {
			Name: AttrAzimuth,
			Type: common.AttributeFloat,
		},
		AttrPhase: {
			Name: AttrPhase,
			Type: common.AttributeString,
		},
		AttrHorizonState: {
			Name: AttrHorizonState,
			Type: common.AttributeString,
		},
	}
}

// NewSettings ...
func NewSettings() m.Attributes {
	return m.Attributes{
		AttrLat: {
			Name: AttrLat,
			Type: common.AttributeFloat,
		},
		AttrLon: {
			Name: AttrLon,
			Type: common.AttributeFloat,
		},
	}
}

// NewStates ...
func NewStates() (states map[string]entity_manager.ActorState) {

	states = map[string]entity_manager.ActorState{
		StateAboveHorizon: {
			Name:        StateAboveHorizon,
			Description: "above horizon",
		},
		StateBelowHorizon: {
			Name:        StateBelowHorizon,
			Description: "below horizon",
		},
	}

	return
}
