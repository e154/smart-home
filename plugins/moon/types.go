// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
)

const (
	Name       = "moon"
	EntityMoon = common.EntityType("moon")
)

const (
	StateAboveHorizon = "aboveHorizon"
	StateBelowHorizon = "belowHorizon"
)

const (
	StateNewMoon        = "new_moon"        // New Moon
	StateWaxingCrescent = "waxing_crescent" // Waxing Crescent
	StateFirstQuarter   = "first_quarter"   // First Quarter
	StateWaxingGibbous  = "waxing_gibbous"  // Waxing Gibbous
	StateFullMoon       = "full_moon"       // Full Moon
	StateWaningGibbous  = "waning_gibbous"  // Waning Gibbous
	StateThirdQuarter   = "third_quarter"   // Third Quarter
	StateWaningCrescent = "waning_crescent" // Waning Crescent
)

const (
	AttrHorizonState = "horizonState" // aboveHorizon|belowHorizon
	AttrPhase        = "phase"
	AttrAzimuth      = "azimuth"
	AttrElevation    = "elevation"
)

func NewAttr()  m.EntityAttributes {
	return m.EntityAttributes{

		AttrElevation: {
			Name: AttrElevation,
			Type: common.EntityAttributeFloat,
		},
		AttrAzimuth: {
			Name: AttrAzimuth,
			Type: common.EntityAttributeFloat,
		},
		AttrPhase: {
			Name: AttrPhase,
			Type: common.EntityAttributeString,
		},
		AttrHorizonState: {
			Name: AttrHorizonState,
			Type: common.EntityAttributeString,
		},
	}
}
