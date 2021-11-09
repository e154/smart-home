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

package weather

import (
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"path"
)

const (
	StatePositionUpdate = "positionUpdate"
)

const (
	TopicPluginWeather = string("plugin.weather")
)

const (
	AttrForecast              = "forecast"
	AttrForecastDay1          = "forecast_day1"
	AttrForecastDay2          = "forecast_day2"
	AttrForecastDay3          = "forecast_day3"
	AttrForecastDay4          = "forecast_day4"
	AttrForecastDay5          = "forecast_day5"
	AttrWeatherCondition      = "condition"
	AttrWeatherDatetime       = "datetime"
	AttrWeatherAttribution    = "attribution"
	AttrWeatherHumidity       = "humidity"
	AttrWeatherOzone          = "ozone"
	AttrWeatherPressure       = "pressure"
	AttrWeatherTemperature    = "temperature"
	AttrWeatherMinTemperature = "min_temperature"
	AttrWeatherMaxTemperature = "max_temperature"
	AttrWeatherVisibility     = "visibility"
	AttrWeatherWindBearing    = "wind_bearing"
	AttrWeatherWindSpeed      = "wind_speed"
	AttrWeatherMain           = "main"
	AttrWeatherDescription    = "description"
	AttrWeatherIcon           = "icon"

	Attribution = ""

	AttrLat    = "lat"
	AttrLon    = "lon"
	AttrPlugin = "plugin"
)

func BaseForecast() m.Attributes {
	return NewForecast(6)
}

func NewForecast(days int) m.Attributes {
	attributes := NewAttr()

	attributes[AttrWeatherAttribution] = &m.Attribute{
		Name:  AttrWeatherAttribution,
		Type:  common.AttributeString,
		Value: Attribution,
	}

	for i := 1; i < days; i++ {
		attributes[fmt.Sprintf("forecast_day%d", i)] = &m.Attribute{
			Name:  fmt.Sprintf("forecast_day%d", i),
			Type:  common.AttributeMap,
			Value: NewAttr(),
		}
	}

	return attributes
}

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
			Type: common.AttributeString,
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
		AttrPlugin: {
			Name: AttrPlugin,
			Type: common.AttributeString,
		},
	}
}

type EventStateChanged struct {
	Type       common.EntityType `json:"type"`
	EntityId   common.EntityId   `json:"entity_id"`
	State      string            `json:"state"`
	Attributes m.Attributes      `json:"attributes"`
	Settings   m.Attributes      `json:"settings"`
}

