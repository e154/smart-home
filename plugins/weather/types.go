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

package weather

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	Version = "0.0.1"

	AttrForecast     = "forecast"
	AttrForecastDay1 = "forecast_day1"
	AttrForecastDay2 = "forecast_day2"
	AttrForecastDay3 = "forecast_day3"
	AttrForecastDay4 = "forecast_day4"
	AttrForecastDay5 = "forecast_day5"

	AttrWeatherCondition   = "condition"
	AttrWeatherOzone       = "ozone"
	AttrWeatherAttribution = "attribution"
	AttrWeatherTemperature = "temperature"
	AttrWeatherVisibility  = "visibility"
	AttrWeatherDescription = "description"
	AttrWeatherIcon        = "icon"

	AttrWeatherMain               = "main"
	AttrWeatherDatetime           = "datetime"
	AttrWeatherHumidity           = "humidity"
	AttrWeatherMaxTemperature     = "max_temperature"
	AttrWeatherMinTemperature     = "min_temperature"
	AttrWeatherPressure           = "pressure"
	AttrWeatherWindBearing        = "wind_bearing"
	AttrWeatherWindSpeed          = "wind_speed"
	Day1AttrWeatherMain           = "day1_main"
	Dya1AttrWeatherIcon           = "day1_icon"
	Day1AttrWeatherDatetime       = "day1_datetime"
	Day1AttrWeatherHumidity       = "day1_humidity"
	Day1AttrWeatherMaxTemperature = "day1_max_temperature"
	Day1AttrWeatherMinTemperature = "day1_min_temperature"
	Day1AttrWeatherPressure       = "day1_pressure"
	Day1AttrWeatherWindBearing    = "day1_wind_bearing"
	Day1AttrWeatherWindSpeed      = "day1_wind_speed"
	Day2AttrWeatherMain           = "day2_main"
	Dya2AttrWeatherIcon           = "day2_icon"
	Day2AttrWeatherDatetime       = "day2_datetime"
	Day2AttrWeatherHumidity       = "day2_humidity"
	Day2AttrWeatherMaxTemperature = "day2_max_temperature"
	Day2AttrWeatherMinTemperature = "day2_min_temperature"
	Day2AttrWeatherPressure       = "day2_pressure"
	Day2AttrWeatherWindBearing    = "day2_wind_bearing"
	Day2AttrWeatherWindSpeed      = "day2_wind_speed"
	Day3AttrWeatherMain           = "day3_main"
	Dya3AttrWeatherIcon           = "day3_icon"
	Day3AttrWeatherDatetime       = "day3_datetime"
	Day3AttrWeatherHumidity       = "day3_humidity"
	Day3AttrWeatherMaxTemperature = "day3_max_temperature"
	Day3AttrWeatherMinTemperature = "day3_min_temperature"
	Day3AttrWeatherPressure       = "day3_pressure"
	Day3AttrWeatherWindBearing    = "day3_wind_bearing"
	Day3AttrWeatherWindSpeed      = "day3_wind_speed"
	Day4AttrWeatherMain           = "day4_main"
	Dya4AttrWeatherIcon           = "day4_icon"
	Day4AttrWeatherDatetime       = "day4_datetime"
	Day4AttrWeatherHumidity       = "day4_humidity"
	Day4AttrWeatherMaxTemperature = "day4_max_temperature"
	Day4AttrWeatherMinTemperature = "day4_min_temperature"
	Day4AttrWeatherPressure       = "day4_pressure"
	Day4AttrWeatherWindBearing    = "day4_wind_bearing"
	Day4AttrWeatherWindSpeed      = "day4_wind_speed"
	Day5AttrWeatherMain           = "day5_main"
	Dya5AttrWeatherIcon           = "day5_icon"
	Day5AttrWeatherDatetime       = "day5_datetime"
	Day5AttrWeatherHumidity       = "day5_humidity"
	Day5AttrWeatherMaxTemperature = "day5_max_temperature"
	Day5AttrWeatherMinTemperature = "day5_min_temperature"
	Day5AttrWeatherPressure       = "day5_pressure"
	Day5AttrWeatherWindBearing    = "day5_wind_bearing"
	Day5AttrWeatherWindSpeed      = "day5_wind_speed"

	// Attribution ...
	Attribution = ""

	// AttrLat ...
	AttrLat = "lat"
	// AttrLon ...
	AttrLon = "lon"
	// AttrTheme ...
	AttrTheme = "theme"
	// AttrWinter ...
	AttrWinter = "winter"
)

