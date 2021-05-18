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

type Zone struct {
	Name      string  `json:"name"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Elevation float64 `json:"elevation"`
}

const (
	SubStatePositionUpdate = "positionUpdate"
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

	Attribution = "Weather forecast from met.no, delivered by the Norwegian Meteorological Institute."
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
		AttrWeatherCondition: {
			Name: AttrWeatherCondition,
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
		AttrWeatherOzone: {
			Name: AttrWeatherOzone,
			Type: common.AttributeFloat,
		},
		AttrWeatherVisibility: {
			Name: AttrWeatherVisibility,
			Type: common.AttributeFloat,
		},
	}
}

type EventSubStateChanged struct {
	Type       common.EntityType `json:"type"`
	EntityId   common.EntityId   `json:"entity_id"`
	State      string            `json:"state"`
	Attributes m.Attributes      `json:"attributes"`
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

func States(n, w bool) (states map[string]entity_manager.ActorState) {

	states = map[string]entity_manager.ActorState{
		StateClearSky:                    {Name: StateClearSky, Description: "clear sky", ImageUrl: GetImagePath(StateClearSky, n, w)},
		StateFair:                        {Name: StateFair, Description: "fair", ImageUrl: GetImagePath(StateFair, n, w)},
		StatePartlyCloudy:                {Name: StatePartlyCloudy, Description: "partly cloudy", ImageUrl: GetImagePath(StatePartlyCloudy, n, w)},
		StateCloudy:                      {Name: StateCloudy, Description: "cloudy", ImageUrl: GetImagePath(StateCloudy, n, w)},
		StateRainShowers:                 {Name: StateRainShowers, Description: "rain showers", ImageUrl: GetImagePath(StateRainShowers, n, w)},
		StateRainShowersAndThunder:       {Name: StateRainShowersAndThunder, Description: "rain showers and thunder", ImageUrl: GetImagePath(StateRainShowersAndThunder, n, w)},
		StateSleetShowers:                {Name: StateSleetShowers, Description: "sleet showers", ImageUrl: GetImagePath(StateSleetShowers, n, w)},
		StateSnowShowers:                 {Name: StateSnowShowers, Description: "snow showers", ImageUrl: GetImagePath(StateSnowShowers, n, w)},
		StateRain:                        {Name: StateRain, Description: "rain", ImageUrl: GetImagePath(StateRain, n, w)},
		StateHeavyRain:                   {Name: StateHeavyRain, Description: "heavy rain", ImageUrl: GetImagePath(StateHeavyRain, n, w)},
		StateHeavyRainAndThunder:         {Name: StateHeavyRainAndThunder, Description: "heavy rain and thunder", ImageUrl: GetImagePath(StateHeavyRainAndThunder, n, w)},
		StateSleet:                       {Name: StateSleet, Description: "sleet", ImageUrl: GetImagePath(StateSleet, n, w)},
		StateSnow:                        {Name: StateSnow, Description: "snow", ImageUrl: GetImagePath(StateSnow, n, w)},
		StateSnowAndThunder:              {Name: StateSnowAndThunder, Description: "snow and thunder", ImageUrl: GetImagePath(StateSnowAndThunder, n, w)},
		StateFog:                         {Name: StateFog, Description: "fog", ImageUrl: GetImagePath(StateFog, n, w)},
		StateSleetShowersAndThunder:      {Name: StateSleetShowersAndThunder, Description: "sleet showers and thunder", ImageUrl: GetImagePath(StateSleetShowersAndThunder, n, w)},
		StateSnowShowersAndThunder:       {Name: StateSnowShowersAndThunder, Description: "snow showers and thunder", ImageUrl: GetImagePath(StateSnowShowersAndThunder, n, w)},
		StateRainAndThunder:              {Name: StateRainAndThunder, Description: "rain and thunder", ImageUrl: GetImagePath(StateRainAndThunder, n, w)},
		StateSleetAndThunder:             {Name: StateSleetAndThunder, Description: "sleet and thunder", ImageUrl: GetImagePath(StateSleetAndThunder, n, w)},
		StateLightRainShowersAndThunder:  {Name: StateLightRainShowersAndThunder, Description: "light rain showers and thunder", ImageUrl: GetImagePath(StateLightRainShowersAndThunder, n, w)},
		StateHeavyRainShowersAndThunder:  {Name: StateHeavyRainShowersAndThunder, Description: "heavy rain showers and thunder", ImageUrl: GetImagePath(StateHeavyRainShowersAndThunder, n, w)},
		StateLightSleetShowersAndThunder: {Name: StateLightSleetShowersAndThunder, Description: "light sleet showers and thunder", ImageUrl: GetImagePath(StateLightSleetShowersAndThunder, n, w)},
		StateHeavySleetShowersAndThunder: {Name: StateHeavySleetShowersAndThunder, Description: "heavy sleet showers and thunder", ImageUrl: GetImagePath(StateHeavySleetShowersAndThunder, n, w)},
		StateLightSnowShowersAndThunder:  {Name: StateLightSnowShowersAndThunder, Description: "light snow showers and thunder", ImageUrl: GetImagePath(StateLightSnowShowersAndThunder, n, w)},
		StateHeavySnowShowersAndThunder:  {Name: StateHeavySnowShowersAndThunder, Description: "heavy snow showers and thunder", ImageUrl: GetImagePath(StateHeavySnowShowersAndThunder, n, w)},
		StateLightRainAndThunder:         {Name: StateLightRainAndThunder, Description: "light rain and thunder", ImageUrl: GetImagePath(StateLightRainAndThunder, n, w)},
		StateLightSleetAndThunder:        {Name: StateLightSleetAndThunder, Description: "light sleet and thunder", ImageUrl: GetImagePath(StateLightSleetAndThunder, n, w)},
		StateHeavySleetAndThunder:        {Name: StateHeavySleetAndThunder, Description: "heavy sleet and thunder", ImageUrl: GetImagePath(StateHeavySleetAndThunder, n, w)},
		StateLightSnowAndThunder:         {Name: StateLightSnowAndThunder, Description: "light snow and thunder", ImageUrl: GetImagePath(StateLightSnowAndThunder, n, w)},
		StateHeavySnowAndThunder:         {Name: StateHeavySnowAndThunder, Description: "heavy snow and thunder", ImageUrl: GetImagePath(StateHeavySnowAndThunder, n, w)},
		StateLightRainShowers:            {Name: StateLightRainShowers, Description: "light rain showers", ImageUrl: GetImagePath(StateLightRainShowers, n, w)},
		StateHeavyRainShowers:            {Name: StateHeavyRainShowers, Description: "heavy rain showers", ImageUrl: GetImagePath(StateHeavyRainShowers, n, w)},
		StateLightSleetShowers:           {Name: StateLightSleetShowers, Description: "light sleet showers", ImageUrl: GetImagePath(StateLightSleetShowers, n, w)},
		StateHeavySleetShowers:           {Name: StateHeavySleetShowers, Description: "heavy sleet showers", ImageUrl: GetImagePath(StateHeavySleetShowers, n, w)},
		StateLightSnowShowers:            {Name: StateLightSnowShowers, Description: "light snow showers", ImageUrl: GetImagePath(StateLightSnowShowers, n, w)},
		StateHeavySnowShowers:            {Name: StateHeavySnowShowers, Description: "heavy snow showers", ImageUrl: GetImagePath(StateHeavySnowShowers, n, w)},
		StateLightRain:                   {Name: StateLightRain, Description: "light rain", ImageUrl: GetImagePath(StateLightRain, n, w)},
		StateLightSleet:                  {Name: StateLightSleet, Description: "light sleet", ImageUrl: GetImagePath(StateLightSleet, n, w)},
		StateHeavySleet:                  {Name: StateHeavySleet, Description: "heavy sleet", ImageUrl: GetImagePath(StateHeavySleet, n, w)},
		StateLightSnow:                   {Name: StateLightSnow, Description: "light snow", ImageUrl: GetImagePath(StateLightSnow, n, w)},
		StateHeavySnow:                   {Name: StateHeavySnow, Description: "heavy snow", ImageUrl: GetImagePath(StateHeavySnow, n, w)},
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
LOOP:
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
