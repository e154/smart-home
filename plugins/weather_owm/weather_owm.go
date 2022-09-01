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

package weather_owm

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/common/web"
	"net/url"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/system/entity_manager"
)

const (
	timeout = time.Second * 5
)

// WeatherOwm ...
type WeatherOwm struct {
	adaptors *adaptors.Adaptors
	crawler  web.Crawler
}

// NewWeatherOwm ...
func NewWeatherOwm(adaptors *adaptors.Adaptors, crawler web.Crawler) (weather *WeatherOwm) {
	weather = &WeatherOwm{
		adaptors: adaptors,
		crawler:  crawler,
	}

	return
}

// UpdateForecast ...
func (p *WeatherOwm) UpdateForecast(zone Zone, settings map[string]*m.Attribute) (forecast m.AttributeValue, err error) {
	forecast, err = p.GetForecast(zone, time.Now(), settings)
	return
}

// GetForecast ...
func (p *WeatherOwm) GetForecast(params Zone, now time.Time, settings map[string]*m.Attribute) (forecast m.AttributeValue, err error) {

	var zone Zone
	if zone, err = p.FetchData(params.Name, params.Lat, params.Lon, now, settings); err != nil {
		log.Errorf("%+v", err)
		return
	}

	// current weather
	if len(zone.Weatherdata.Daily) > 0 {
		var state entity_manager.ActorState
		if len(zone.Weatherdata.Current.Weather) > 0 {
			state = WeatherCondition(zone.Weatherdata.Current.Weather[0])
		}

		forecast = m.AttributeValue{
			weather.AttrWeatherDatetime:       time.Unix(zone.Weatherdata.Current.Dt, 0),
			weather.AttrWeatherMinTemperature: zone.Weatherdata.Daily[0].Temp.Min,
			weather.AttrWeatherMaxTemperature: zone.Weatherdata.Daily[0].Temp.Max,
			weather.AttrWeatherTemperature:    zone.Weatherdata.Current.Temp,
			weather.AttrWeatherPressure:       zone.Weatherdata.Current.Pressure,
			weather.AttrWeatherHumidity:       zone.Weatherdata.Current.Humidity,
			weather.AttrWeatherWindSpeed:      zone.Weatherdata.Current.WindSpeed,
			weather.AttrWeatherWindBearing:    zone.Weatherdata.Current.WindDeg,
			weather.AttrWeatherVisibility:     zone.Weatherdata.Current.Visibility,
			weather.AttrWeatherMain:           state.Name,
			weather.AttrWeatherDescription:    state.Description,
			weather.AttrWeatherIcon:           common.StringValue(state.ImageUrl),
		}
	}

	// week
	for i := 1; i < 6; i++ {
		if len(zone.Weatherdata.Daily) < i {
			return
		}
		var state entity_manager.ActorState
		if len(zone.Weatherdata.Daily[i].Weather) > 0 {
			state = WeatherCondition(zone.Weatherdata.Daily[i].Weather[0])
		}

		attrs := m.AttributeValue{
			weather.AttrWeatherDatetime:       time.Unix(zone.Weatherdata.Daily[i].Dt, 0),
			weather.AttrWeatherMinTemperature: zone.Weatherdata.Daily[i].Temp.Min,
			weather.AttrWeatherMaxTemperature: zone.Weatherdata.Daily[i].Temp.Max,
			weather.AttrWeatherPressure:       zone.Weatherdata.Daily[i].Pressure,
			weather.AttrWeatherHumidity:       zone.Weatherdata.Daily[i].Humidity,
			weather.AttrWeatherWindSpeed:      zone.Weatherdata.Daily[i].WindSpeed,
			weather.AttrWeatherWindBearing:    zone.Weatherdata.Daily[i].WindDeg,
			weather.AttrWeatherMain:           state.Name,
			weather.AttrWeatherDescription:    state.Description,
			weather.AttrWeatherIcon:           common.StringValue(state.ImageUrl),
		}
		for name, attr := range attrs {
			forecast[fmt.Sprintf("day%d_%s", i, name)] = attr
		}

	}

	forecast[weather.AttrWeatherAttribution] = Attribution

	return
}

// FetchData ...
func (p *WeatherOwm) FetchData(name string, lat, lon float64, now time.Time, settings map[string]*m.Attribute) (zone Zone, err error) {

	if lat == 0 || lon == 0 {
		err = errors.Wrap(apperr.ErrBadRequestParams, "zero positions")
		return
	}

	// get from storage
	if zone, err = p.fetchFromLocalStorage(name); err == nil {
		if zone.LoadedAt != nil && now.Sub(common.TimeValue(zone.LoadedAt)).Minutes() < 60 {
			return
		}
	} else {
		log.Error(err.Error())
	}

	// fetch from server
	var body []byte
	if body, err = p.fetchFromServer(lat, lon, settings); err != nil {
		err = errors.Wrap(errors.New("fetch from server failed"), err.Error())
		return
	}

	zone = Zone{
		Name: name,
		Lat:  lat,
		Lon:  lon,
	}
	zone.Weatherdata = &WeatherFor8Days{}
	zone.LoadedAt = common.Time(time.Now())
	if err = json.Unmarshal(body, zone.Weatherdata); err != nil {
		err = errors.Wrap(errors.New("failed unmarshal response body"), err.Error())
		return
	}

	// save to storage
	if err = p.saveToLocalStorage(zone); err != nil {
		err = errors.Wrap(errors.New("save to local storage failed"), err.Error())
	}

	return
}

func (p *WeatherOwm) fetchFromLocalStorage(name string) (zone Zone, err error) {

	log.Debugf("fetch from local storage")

	var variable m.Variable
	if variable, err = p.adaptors.Variable.GetByName(fmt.Sprintf("weather_owm.%s", name)); err != nil {
		return
	}

	zone = Zone{}
	err = json.Unmarshal([]byte(variable.Value), &zone)

	return
}

func (p *WeatherOwm) saveToLocalStorage(zone Zone) (err error) {

	log.Debugf("save to local storage '%s'", zone.Name)

	var b []byte
	if b, err = json.Marshal(zone); err != nil {
		return
	}

	err = p.adaptors.Variable.CreateOrUpdate(m.Variable{
		System:   true,
		Name:     fmt.Sprintf("weather_owm.%s", zone.Name),
		Value:    string(b),
		EntityId: common.NewEntityId(fmt.Sprintf("weather_owm.%s", zone.Name)),
	})

	return
}

func (p *WeatherOwm) fetchFromServer(lat, lon float64, settings map[string]*m.Attribute) (body []byte, err error) {

	uri, _ := url.Parse(DefaultApiUrl)
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", lat))
	params.Add("lon", fmt.Sprintf("%f", lon))
	params.Add("cnt", "48")

	if appid, ok := settings[AttrAppid]; ok {
		params.Add("appid", appid.String())
	}

	if units, ok := settings[AttrUnits]; ok {
		params.Add("units", units.String())
	}

	if lang, ok := settings[AttrLang]; ok {
		params.Add("lang", lang.String())
	}

	uri.RawQuery = params.Encode()

	log.Debugf("fetch from server %s\n", uri.String())

	_, body, err = p.crawler.Probe(web.Request{Method: "GET", Url: uri.String(), Timeout: timeout})

	return
}