// BaseForecast ...
func BaseForecast() m.Attributes {
	return NewForecast(6)
}

// NewForecast ...
func NewForecast(days int) m.Attributes {
	attributes := NewAttr()

	attributes[AttrWeatherAttribution] = &m.Attribute{
		Name:  AttrWeatherAttribution,
		Type:  common.AttributeString,
		Value: Attribution,
	}

	for i := 1; i < days; i++ {
		attrs := NewAttr()
		for name, attr := range attrs {
			newName := fmt.Sprintf("day%d_%s", i, name)
			attr.Name = newName
			attributes[newName] = attr
		}
	}

	return attributes
}

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrWeatherDatetime: {
			Name: AttrWeatherDatetime,
			Type: common.AttributeTime,
		},
		AttrWeatherTemperature: {
			Name: AttrWeatherTemperature,
			Type: common.AttributeFloat,
		},
		AttrWeatherMain: {
			Name: AttrWeatherMain,
			Type: common.AttributeString,
		},
		AttrWeatherDescription: {
			Name: AttrWeatherDescription,
			Type: common.AttributeString,
		},
		AttrWeatherIcon: {
			Name: AttrWeatherIcon,
			Type: common.AttributeImage,
		},
		AttrWeatherMinTemperature: {
			Name: AttrWeatherMinTemperature,
			Type: common.AttributeFloat,
		},
		AttrWeatherMaxTemperature: {
			Name: AttrWeatherMaxTemperature,
			Type: common.AttributeFloat,
		},
		AttrWeatherHumidity: {
			Name: AttrWeatherHumidity,
			Type: common.AttributeFloat,
		},
		AttrWeatherPressure: {
			Name: AttrWeatherPressure,
			Type: common.AttributeFloat,
		},
		AttrWeatherWindBearing: {
			Name: AttrWeatherWindBearing,
			Type: common.AttributeFloat,
		},
		AttrWeatherWindSpeed: {
			Name: AttrWeatherWindSpeed,
			Type: common.AttributeFloat,
		},
		AttrWeatherVisibility: {
			Name: AttrWeatherVisibility,
			Type: common.AttributeInt,
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
		AttrTheme: {
			Name: AttrTheme,
			Type: common.AttributeString,
		},
		AttrWinter: {
			Name: AttrWinter,
			Type: common.AttributeBool,
		},
	}
}

// EventStateChanged ...
type EventStateChanged struct {
	Type       string          `json:"type"`
	EntityId   common.EntityId `json:"entity_id"`
	State      string          `json:"state"`
	Attributes m.Attributes    `json:"attributes"`
	Settings   m.Attributes    `json:"settings"`
}

