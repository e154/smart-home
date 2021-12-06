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

package weather_met

import (
	"encoding/xml"
	"time"
)

const (
	// Attribution ...
	Attribution = "Weather forecast from met.no, delivered by the Norwegian Meteorological Institute."
)

// MaxTemperature ...
type MaxTemperature struct {
	Id    string `xml:"id,attr"`
	Unit  string `xml:"unit,attr"`
	Value string `xml:"value,attr"`
}

// MinTemperature ...
type MinTemperature struct {
	Id    string `xml:"id,attr"`
	Unit  string `xml:"unit,attr"`
	Value string `xml:"value,attr"`
}

// Symbol ...
type Symbol struct {
	Id     string `xml:"id,attr"`
	Number string `xml:"number,attr"`
}

// Precipitation ...
type Precipitation struct {
	Unit  string `xml:"unit,attr"`
	Value string `xml:"value,attr"`
}

// DewpointTemperature ...
type DewpointTemperature struct {
	Id    string `xml:"id,attr"`
	Unit  string `xml:"unit,attr"`
	Value string `xml:"value,attr"`
}

// HighClouds ...
type HighClouds struct {
	Id      string `xml:"id,attr"`
	Percent string `xml:"percent,attr"`
}

// MediumClouds ...
type MediumClouds struct {
	Id      string `xml:"id,attr"`
	Percent string `xml:"percent,attr"`
}

// LowClouds ...
type LowClouds struct {
	Id      string `xml:"id,attr"`
	Percent string `xml:"percent,attr"`
}

// Fog ...
type Fog struct {
	Id      string `xml:"id,attr"`
	Percent string `xml:"percent,attr"`
}

// Cloudiness ...
type Cloudiness struct {
	Id      string `xml:"id,attr"`
	Percent string `xml:"percent,attr"`
}

// Pressure ...
type Pressure struct {
	Id    string `xml:"id,attr"`
	Unit  string `xml:"unit,attr"`
	Value string `xml:"value,attr"`
}

// Humidity ...
type Humidity struct {
	Unit  string `xml:"unit,attr"`
	Value string `xml:"value,attr"`
}

// WindSpeed ...
type WindSpeed struct {
	Id       string `xml:"id,attr"`
	Mps      string `xml:"mps,attr"`
	Beaufort string `xml:"beaufort,attr"`
	Name     string `xml:"name,attr"`
}

// WindDirection ...
type WindDirection struct {
	Id   string `xml:"id,attr"`
	Deg  string `xml:"deg,attr"`
	Name string `xml:"name,attr"`
}

// Temperature ...
type Temperature struct {
	Id    string `xml:"id,attr"`
	Unit  string `xml:"unit,attr"`
	Value string `xml:"value,attr"`
}

// Location ...
type Location struct {
	Altitude            string               `xml:"altitude,attr"`
	Latitude            string               `xml:"latitude,attr"`
	Longitude           string               `xml:"longitude,attr"`
	Temperature         *Temperature         `xml:"temperature,omitempty" json:"temperature,omitempty"`
	WindDirection       *WindDirection       `xml:"windDirection,omitempty" json:"windDirection,omitempty"`
	WindSpeed           *WindSpeed           `xml:"windSpeed,omitempty" json:"windSpeed,omitempty"`
	Humidity            *Humidity            `xml:"humidity,omitempty" json:"humidity,omitempty"`
	Pressure            *Pressure            `xml:"pressure,omitempty" json:"pressure,omitempty"`
	Cloudiness          *Cloudiness          `xml:"cloudiness,omitempty" json:"cloudiness,omitempty"`
	Fog                 *Fog                 `xml:"fog,omitempty" json:"fog,omitempty"`
	LowClouds           *LowClouds           `xml:"lowClouds,omitempty" json:"lowClouds,omitempty"`
	MediumClouds        *MediumClouds        `xml:"mediumClouds,omitempty" json:"mediumClouds,omitempty"`
	HighClouds          *HighClouds          `xml:"highClouds,omitempty" json:"highClouds,omitempty"`
	DewpointTemperature *DewpointTemperature `xml:"dewpointTemperature,omitempty" json:"dewpointTemperature,omitempty"`
	Precipitation       *Precipitation       `xml:"precipitation,omitempty" json:"precipitation,omitempty"`
	Symbol              *Symbol              `xml:"symbol,omitempty" json:"symbol,omitempty"`
	MinTemperature      *MinTemperature      `xml:"minTemperature,omitempty" json:"minTemperature,omitempty"`
	MaxTemperature      *MaxTemperature      `xml:"maxTemperature,omitempty" json:"maxTemperature,omitempty"`
}

// Product ...
type Product struct {
	Datatype string    `xml:"datatype,attr"`
	From     time.Time `xml:"from,attr"`
	To       time.Time `xml:"to,attr"`
	Location Location  `xml:"location"`
}

// Products ...
type Products []Product

// Len ...
func (p Products) Len() int {
	return len(p)
}

// Swap ...
func (p Products) Swap(a, b int) {
	p[a], p[b] = p[b], p[a]
}

// Less ...
func (p Products) Less(a, b int) bool {
	return p[a].From.Unix() < p[b].From.Unix()
}

// Weatherdata ...
type Weatherdata struct {
	XMLName  xml.Name `xml:"weatherdata"`
	Products Products `xml:"product>time"`
}

// Zone ...
type Zone struct {
	Name        string       `json:"name"`
	Lat         float64      `json:"lat"`
	Lon         float64      `json:"lon"`
	Weatherdata *Weatherdata `json:"weatherdata"`
	LoadedAt    *time.Time   `json:"loaded_at"`
}