const (
	StateClearSky                    = "clearSky"
	StateFair                        = "fair"
	StatePartlyCloudy                = "partlyCloudy"
	StateCloudy                      = "cloudy"
	StateLightRainShowers            = "lightRainShowers"
	StateRainShowers                 = "rainShowers"
	StateHeavyRainShowers            = "heavyRainShowers"
	StateLightRainShowersAndThunder  = "lightRainShowersAndThunder"
	StateRainShowersAndThunder       = "rainShowersAndThunder"
	StateHeavyRainShowersAndThunder  = "heavyRainShowersAndThunder"
	StateLightSleetShowers           = "lightSleetShowers"
	StateSleetShowers                = "sleetShowers"
	StateHeavySleetShowers           = "heavySleetShowers"
	StateLightSleetShowersAndThunder = "lightSleetShowersAndThunder"
	StateSleetShowersAndThunder      = "sleetShowersAndThunder"
	StateHeavySleetShowersAndThunder = "heavySleetShowersAndThunder"
	StateLightSnowShowers            = "lightSnowShowers"
	StateSnowShowers                 = "snowShowers"
	StateHeavySnowShowers            = "heavySnowShowers"
	StateLightSnowShowersAndThunder  = "lightSnowShowersAndThunder"
	StateSnowShowersAndThunder       = "snowShowersAndThunder"
	StateHeavySnowShowersAndThunder  = "heavySnowShowersAndThunder"
	StateLightRain                   = "lightRain"
	StateRain                        = "rain"
	StateHeavyRain                   = "heavyRain"
	StateLightRainAndThunder         = "lightRainAndThunder"
	StateRainAndThunder              = "rainAndThunder"
	StateHeavyRainAndThunder         = "heavyRainAndThunder"
	StateLightSleet                  = "lightSleet"
	StateSleet                       = "sleet"
	StateHeavySleet                  = "heavySleet"
	StateLightSleetAndThunder        = "lightSleetAndThunder"
	StateSleetAndThunder             = "sleetAndThunder"
	StateHeavySleetAndThunder        = "heavySleetAndThunder"
	StateLightSnow                   = "lightSnow"
	StateSnow                        = "snow"
	StateHeavySnow                   = "heavySnow"
	StateLightSnowAndThunder         = "lightSnowAndThunder"
	StateSnowAndThunder              = "snowAndThunder"
	StateHeavySnowAndThunder         = "heavySnowAndThunder"
	StateFog                         = "fog"
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

func GetActorState(state string, n, w bool) (actorState entity_manager.ActorState) {
	switch state {
	case StateClearSky:
		actorState = entity_manager.ActorState{Name: StateClearSky, Description: "clear sky", ImageUrl: GetImagePath(StateClearSky, n, w)}
	case StateFair:
		actorState = entity_manager.ActorState{Name: StateFair, Description: "fair", ImageUrl: GetImagePath(StateFair, n, w)}
	case StatePartlyCloudy:
		actorState = entity_manager.ActorState{Name: StatePartlyCloudy, Description: "partly cloudy", ImageUrl: GetImagePath(StatePartlyCloudy, n, w)}
	case StateCloudy:
		actorState = entity_manager.ActorState{Name: StateCloudy, Description: "cloudy", ImageUrl: GetImagePath(StateCloudy, n, w)}
	case StateRainShowers:
		actorState = entity_manager.ActorState{Name: StateRainShowers, Description: "rain showers", ImageUrl: GetImagePath(StateRainShowers, n, w)}
	case StateRainShowersAndThunder:
		actorState = entity_manager.ActorState{Name: StateRainShowersAndThunder, Description: "rain showers and thunder", ImageUrl: GetImagePath(StateRainShowersAndThunder, n, w)}
	case StateSleetShowers:
		actorState = entity_manager.ActorState{Name: StateSleetShowers, Description: "sleet showers", ImageUrl: GetImagePath(StateSleetShowers, n, w)}
	case StateSnowShowers:
		actorState = entity_manager.ActorState{Name: StateSnowShowers, Description: "snow showers", ImageUrl: GetImagePath(StateSnowShowers, n, w)}
	case StateRain:
		actorState = entity_manager.ActorState{Name: StateRain, Description: "rain", ImageUrl: GetImagePath(StateRain, n, w)}
	case StateHeavyRain:
		actorState = entity_manager.ActorState{Name: StateHeavyRain, Description: "heavy rain", ImageUrl: GetImagePath(StateHeavyRain, n, w)}
	case StateHeavyRainAndThunder:
		actorState = entity_manager.ActorState{Name: StateHeavyRainAndThunder, Description: "heavy rain and thunder", ImageUrl: GetImagePath(StateHeavyRainAndThunder, n, w)}
	case StateSleet:
		actorState = entity_manager.ActorState{Name: StateSleet, Description: "sleet", ImageUrl: GetImagePath(StateSleet, n, w)}
	case StateSnow:
		actorState = entity_manager.ActorState{Name: StateSnow, Description: "snow", ImageUrl: GetImagePath(StateSnow, n, w)}
	case StateSnowAndThunder:
		actorState = entity_manager.ActorState{Name: StateSnowAndThunder, Description: "snow and thunder", ImageUrl: GetImagePath(StateSnowAndThunder, n, w)}
	case StateFog:
		actorState = entity_manager.ActorState{Name: StateFog, Description: "fog", ImageUrl: GetImagePath(StateFog, n, w)}
	case StateSleetShowersAndThunder:
		actorState = entity_manager.ActorState{Name: StateSleetShowersAndThunder, Description: "sleet showers and thunder", ImageUrl: GetImagePath(StateSleetShowersAndThunder, n, w)}
	case StateSnowShowersAndThunder:
		actorState = entity_manager.ActorState{Name: StateSnowShowersAndThunder, Description: "snow showers and thunder", ImageUrl: GetImagePath(StateSnowShowersAndThunder, n, w)}
	case StateRainAndThunder:
		actorState = entity_manager.ActorState{Name: StateRainAndThunder, Description: "rain and thunder", ImageUrl: GetImagePath(StateRainAndThunder, n, w)}
	case StateSleetAndThunder:
		actorState = entity_manager.ActorState{Name: StateSleetAndThunder, Description: "sleet and thunder", ImageUrl: GetImagePath(StateSleetAndThunder, n, w)}
	case StateLightRainShowersAndThunder:
		actorState = entity_manager.ActorState{Name: StateLightRainShowersAndThunder, Description: "light rain showers and thunder", ImageUrl: GetImagePath(StateLightRainShowersAndThunder, n, w)}
	case StateHeavyRainShowersAndThunder:
		actorState = entity_manager.ActorState{Name: StateHeavyRainShowersAndThunder, Description: "heavy rain showers and thunder", ImageUrl: GetImagePath(StateHeavyRainShowersAndThunder, n, w)}
	case StateLightSleetShowersAndThunder:
		actorState = entity_manager.ActorState{Name: StateLightSleetShowersAndThunder, Description: "light sleet showers and thunder", ImageUrl: GetImagePath(StateLightSleetShowersAndThunder, n, w)}
	case StateHeavySleetShowersAndThunder:
		actorState = entity_manager.ActorState{Name: StateHeavySleetShowersAndThunder, Description: "heavy sleet showers and thunder", ImageUrl: GetImagePath(StateHeavySleetShowersAndThunder, n, w)}
	case StateLightSnowShowersAndThunder:
		actorState = entity_manager.ActorState{Name: StateLightSnowShowersAndThunder, Description: "light snow showers and thunder", ImageUrl: GetImagePath(StateLightSnowShowersAndThunder, n, w)}
	case StateHeavySnowShowersAndThunder:
		actorState = entity_manager.ActorState{Name: StateHeavySnowShowersAndThunder, Description: "heavy snow showers and thunder", ImageUrl: GetImagePath(StateHeavySnowShowersAndThunder, n, w)}
	case StateLightRainAndThunder:
		actorState = entity_manager.ActorState{Name: StateLightRainAndThunder, Description: "light rain and thunder", ImageUrl: GetImagePath(StateLightRainAndThunder, n, w)}
	case StateLightSleetAndThunder:
		actorState = entity_manager.ActorState{Name: StateLightSleetAndThunder, Description: "light sleet and thunder", ImageUrl: GetImagePath(StateLightSleetAndThunder, n, w)}
	case StateHeavySleetAndThunder:
		actorState = entity_manager.ActorState{Name: StateHeavySleetAndThunder, Description: "heavy sleet and thunder", ImageUrl: GetImagePath(StateHeavySleetAndThunder, n, w)}
	case StateLightSnowAndThunder:
		actorState = entity_manager.ActorState{Name: StateLightSnowAndThunder, Description: "light snow and thunder", ImageUrl: GetImagePath(StateLightSnowAndThunder, n, w)}
	case StateHeavySnowAndThunder:
		actorState = entity_manager.ActorState{Name: StateHeavySnowAndThunder, Description: "heavy snow and thunder", ImageUrl: GetImagePath(StateHeavySnowAndThunder, n, w)}
	case StateLightRainShowers:
		actorState = entity_manager.ActorState{Name: StateLightRainShowers, Description: "light rain showers", ImageUrl: GetImagePath(StateLightRainShowers, n, w)}
	case StateHeavyRainShowers:
		actorState = entity_manager.ActorState{Name: StateHeavyRainShowers, Description: "heavy rain showers", ImageUrl: GetImagePath(StateHeavyRainShowers, n, w)}
	case StateLightSleetShowers:
		actorState = entity_manager.ActorState{Name: StateLightSleetShowers, Description: "light sleet showers", ImageUrl: GetImagePath(StateLightSleetShowers, n, w)}
	case StateHeavySleetShowers:
		actorState = entity_manager.ActorState{Name: StateHeavySleetShowers, Description: "heavy sleet showers", ImageUrl: GetImagePath(StateHeavySleetShowers, n, w)}
	case StateLightSnowShowers:
		actorState = entity_manager.ActorState{Name: StateLightSnowShowers, Description: "light snow showers", ImageUrl: GetImagePath(StateLightSnowShowers, n, w)}
	case StateHeavySnowShowers:
		actorState = entity_manager.ActorState{Name: StateHeavySnowShowers, Description: "heavy snow showers", ImageUrl: GetImagePath(StateHeavySnowShowers, n, w)}
	case StateLightRain:
		actorState = entity_manager.ActorState{Name: StateLightRain, Description: "light rain", ImageUrl: GetImagePath(StateLightRain, n, w)}
	case StateLightSleet:
		actorState = entity_manager.ActorState{Name: StateLightSleet, Description: "light sleet", ImageUrl: GetImagePath(StateLightSleet, n, w)}
	case StateHeavySleet:
		actorState = entity_manager.ActorState{Name: StateHeavySleet, Description: "heavy sleet", ImageUrl: GetImagePath(StateHeavySleet, n, w)}
	case StateLightSnow:
		actorState = entity_manager.ActorState{Name: StateLightSnow, Description: "light snow", ImageUrl: GetImagePath(StateLightSnow, n, w)}
	case StateHeavySnow:
		actorState = entity_manager.ActorState{Name: StateHeavySnow, Description: "heavy snow", ImageUrl: GetImagePath(StateHeavySnow, n, w)}
	}

	return
}

func NewActorStates(n, w bool) (states map[string]entity_manager.ActorState) {

	states = map[string]entity_manager.ActorState{
		StateClearSky:                    GetActorState(StateClearSky, n, w),
		StateFair:                        GetActorState(StateFair, n, w),
		StatePartlyCloudy:                GetActorState(StatePartlyCloudy, n, w),
		StateCloudy:                      GetActorState(StateCloudy, n, w),
		StateRainShowers:                 GetActorState(StateRainShowers, n, w),
		StateRainShowersAndThunder:       GetActorState(StateRainShowersAndThunder, n, w),
		StateSleetShowers:                GetActorState(StateSleetShowers, n, w),
		StateSnowShowers:                 GetActorState(StateSnowShowers, n, w),
		StateRain:                        GetActorState(StateRain, n, w),
		StateHeavyRain:                   GetActorState(StateHeavyRain, n, w),
		StateHeavyRainAndThunder:         GetActorState(StateHeavyRainAndThunder, n, w),
		StateSleet:                       GetActorState(StateSleet, n, w),
		StateSnow:                        GetActorState(StateSnow, n, w),
		StateSnowAndThunder:              GetActorState(StateSnowAndThunder, n, w),
		StateFog:                         GetActorState(StateFog, n, w),
		StateSleetShowersAndThunder:      GetActorState(StateSleetShowersAndThunder, n, w),
		StateSnowShowersAndThunder:       GetActorState(StateSnowShowersAndThunder, n, w),
		StateRainAndThunder:              GetActorState(StateRainAndThunder, n, w),
		StateSleetAndThunder:             GetActorState(StateSleetAndThunder, n, w),
		StateLightRainShowersAndThunder:  GetActorState(StateLightRainShowersAndThunder, n, w),
		StateHeavyRainShowersAndThunder:  GetActorState(StateHeavyRainShowersAndThunder, n, w),
		StateLightSleetShowersAndThunder: GetActorState(StateLightSleetShowersAndThunder, n, w),
		StateHeavySleetShowersAndThunder: GetActorState(StateHeavySleetShowersAndThunder, n, w),
		StateLightSnowShowersAndThunder:  GetActorState(StateLightSnowShowersAndThunder, n, w),
		StateHeavySnowShowersAndThunder:  GetActorState(StateHeavySnowShowersAndThunder, n, w),
		StateLightRainAndThunder:         GetActorState(StateLightRainAndThunder, n, w),
		StateLightSleetAndThunder:        GetActorState(StateLightSleetAndThunder, n, w),
		StateHeavySleetAndThunder:        GetActorState(StateHeavySleetAndThunder, n, w),
		StateLightSnowAndThunder:         GetActorState(StateLightSnowAndThunder, n, w),
		StateHeavySnowAndThunder:         GetActorState(StateHeavySnowAndThunder, n, w),
		StateLightRainShowers:            GetActorState(StateLightRainShowers, n, w),
		StateHeavyRainShowers:            GetActorState(StateHeavyRainShowers, n, w),
		StateLightSleetShowers:           GetActorState(StateLightSleetShowers, n, w),
		StateHeavySleetShowers:           GetActorState(StateHeavySleetShowers, n, w),
		StateLightSnowShowers:            GetActorState(StateLightSnowShowers, n, w),
		StateHeavySnowShowers:            GetActorState(StateHeavySnowShowers, n, w),
		StateLightRain:                   GetActorState(StateLightRain, n, w),
		StateLightSleet:                  GetActorState(StateLightSleet, n, w),
		StateHeavySleet:                  GetActorState(StateHeavySleet, n, w),
		StateLightSnow:                   GetActorState(StateLightSnow, n, w),
		StateHeavySnow:                   GetActorState(StateHeavySnow, n, w),
	}

	return
}

func GetIndexByState(state string) (idx int) {
	idx = conditions[state]
	return
}

func GetStateByIndex(idx int) (state string) {
	for s, i := range conditions {
		if idx == i {
			state = s
			return
		}
	}
	return
}

func GetImagePath(state string, night, winter bool) (imagePath *string) {
	idx := GetIndexByState(state)
	var sfx = "d"
	if winter {
		sfx = "m"
	}
	if night {
		sfx = "n"
	}

	switch idx {
	case 4, 46, 9, 10, 30, 22, 11, 47, 12, 48, 31, 23, 32, 49, 13, 50, 33, 14, 34, 15:
		sfx = ""
	}

LOOP:
	// https://github.com/nrkno/yr-weather-symbols
	p := path.Join(common.StaticPath(), "weather", "yr", fmt.Sprintf("%02d%s.svg", idx, sfx))
	if !common.FileExist(p) {
		if sfx == "m" || sfx == "n" {
			sfx = "d"
			goto LOOP
		}
		if sfx == "d" {
			sfx = ""
			goto LOOP
		}
	}
	imagePath = common.String(p)
	return
}