const (
	// StateClearSky ...
	StateClearSky = "clearSky"
	// StateFair ...
	StateFair = "fair"
	// StatePartlyCloudy ...
	StatePartlyCloudy = "partlyCloudy"
	// StateCloudy ...
	StateCloudy = "cloudy"
	// StateLightRainShowers ...
	StateLightRainShowers = "lightRainShowers"
	// StateRainShowers ...
	StateRainShowers = "rainShowers"
	// StateHeavyRainShowers ...
	StateHeavyRainShowers = "heavyRainShowers"
	// StateLightRainShowersAndThunder ...
	StateLightRainShowersAndThunder = "lightRainShowersAndThunder"
	// StateRainShowersAndThunder ...
	StateRainShowersAndThunder = "rainShowersAndThunder"
	// StateHeavyRainShowersAndThunder ...
	StateHeavyRainShowersAndThunder = "heavyRainShowersAndThunder"
	// StateLightSleetShowers ...
	StateLightSleetShowers = "lightSleetShowers"
	// StateSleetShowers ...
	StateSleetShowers = "sleetShowers"
	// StateHeavySleetShowers ...
	StateHeavySleetShowers = "heavySleetShowers"
	// StateLightSleetShowersAndThunder ...
	StateLightSleetShowersAndThunder = "lightSleetShowersAndThunder"
	// StateSleetShowersAndThunder ...
	StateSleetShowersAndThunder = "sleetShowersAndThunder"
	// StateHeavySleetShowersAndThunder ...
	StateHeavySleetShowersAndThunder = "heavySleetShowersAndThunder"
	// StateLightSnowShowers ...
	StateLightSnowShowers = "lightSnowShowers"
	// StateSnowShowers ...
	StateSnowShowers = "snowShowers"
	// StateHeavySnowShowers ...
	StateHeavySnowShowers = "heavySnowShowers"
	// StateLightSnowShowersAndThunder ...
	StateLightSnowShowersAndThunder = "lightSnowShowersAndThunder"
	// StateSnowShowersAndThunder ...
	StateSnowShowersAndThunder = "snowShowersAndThunder"
	// StateHeavySnowShowersAndThunder ...
	StateHeavySnowShowersAndThunder = "heavySnowShowersAndThunder"
	// StateLightRain ...
	StateLightRain = "lightRain"
	// StateRain ...
	StateRain = "rain"
	// StateHeavyRain ...
	StateHeavyRain = "heavyRain"
	// StateLightRainAndThunder ...
	StateLightRainAndThunder = "lightRainAndThunder"
	// StateRainAndThunder ...
	StateRainAndThunder = "rainAndThunder"
	// StateHeavyRainAndThunder ...
	StateHeavyRainAndThunder = "heavyRainAndThunder"
	// StateLightSleet ...
	StateLightSleet = "lightSleet"
	// StateSleet ...
	StateSleet = "sleet"
	// StateHeavySleet ...
	StateHeavySleet = "heavySleet"
	// StateLightSleetAndThunder ...
	StateLightSleetAndThunder = "lightSleetAndThunder"
	// StateSleetAndThunder ...
	StateSleetAndThunder = "sleetAndThunder"
	// StateHeavySleetAndThunder ...
	StateHeavySleetAndThunder = "heavySleetAndThunder"
	// StateLightSnow ...
	StateLightSnow = "lightSnow"
	// StateSnow ...
	StateSnow = "snow"
	// StateHeavySnow ...
	StateHeavySnow = "heavySnow"
	// StateLightSnowAndThunder ...
	StateLightSnowAndThunder = "lightSnowAndThunder"
	// StateSnowAndThunder ...
	StateSnowAndThunder = "snowAndThunder"
	// StateHeavySnowAndThunder ...
	StateHeavySnowAndThunder = "heavySnowAndThunder"
	// StateFog ...
	StateFog = "fog"
)

