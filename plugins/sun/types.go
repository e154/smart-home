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
	"time"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
)

const (
	// Name ...
	Name = "sun"
	// EntitySun ...
	EntitySun = string("sun")
)

const (
	// StateAboveHorizon ...
	StateAboveHorizon = "aboveHorizon"
	// StateBelowHorizon ...
	StateBelowHorizon = "belowHorizon"
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
	// AttrSunrise ...
	AttrSunrise = "sunrise"
	// AttrSunset ...
	AttrSunset = "sunset"
	// AttrSunriseEnd ...
	AttrSunriseEnd = "sunriseEnd"
	// AttrSunsetStart ...
	AttrSunsetStart = "sunsetStart"
	// AttrDawn ...
	AttrDawn = "dawn"
	// AttrDusk ...
	AttrDusk = "dusk"
	// AttrNauticalDawn ...
	AttrNauticalDawn = "nauticalDawn"
	// AttrNauticalDusk ...
	AttrNauticalDusk = "nauticalDusk"
	// AttrNightEnd ...
	AttrNightEnd = "nightEnd"
	// AttrNight ...
	AttrNight = "night"
	// AttrGoldenHourEnd ...
	AttrGoldenHourEnd = "goldenHourEnd"
	// AttrGoldenHour ...
	AttrGoldenHour = "goldenHour"
	// AttrSolarNoon ...
	AttrSolarNoon = "solarNoon"
	// AttrNadir ...
	AttrNadir = "nadir"

	// AttrLat ...
	AttrLat = "lat"
	// AttrLon ...
	AttrLon = "lon"
)

// DayTime ...
type DayTime struct {
	MorningName string
	Time        time.Time
}

// DayTimes ...
type DayTimes []DayTime

// Len ...
func (l DayTimes) Len() int { return len(l) }

// Swap ...
func (l DayTimes) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

// Less ...
func (l DayTimes) Less(i, j int) bool { return l[i].Time.UnixNano() < l[j].Time.UnixNano() }

// NewAttr ...
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
			Type: common.AttributeTime,
		},
		AttrNauticalDusk: {
			Name: AttrNauticalDusk,
			Type: common.AttributeTime,
		},
		AttrNightEnd: {
			Name: AttrNightEnd,
			Type: common.AttributeTime,
		},
		AttrNight: {
			Name: AttrNight,
			Type: common.AttributeTime,
		},
		AttrGoldenHourEnd: {
			Name: AttrGoldenHourEnd,
			Type: common.AttributeTime,
		},
		AttrGoldenHour: {
			Name: AttrGoldenHour,
			Type: common.AttributeTime,
		},
		AttrSolarNoon: {
			Name: AttrSolarNoon,
			Type: common.AttributeTime,
		},
		AttrNadir: {
			Name: AttrNadir,
			Type: common.AttributeTime,
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

const (
	// Sunrise ...
	Sunrise = "sunrise"
	// Sunset ...
	Sunset = "sunset"

	// SunriseEnd ...
	SunriseEnd = "sunriseEnd"
	// SunsetStart ...
	SunsetStart = "sunsetStart"

	// Dawn ...
	Dawn = "dawn"
	// Dusk ...
	Dusk = "dusk"

	// NauticalDawn ...
	NauticalDawn = "nauticalDawn"
	// NauticalDusk ...
	NauticalDusk = "nauticalDusk"

	// NightEnd ...
	NightEnd = "nightEnd"
	// Night ...
	Night = "night"

	// GoldenHourEnd ...
	GoldenHourEnd = "goldenHourEnd"
	// GoldenHour ...
	GoldenHour = "goldenHour"

	// SolarNoon ...
	SolarNoon = "solarNoon"
	// Nadir ...
	Nadir = "nadir"
)

// NewStates ...
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
