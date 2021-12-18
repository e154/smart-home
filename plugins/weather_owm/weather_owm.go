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
	"encoding/xml"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/web"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
)

// WeatherOwm ...
type WeatherOwm struct {
	adaptors *adaptors.Adaptors
	eventBus event_bus.EventBus
	settings map[string]*m.Attribute
	lock     *sync.Mutex
	zones    *sync.Map
}

// NewWeatherOwm ...
func NewWeatherOwm(eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors,
	settings map[string]*m.Attribute) (weather *WeatherOwm) {
	weather = &WeatherOwm{
		eventBus: eventBus,
		adaptors: adaptors,
		settings: settings,
		lock:     &sync.Mutex{},
		zones:    &sync.Map{},
	}

	return
}

// AddWeather ...
func (p *WeatherOwm) AddWeather(entityId common.EntityId, settings m.Attributes) {
	p.UpdateWeatherList(entityId, settings)
}

// UpdateWeatherList ...
func (p *WeatherOwm) UpdateWeatherList(entityId common.EntityId, settings m.Attributes) {

	zone := Zone{
		Name: entityId.Name(),
		Lat:  settings[weather.AttrLat].Float64(),
		Lon:  settings[weather.AttrLon].Float64(),
	}

	var update bool
	if _, ok := p.zones.Load(entityId.Name()); !ok {
		update = true
	}
	p.zones.Store(entityId, zone)

	if !update {
		return
	}
	p.UpdateForecastForAll()
}

// RemoveWeather ...
func (p *WeatherOwm) RemoveWeather(entityId common.EntityId) {
	p.zones.Delete(entityId.Name())
	log.Infof("unload weather_owm.%s", entityId.Name())

	p.eventBus.Publish(event_bus.TopicEntities, event_bus.EventRemoveActor{
		Type:     "weather_owm",
		EntityId: common.EntityId(fmt.Sprintf("weather_owm.%s", entityId.Name())),
	})
}

// UpdateForecastForAll ...
func (p *WeatherOwm) UpdateForecastForAll() (err error) {

	p.lock.Lock()
	defer p.lock.Unlock()

	p.zones.Range(func(key, value interface{}) bool {
		zone, ok := value.(Zone)
		if !ok {
			return true
		}

		if err = p.UpdateForecast(zone); err != nil {
			log.Error(err.Error())
		}

		return true
	})

	return nil
}

// UpdateForecast ...
func (p *WeatherOwm) UpdateForecast(zone Zone) (err error) {

	var forecast m.AttributeValue
	if forecast, err = p.GetForecast(zone, time.Now()); err != nil {
		return
	}

	attr := weather.BaseForecast()
	attr.Deserialize(forecast)

	p.eventBus.Publish(event_bus.TopicEntities, event_bus.EventRequestState{
		From:       common.EntityId(fmt.Sprintf("weather_owm.%s", zone.Name)),
		To:         common.EntityId(fmt.Sprintf("weather.%s", zone.Name)),
		Attributes: attr,
	})

	return nil
}

// GetForecast ...
func (p *WeatherOwm) GetForecast(params Zone, now time.Time) (forecast m.AttributeValue, err error) {

	var zone Zone
	if zone, err = p.FetchData(params.Name, params.Lat, params.Lon, now); err != nil {
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
		forecast[fmt.Sprintf("forecast_day%d", i)] = m.AttributeValue{
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
	}

	forecast[weather.AttrWeatherAttribution] = Attribution

	return
}

// FetchData ...
func (p *WeatherOwm) FetchData(name string, lat, lon float64, now time.Time) (zone Zone, err error) {

	if lat == 0 || lon == 0 {
		err = fmt.Errorf("zero positions")
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
	if body, err = p.fetchFromServer(lat, lon); err != nil {
		return
	}

	zone = Zone{
		Name: name,
		Lat:  lat,
		Lon:  lon,
	}
	zone.Weatherdata = &WeatherFor8Days{}
	zone.LoadedAt = common.Time(time.Now())
	if err = xml.Unmarshal(body, zone.Weatherdata); err != nil {
		return
	}

	// save to storage
	err = p.saveToLocalStorage(zone)

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
		Name:     fmt.Sprintf("weather_owm.%s", zone.Name),
		Value:    string(b),
		EntityId: common.NewEntityId(fmt.Sprintf("weather.%s", zone.Name)),
	})

	return
}

func (p *WeatherOwm) fetchFromServer(lat, lon float64) (body []byte, err error) {

	uri, _ := url.Parse(DefaultApiUrl)
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", lat))
	params.Add("lon", fmt.Sprintf("%f", lon))
	params.Add("appid", p.settings[AttrAppid].String())
	params.Add("units", p.settings[AttrUnits].String())
	params.Add("lang", p.settings[AttrLang].String())
	params.Add("cnt", "48")

	uri.RawQuery = params.Encode()

	log.Debugf("fetch from server %s\n", uri.String())

	body, err = web.Crawler(web.Request{Method: "GET", Url: uri.String()})

	return
}