var (
	conditions = map[string]int{
		StateClearSky:                    1,
		StateFair:                        2,
		StatePartlyCloudy:                3,
		StateCloudy:                      4,
		StateRainShowers:                 5,
		StateRainShowersAndThunder:       6,
		StateSleetShowers:                7,
		StateSnowShowers:                 8,
		StateRain:                        9,
		StateHeavyRain:                   10,
		StateHeavyRainAndThunder:         11,
		StateSleet:                       12,
		StateSnow:                        13,
		StateSnowAndThunder:              14,
		StateFog:                         15,
		StateSleetShowersAndThunder:      20,
		StateSnowShowersAndThunder:       21,
		StateRainAndThunder:              22,
		StateSleetAndThunder:             23,
		StateLightRainShowersAndThunder:  24,
		StateHeavyRainShowersAndThunder:  25,
		StateLightSleetShowersAndThunder: 26,
		StateHeavySleetShowersAndThunder: 27,
		StateLightSnowShowersAndThunder:  28,
		StateHeavySnowShowersAndThunder:  29,
		StateLightRainAndThunder:         30,
		StateLightSleetAndThunder:        31,
		StateHeavySleetAndThunder:        32,
		StateLightSnowAndThunder:         33,
		StateHeavySnowAndThunder:         34,
		StateLightRainShowers:            40,
		StateHeavyRainShowers:            41,
		StateLightSleetShowers:           42,
		StateHeavySleetShowers:           43,
		StateLightSnowShowers:            44,
		StateHeavySnowShowers:            45,
		StateLightRain:                   46,
		StateLightSleet:                  47,
		StateHeavySleet:                  48,
		StateLightSnow:                   49,
		StateHeavySnow:                   50,
	}
)

