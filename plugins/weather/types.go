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

func BaseForecast() m.EntityAttributes {
	return NewForecast(6)
}

func NewForecast(days int) m.EntityAttributes {
	attributes := NewAttr()

	attributes[AttrWeatherAttribution] = &m.EntityAttribute{
		Name:  AttrWeatherAttribution,
		Type:  common.EntityAttributeString,
		Value: Attribution,
	}

	for i := 1; i < days; i++ {
		attributes[fmt.Sprintf("forecast_day%d", i)] = &m.EntityAttribute{
			Name:  fmt.Sprintf("forecast_day%d", i),
			Type:  common.EntityAttributeMap,
			Value: NewAttr(),
		}
	}

	return attributes
}

func NewAttr() m.EntityAttributes {
	return m.EntityAttributes{
		AttrWeatherDatetime: {
			Name: AttrWeatherDatetime,
			Type: common.EntityAttributeTime,
		},
		AttrWeatherTemperature: {
			Name: AttrWeatherTemperature,
			Type: common.EntityAttributeFloat,
		},
		AttrWeatherMinTemperature: {
			Name: AttrWeatherMinTemperature,
			Type: common.EntityAttributeFloat,
		},
		AttrWeatherMaxTemperature: {
			Name: AttrWeatherMaxTemperature,
			Type: common.EntityAttributeFloat,
		},
		AttrWeatherHumidity: {
			Name: AttrWeatherHumidity,
			Type: common.EntityAttributeFloat,
		},
		AttrWeatherPressure: {
			Name: AttrWeatherPressure,
			Type: common.EntityAttributeFloat,
		},
		AttrWeatherWindBearing: {
			Name: AttrWeatherWindBearing,
			Type: common.EntityAttributeFloat,
		},
		AttrWeatherWindSpeed: {
			Name: AttrWeatherWindSpeed,
			Type: common.EntityAttributeFloat,
		},
		AttrWeatherOzone: {
			Name: AttrWeatherOzone,
			Type: common.EntityAttributeFloat,
		},
		AttrWeatherVisibility: {
			Name: AttrWeatherVisibility,
			Type: common.EntityAttributeFloat,
		},
	}
}

type EventSubStateChanged struct {
	Type       common.EntityType  `json:"type"`
	EntityId   common.EntityId    `json:"entity_id"`
	State      string             `json:"state"`
	Attributes m.EntityAttributes `json:"attributes"`
}
