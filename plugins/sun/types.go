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

package sun

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"time"
)

const (
	Name      = "sun"
	EntitySun = common.EntityType("sun")
)

const (
	StateAboveHorizon = "aboveHorizon"
	StateBelowHorizon = "belowHorizon"
)

const (
	AttrHorizonState  = "horizonState" // aboveHorizon|belowHorizon
	AttrPhase         = "phase"
	AttrAzimuth       = "azimuth"
	AttrElevation     = "elevation"
	AttrSunrise       = "sunrise"       // sunrise (top edge of the sun appears on the horizon)
	AttrSunset        = "sunset"        // sunset (sun disappears below the horizon, evening civil twilight starts)
	AttrSunriseEnd    = "sunriseEnd"    // sunrise ends (bottom edge of the sun touches the horizon)
	AttrSunsetStart   = "sunsetStart"   // sunset starts (bottom edge of the sun touches the horizon)
	AttrDawn          = "dawn"          // dawn (morning nautical twilight ends, morning civil twilight starts)
	AttrDusk          = "dusk"          // dusk (evening nautical twilight starts)
	AttrNauticalDawn  = "nauticalDawn"  // nautical dawn (morning nautical twilight starts)
	AttrNauticalDusk  = "nauticalDusk"  // nautical dusk (evening astronomical twilight starts)
	AttrNightEnd      = "nightEnd"      // night ends (morning astronomical twilight starts)
	AttrNight         = "night"         // night starts (dark enough for astronomical observations)
	AttrGoldenHourEnd = "goldenHourEnd" // morning golden hour (soft light, best DayTime for photography) ends
	AttrGoldenHour    = "goldenHour"    // evening golden hour starts
	AttrSolarNoon     = "solarNoon"     // solar noon (sun is in the highest position)
	AttrNadir         = "nadir"         // nadir (darkest moment of the night, sun is in the lowest position)

	AttrLat          = "lat"
	AttrLon          = "lon"
)

type DayTime struct {
	MorningName string
	Time        time.Time
}

type DayTimes []DayTime

// Len ...
func (l DayTimes) Len() int { return len(l) }

// Swap ...
func (l DayTimes) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

// Less ...
func (l DayTimes) Less(i, j int) bool { return l[i].Time.UnixNano() < l[j].Time.UnixNano() }

func NewAttr() m.Attributes {
	return m.Attributes{
		AttrSunrise: {
			Name: AttrSunrise,
			Type: common.AttributeTime,
		},
		AttrSunset: {
			Name: AttrSunset,
			Type: common.AttributeTime,
		},
		AttrSunriseEnd: {
			Name: AttrSunriseEnd,
			Type: common.AttributeTime,
		},
		AttrSunsetStart: {
			Name: AttrSunsetStart,
			Type: common.AttributeTime,
		},
		AttrDawn: {
			Name: AttrDawn,
			Type: common.AttributeTime,
		},
		AttrDusk: {
			Name: AttrDusk,
			Type: common.AttributeTime,
		},
		AttrNauticalDawn: {
			Name: AttrNauticalDawn,
			Type: common.AttributeBool,
		},
		AttrNauticalDusk: {
			Name: AttrNauticalDusk,
			Type: common.AttributeBool,
		},
		AttrNightEnd: {
			Name: AttrNightEnd,
			Type: common.AttributeBool,
		},
		AttrNight: {
			Name: AttrNight,
			Type: common.AttributeBool,
		},
		AttrGoldenHourEnd: {
			Name: AttrGoldenHourEnd,
			Type: common.AttributeBool,
		},
		AttrGoldenHour: {
			Name: AttrGoldenHour,
			Type: common.AttributeBool,
		},
		AttrSolarNoon: {
			Name: AttrSolarNoon,
			Type: common.AttributeBool,
		},
		AttrNadir: {
			Name: AttrNadir,
			Type: common.AttributeBool,
		},
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

const (
	Sunrise = "sunrise" // sunrise (top edge of the sun appears on the horizon)
	Sunset  = "sunset"  // sunset (sun disappears below the horizon, evening civil twilight starts)

	SunriseEnd  = "sunriseEnd"  // sunrise ends (bottom edge of the sun touches the horizon)
	SunsetStart = "sunsetStart" // sunset starts (bottom edge of the sun touches the horizon)

	Dawn = "dawn" // dawn (morning nautical twilight ends, morning civil twilight starts)
	Dusk = "dusk" // dusk (evening nautical twilight starts)

	NauticalDawn = "nauticalDawn" // nautical dawn (morning nautical twilight starts)
	NauticalDusk = "nauticalDusk" // nautical dusk (evening astronomical twilight starts)

	NightEnd = "nightEnd" // night ends (morning astronomical twilight starts)
	Night    = "night"    // night starts (dark enough for astronomical observations)

	GoldenHourEnd = "goldenHourEnd" // morning golden hour (soft light, best DayTime for photography) ends
	GoldenHour    = "goldenHour"    // evening golden hour starts

	SolarNoon = "solarNoon" // solar noon (sun is in the highest position)
	Nadir     = "nadir"     // nadir (darkest moment of the night, sun is in the lowest position)
)

func NewStates() (states map[string]entity_manager.ActorState) {

	states = map[string]entity_manager.ActorState{
		Sunrise: {
			Name:        Sunrise,
			Description: "sunrise (top edge of the sun appears on the horizon)",
		},
		Sunset: {
			Name:        Sunset,
			Description: "sunset (sun disappears below the horizon, evening civil twilight starts)",
		},
		SunriseEnd: {
			Name:        SunriseEnd,
			Description: "sunrise ends (bottom edge of the sun touches the horizon)",
		},
		SunsetStart: {
			Name:        SunsetStart,
			Description: "sunset starts (bottom edge of the sun touches the horizon)",
		},
		Dawn: {
			Name:        Dawn,
			Description: "dawn (morning nautical twilight ends, morning civil twilight starts)",
		},
		Dusk: {
			Name:        Dusk,
			Description: "dusk (evening nautical twilight starts)",
		},
		NauticalDawn: {
			Name:        NauticalDawn,
			Description: "nautical dawn (morning nautical twilight starts)",
		},
		NauticalDusk: {
			Name:        NauticalDusk,
			Description: "nautical dusk (evening astronomical twilight starts)",
		},
		NightEnd: {
			Name:        NightEnd,
			Description: "night ends (morning astronomical twilight starts)",
		},
		Night: {
			Name:        Night,
			Description: "night starts (dark enough for astronomical observations)",
		},
		GoldenHourEnd: {
			Name:        GoldenHourEnd,
			Description: "morning golden hour (soft light, best DayTime for photography) ends",
		},
		GoldenHour: {
			Name:        GoldenHour,
			Description: "evening golden hour starts",
		},
		SolarNoon: {
			Name:        SolarNoon,
			Description: "solar noon (sun is in the highest position)",
		},
		Nadir: {
			Name:        Nadir,
			Description: "nadir (darkest moment of the night, sun is in the lowest position)",
		},
	}

	return
}