// GetActorState ...
func GetActorState(state, theme string, n, w bool) (actorState supervisor.ActorState) {
	switch state {
	case StateClearSky:
		actorState = supervisor.ActorState{Name: StateClearSky, Description: "clear sky", ImageUrl: GetImagePath(StateClearSky, theme, n, w)}
	case StateFair:
		actorState = supervisor.ActorState{Name: StateFair, Description: "fair", ImageUrl: GetImagePath(StateFair, theme, n, w)}
	case StatePartlyCloudy:
		actorState = supervisor.ActorState{Name: StatePartlyCloudy, Description: "partly cloudy", ImageUrl: GetImagePath(StatePartlyCloudy, theme, n, w)}
	case StateCloudy:
		actorState = supervisor.ActorState{Name: StateCloudy, Description: "cloudy", ImageUrl: GetImagePath(StateCloudy, theme, n, w)}
	case StateRainShowers:
		actorState = supervisor.ActorState{Name: StateRainShowers, Description: "rain showers", ImageUrl: GetImagePath(StateRainShowers, theme, n, w)}
	case StateRainShowersAndThunder:
		actorState = supervisor.ActorState{Name: StateRainShowersAndThunder, Description: "rain showers and thunder", ImageUrl: GetImagePath(StateRainShowersAndThunder, theme, n, w)}
	case StateSleetShowers:
		actorState = supervisor.ActorState{Name: StateSleetShowers, Description: "sleet showers", ImageUrl: GetImagePath(StateSleetShowers, theme, n, w)}
	case StateSnowShowers:
		actorState = supervisor.ActorState{Name: StateSnowShowers, Description: "snow showers", ImageUrl: GetImagePath(StateSnowShowers, theme, n, w)}
	case StateRain:
		actorState = supervisor.ActorState{Name: StateRain, Description: "rain", ImageUrl: GetImagePath(StateRain, theme, n, w)}
	case StateHeavyRain:
		actorState = supervisor.ActorState{Name: StateHeavyRain, Description: "heavy rain", ImageUrl: GetImagePath(StateHeavyRain, theme, n, w)}
	case StateHeavyRainAndThunder:
		actorState = supervisor.ActorState{Name: StateHeavyRainAndThunder, Description: "heavy rain and thunder", ImageUrl: GetImagePath(StateHeavyRainAndThunder, theme, n, w)}
	case StateSleet:
		actorState = supervisor.ActorState{Name: StateSleet, Description: "sleet", ImageUrl: GetImagePath(StateSleet, theme, n, w)}
	case StateSnow:
		actorState = supervisor.ActorState{Name: StateSnow, Description: "snow", ImageUrl: GetImagePath(StateSnow, theme, n, w)}
	case StateSnowAndThunder:
		actorState = supervisor.ActorState{Name: StateSnowAndThunder, Description: "snow and thunder", ImageUrl: GetImagePath(StateSnowAndThunder, theme, n, w)}
	case StateFog:
		actorState = supervisor.ActorState{Name: StateFog, Description: "fog", ImageUrl: GetImagePath(StateFog, theme, n, w)}
	case StateSleetShowersAndThunder:
		actorState = supervisor.ActorState{Name: StateSleetShowersAndThunder, Description: "sleet showers and thunder", ImageUrl: GetImagePath(StateSleetShowersAndThunder, theme, n, w)}
	case StateSnowShowersAndThunder:
		actorState = supervisor.ActorState{Name: StateSnowShowersAndThunder, Description: "snow showers and thunder", ImageUrl: GetImagePath(StateSnowShowersAndThunder, theme, n, w)}
	case StateRainAndThunder:
		actorState = supervisor.ActorState{Name: StateRainAndThunder, Description: "rain and thunder", ImageUrl: GetImagePath(StateRainAndThunder, theme, n, w)}
	case StateSleetAndThunder:
		actorState = supervisor.ActorState{Name: StateSleetAndThunder, Description: "sleet and thunder", ImageUrl: GetImagePath(StateSleetAndThunder, theme, n, w)}
	case StateLightRainShowersAndThunder:
		actorState = supervisor.ActorState{Name: StateLightRainShowersAndThunder, Description: "light rain showers and thunder", ImageUrl: GetImagePath(StateLightRainShowersAndThunder, theme, n, w)}
	case StateHeavyRainShowersAndThunder:
		actorState = supervisor.ActorState{Name: StateHeavyRainShowersAndThunder, Description: "heavy rain showers and thunder", ImageUrl: GetImagePath(StateHeavyRainShowersAndThunder, theme, n, w)}
	case StateLightSleetShowersAndThunder:
		actorState = supervisor.ActorState{Name: StateLightSleetShowersAndThunder, Description: "light sleet showers and thunder", ImageUrl: GetImagePath(StateLightSleetShowersAndThunder, theme, n, w)}
	case StateHeavySleetShowersAndThunder:
		actorState = supervisor.ActorState{Name: StateHeavySleetShowersAndThunder, Description: "heavy sleet showers and thunder", ImageUrl: GetImagePath(StateHeavySleetShowersAndThunder, theme, n, w)}
	case StateLightSnowShowersAndThunder:
		actorState = supervisor.ActorState{Name: StateLightSnowShowersAndThunder, Description: "light snow showers and thunder", ImageUrl: GetImagePath(StateLightSnowShowersAndThunder, theme, n, w)}
	case StateHeavySnowShowersAndThunder:
		actorState = supervisor.ActorState{Name: StateHeavySnowShowersAndThunder, Description: "heavy snow showers and thunder", ImageUrl: GetImagePath(StateHeavySnowShowersAndThunder, theme, n, w)}
	case StateLightRainAndThunder:
		actorState = supervisor.ActorState{Name: StateLightRainAndThunder, Description: "light rain and thunder", ImageUrl: GetImagePath(StateLightRainAndThunder, theme, n, w)}
	case StateLightSleetAndThunder:
		actorState = supervisor.ActorState{Name: StateLightSleetAndThunder, Description: "light sleet and thunder", ImageUrl: GetImagePath(StateLightSleetAndThunder, theme, n, w)}
	case StateHeavySleetAndThunder:
		actorState = supervisor.ActorState{Name: StateHeavySleetAndThunder, Description: "heavy sleet and thunder", ImageUrl: GetImagePath(StateHeavySleetAndThunder, theme, n, w)}
	case StateLightSnowAndThunder:
		actorState = supervisor.ActorState{Name: StateLightSnowAndThunder, Description: "light snow and thunder", ImageUrl: GetImagePath(StateLightSnowAndThunder, theme, n, w)}
	case StateHeavySnowAndThunder:
		actorState = supervisor.ActorState{Name: StateHeavySnowAndThunder, Description: "heavy snow and thunder", ImageUrl: GetImagePath(StateHeavySnowAndThunder, theme, n, w)}
	case StateLightRainShowers:
		actorState = supervisor.ActorState{Name: StateLightRainShowers, Description: "light rain showers", ImageUrl: GetImagePath(StateLightRainShowers, theme, n, w)}
	case StateHeavyRainShowers:
		actorState = supervisor.ActorState{Name: StateHeavyRainShowers, Description: "heavy rain showers", ImageUrl: GetImagePath(StateHeavyRainShowers, theme, n, w)}
	case StateLightSleetShowers:
		actorState = supervisor.ActorState{Name: StateLightSleetShowers, Description: "light sleet showers", ImageUrl: GetImagePath(StateLightSleetShowers, theme, n, w)}
	case StateHeavySleetShowers:
		actorState = supervisor.ActorState{Name: StateHeavySleetShowers, Description: "heavy sleet showers", ImageUrl: GetImagePath(StateHeavySleetShowers, theme, n, w)}
	case StateLightSnowShowers:
		actorState = supervisor.ActorState{Name: StateLightSnowShowers, Description: "light snow showers", ImageUrl: GetImagePath(StateLightSnowShowers, theme, n, w)}
	case StateHeavySnowShowers:
		actorState = supervisor.ActorState{Name: StateHeavySnowShowers, Description: "heavy snow showers", ImageUrl: GetImagePath(StateHeavySnowShowers, theme, n, w)}
	case StateLightRain:
		actorState = supervisor.ActorState{Name: StateLightRain, Description: "light rain", ImageUrl: GetImagePath(StateLightRain, theme, n, w)}
	case StateLightSleet:
		actorState = supervisor.ActorState{Name: StateLightSleet, Description: "light sleet", ImageUrl: GetImagePath(StateLightSleet, theme, n, w)}
	case StateHeavySleet:
		actorState = supervisor.ActorState{Name: StateHeavySleet, Description: "heavy sleet", ImageUrl: GetImagePath(StateHeavySleet, theme, n, w)}
	case StateLightSnow:
		actorState = supervisor.ActorState{Name: StateLightSnow, Description: "light snow", ImageUrl: GetImagePath(StateLightSnow, theme, n, w)}
	case StateHeavySnow:
		actorState = supervisor.ActorState{Name: StateHeavySnow, Description: "heavy snow", ImageUrl: GetImagePath(StateHeavySnow, theme, n, w)}
	}

	return
}

// NewActorStates ...
func NewActorStates(n, w bool) (states map[string]supervisor.ActorState) {
	const theme = "yr"
	states = map[string]supervisor.ActorState{
		StateClearSky:                    GetActorState(StateClearSky, theme, n, w),
		StateFair:                        GetActorState(StateFair, theme, n, w),
		StatePartlyCloudy:                GetActorState(StatePartlyCloudy, theme, n, w),
		StateCloudy:                      GetActorState(StateCloudy, theme, n, w),
		StateRainShowers:                 GetActorState(StateRainShowers, theme, n, w),
		StateRainShowersAndThunder:       GetActorState(StateRainShowersAndThunder, theme, n, w),
		StateSleetShowers:                GetActorState(StateSleetShowers, theme, n, w),
		StateSnowShowers:                 GetActorState(StateSnowShowers, theme, n, w),
		StateRain:                        GetActorState(StateRain, theme, n, w),
		StateHeavyRain:                   GetActorState(StateHeavyRain, theme, n, w),
		StateHeavyRainAndThunder:         GetActorState(StateHeavyRainAndThunder, theme, n, w),
		StateSleet:                       GetActorState(StateSleet, theme, n, w),
		StateSnow:                        GetActorState(StateSnow, theme, n, w),
		StateSnowAndThunder:              GetActorState(StateSnowAndThunder, theme, n, w),
		StateFog:                         GetActorState(StateFog, theme, n, w),
		StateSleetShowersAndThunder:      GetActorState(StateSleetShowersAndThunder, theme, n, w),
		StateSnowShowersAndThunder:       GetActorState(StateSnowShowersAndThunder, theme, n, w),
		StateRainAndThunder:              GetActorState(StateRainAndThunder, theme, n, w),
		StateSleetAndThunder:             GetActorState(StateSleetAndThunder, theme, n, w),
		StateLightRainShowersAndThunder:  GetActorState(StateLightRainShowersAndThunder, theme, n, w),
		StateHeavyRainShowersAndThunder:  GetActorState(StateHeavyRainShowersAndThunder, theme, n, w),
		StateLightSleetShowersAndThunder: GetActorState(StateLightSleetShowersAndThunder, theme, n, w),
		StateHeavySleetShowersAndThunder: GetActorState(StateHeavySleetShowersAndThunder, theme, n, w),
		StateLightSnowShowersAndThunder:  GetActorState(StateLightSnowShowersAndThunder, theme, n, w),
		StateHeavySnowShowersAndThunder:  GetActorState(StateHeavySnowShowersAndThunder, theme, n, w),
		StateLightRainAndThunder:         GetActorState(StateLightRainAndThunder, theme, n, w),
		StateLightSleetAndThunder:        GetActorState(StateLightSleetAndThunder, theme, n, w),
		StateHeavySleetAndThunder:        GetActorState(StateHeavySleetAndThunder, theme, n, w),
		StateLightSnowAndThunder:         GetActorState(StateLightSnowAndThunder, theme, n, w),
		StateHeavySnowAndThunder:         GetActorState(StateHeavySnowAndThunder, theme, n, w),
		StateLightRainShowers:            GetActorState(StateLightRainShowers, theme, n, w),
		StateHeavyRainShowers:            GetActorState(StateHeavyRainShowers, theme, n, w),
		StateLightSleetShowers:           GetActorState(StateLightSleetShowers, theme, n, w),
		StateHeavySleetShowers:           GetActorState(StateHeavySleetShowers, theme, n, w),
		StateLightSnowShowers:            GetActorState(StateLightSnowShowers, theme, n, w),
		StateHeavySnowShowers:            GetActorState(StateHeavySnowShowers, theme, n, w),
		StateLightRain:                   GetActorState(StateLightRain, theme, n, w),
		StateLightSleet:                  GetActorState(StateLightSleet, theme, n, w),
		StateHeavySleet:                  GetActorState(StateHeavySleet, theme, n, w),
		StateLightSnow:                   GetActorState(StateLightSnow, theme, n, w),
		StateHeavySnow:                   GetActorState(StateHeavySnow, theme, n, w),
	}

	return
}

// GetIndexByState ...
func GetIndexByState(state string) (idx int) {
	idx = conditions[state]
	return
}

// GetStateByIndex ...
func GetStateByIndex(idx int) (state string) {
	for s, i := range conditions {
		if idx == i {
			state = s
			return
		}
	}
	return
}

// GetImagePath ...
func GetImagePath(state, theme string, night, winter bool) (imagePath *string) {
	if theme == "" {
		theme = "yr"
	}
	idx := GetIndexByState(state)
	var sfx = "d"
	switch {
	case winter:
		sfx = "m"
	case night:
		sfx = "n"
	}

	switch idx {
	case 4, 46, 9, 10, 30, 22, 11, 47, 12, 48, 31, 23, 32, 49, 13, 50, 33, 14, 34, 15:
		sfx = ""
	}

LOOP:
	// https://github.com/nrkno/yr-weather-symbols
	p := path.Join("static", "weather", theme, fmt.Sprintf("%02d%s.svg", idx, sfx))
	if !common.FileExist(path.Join("data", p)) {
		if sfx == "m" || sfx == "n" {
			sfx = "d"
			goto LOOP
		}
		if sfx == "d" {
			sfx = ""
			goto LOOP
		}
	}
	p = path.Join(string(filepath.Separator), p)
	imagePath = common.String(p)
	return
}
