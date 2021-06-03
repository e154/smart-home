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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/web"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/plugins/zone"
	"github.com/e154/smart-home/system/event_bus"
	"math"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"
)

type WeatherMet struct {
	adaptors *adaptors.Adaptors
	eventBus event_bus.EventBus
	lock     *sync.Mutex
	zones    *sync.Map
}

func NewWeatherMet(eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors) (weather *WeatherMet) {
	weather = &WeatherMet{
		eventBus: eventBus,
		adaptors: adaptors,
		lock:     &sync.Mutex{},
		zones:    &sync.Map{},
	}

	return
}

func (p *WeatherMet) AddWeather(entityId common.EntityId, settings m.Attributes) {
	p.UpdateWeatherList(entityId, settings)
}

func (p *WeatherMet) UpdateWeatherList(entityId common.EntityId, settings m.Attributes) {

	zone := Zone{
		Name: entityId.Name(),
		Lat:  settings[zone.AttrLat].Float64(),
		Lon:  settings[zone.AttrLon].Float64(),
	}

	var update bool
	if _, ok := p.zones.Load(entityId); !ok {
		update = true
	}
	p.zones.Store(entityId, zone)

	if !update {
		return
	}
	p.UpdateForecastForAll()
}

func (p *WeatherMet) RemoveWeather(entityId common.EntityId) {
	p.zones.Delete(entityId.Name())
	log.Infof("unload weather_met.%s", entityId.Name())

	p.eventBus.Publish(event_bus.TopicEntities, event_bus.EventRemoveActor{
		Type:     "weather_met",
		EntityId: common.EntityId(fmt.Sprintf("weather_met.%s", entityId.Name())),
	})
}

func (p *WeatherMet) UpdateForecastForAll() (err error) {

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

func (p *WeatherMet) UpdateForecast(zone Zone) (err error) {

	var forecast m.AttributeValue
	if forecast, err = p.GetForecast(zone, time.Now()); err != nil {
		return
	}

	attr := weather.BaseForecast()
	attr.Deserialize(forecast)

	p.eventBus.Publish(event_bus.TopicEntities, event_bus.EventRequestState{
		From:       common.EntityId(fmt.Sprintf("weather_met.%s", zone.Name)),
		To:         common.EntityId(fmt.Sprintf("weather.%s", zone.Name)),
		Attributes: attr,
	})

	return nil
}

func (p *WeatherMet) GetForecast(params Zone, now time.Time) (forecast m.AttributeValue, err error) {

	var zone Zone
	if zone, err = p.FetchData(params.Name, params.Lat, params.Lon, now); err != nil {
		return
	}

	if forecast, err = p.getCurrentWeather(zone, now); err != nil {
		return
	}

	rounded := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for i := 1; i < 6; i++ {
		t := rounded.AddDate(0, 0, i)
		weather, err := p.getWeather(zone, t, 24)
		if err != nil {
			log.Error(err.Error())
		}
		forecast[fmt.Sprintf("forecast_day%d", i)] = weather
	}

	return
}

func (p *WeatherMet) FetchData(name string, lat, lon float64, now time.Time) (zone Zone, err error) {

	if lat == 0 || lon == 0 {
		err = fmt.Errorf("zero positions")
		return
	}

	// get from storage
	if zone, err = p.fetchFromLocalStorage(name); err == nil {
		if zone.LoadetAt != nil && now.Sub(common.TimeValue(zone.LoadetAt)).Minutes() < 60 {
			return
		}
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
	zone.Weatherdata = &Weatherdata{}
	zone.LoadetAt = common.Time(time.Now())
	if err = xml.Unmarshal(body, zone.Weatherdata); err != nil {
		return
	}

	// save to storage
	err = p.saveToLocalStorage(zone)

	return
}

func (p *WeatherMet) fetchFromLocalStorage(name string) (zone Zone, err error) {

	log.Debugf("fetch from local storage")

	var variable m.Variable
	if variable, err = p.adaptors.Variable.GetByName(fmt.Sprintf("weather_met.%s", name)); err != nil {
		return
	}

	zone = Zone{}
	err = json.Unmarshal([]byte(variable.Value), &zone)

	return
}

func (p *WeatherMet) saveToLocalStorage(zone Zone) (err error) {

	log.Debugf("save to local storage '%s'", zone.Name)

	var b []byte
	if b, err = json.Marshal(zone); err != nil {
		return
	}

	err = p.adaptors.Variable.CreateOrUpdate(m.Variable{
		Name:     fmt.Sprintf("weather_met.%s", zone.Name),
		Value:    string(b),
		EntityId: common.NewEntityId(fmt.Sprintf("weather.%s", zone.Name)),
	})

	return
}

func (p *WeatherMet) fetchFromServer(lat, lon float64) (body []byte, err error) {

	uri, _ := url.Parse(DefaultApiUrl)
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", lat))
	params.Add("lon", fmt.Sprintf("%f", lon))

	uri.RawQuery = params.Encode()

	log.Debugf("fetch from server %s", uri.String())

	body, err = web.Crawler(uri.String())

	return
}

func (p *WeatherMet) getCurrentWeather(zone Zone, now time.Time) (w6 m.AttributeValue, err error) {

	if w6, err = p.getWeather(zone, now, 6); err != nil {
		return
	}

	var w24 m.AttributeValue
	if w24, err = p.getWeather(zone, now, 24); err != nil {
		return
	}

	w6[weather.AttrWeatherMaxTemperature] = w24[weather.AttrWeatherMaxTemperature]
	w6[weather.AttrWeatherMinTemperature] = w24[weather.AttrWeatherMinTemperature]

	return
}

func (p *WeatherMet) getWeather(zone Zone, t time.Time, maxHour float64) (w m.AttributeValue, err error) {

	var orderedEntries Products
	for _, product := range zone.Weatherdata.Products {
		if t.Sub(product.To).Minutes() > 0 {
			continue
		}

		averageDist := product.To.Sub(t).Hours() + product.From.Sub(t).Hours()

		if averageDist > maxHour {
			continue
		}

		orderedEntries = append(orderedEntries, product)
	}

	sort.Sort(orderedEntries)

	w = m.AttributeValue{
		weather.AttrWeatherDatetime:       t,
		weather.AttrWeatherMinTemperature: p.getMinTemperature(orderedEntries),
		weather.AttrWeatherMaxTemperature: p.getMaxTemperature(orderedEntries),
		weather.AttrWeatherCondition:      p.getCondition(orderedEntries),
		weather.AttrWeatherPressure:       p.getPressure(orderedEntries),
		weather.AttrWeatherHumidity:       p.getHumidity(orderedEntries),
		weather.AttrWeatherWindSpeed:      p.getWindSpeed(orderedEntries),
		weather.AttrWeatherWindBearing:    p.getWindDirection(orderedEntries),
	}

	if maxHour <= 6 {
		w[weather.AttrWeatherTemperature] = p.getTemperature(orderedEntries)
	}

	return
}

func (p *WeatherMet) getTemperature(orderedEntries []Product) (value float64) {

	for _, entry := range orderedEntries {
		if entry.Location.Temperature != nil {
			value, _ = strconv.ParseFloat(entry.Location.Temperature.Value, 32)
			value = math.Round(value*100) / 100
			return
		}
	}
	return
}

func (p *WeatherMet) getMinTemperature(orderedEntries []Product) (value float64) {

	value = 99
	for _, entry := range orderedEntries {
		if entry.Location.MinTemperature != nil {
			_value, _ := strconv.ParseFloat(entry.Location.MinTemperature.Value, 32)
			if value == 99 {
				value = _value
			}
			_value = math.Round(_value*100) / 100
			if _value < value {
				value = _value
			}
		}
	}

	if value == 99 {
		value = 0
	}

	return
}

func (p *WeatherMet) getMaxTemperature(orderedEntries []Product) (value float64) {

	for _, entry := range orderedEntries {
		if entry.Location.MaxTemperature != nil {
			_value, _ := strconv.ParseFloat(entry.Location.MaxTemperature.Value, 32)
			_value = math.Round(_value*100) / 100
			if _value > value {
				value = _value
			}
		}
	}
	return
}

func (p *WeatherMet) getCondition(orderedEntries []Product) (value string) {

	for _, entry := range orderedEntries {
		if entry.Location.Symbol != nil {
			num, _ := strconv.Atoi(entry.Location.Symbol.Number)
			value = weather.GetStateByIndex(num)
		}
	}
	return
}

func (p *WeatherMet) getPressure(orderedEntries []Product) (value float64) {

	for _, entry := range orderedEntries {
		if entry.Location.Pressure != nil {
			_value, _ := strconv.ParseFloat(entry.Location.Pressure.Value, 32)
			_value = math.Round(_value*100) / 100
			if _value > value {
				value = _value
			}
		}
	}
	return
}

func (p *WeatherMet) getHumidity(orderedEntries []Product) (value float64) {

	for _, entry := range orderedEntries {
		if entry.Location.Humidity != nil {
			_value, _ := strconv.ParseFloat(entry.Location.Humidity.Value, 32)
			_value = math.Round(_value*100) / 100
			if _value > value {
				value = _value
			}
		}
	}
	return
}

func (p *WeatherMet) getDewpointTemperature(orderedEntries []Product) (value float64) {

	for _, entry := range orderedEntries {
		if entry.Location.DewpointTemperature != nil {
			value, _ = strconv.ParseFloat(entry.Location.DewpointTemperature.Value, 32)
			value = math.Round(value*100) / 100
			return
		}
	}
	return
}

func (p *WeatherMet) getWindSpeed(orderedEntries []Product) (value float64) {

	for _, entry := range orderedEntries {
		if entry.Location.WindSpeed != nil {
			mps, _ := strconv.ParseFloat(entry.Location.WindSpeed.Mps, 32)
			value = mps * 3.6
			value = math.Round(value*100) / 100
			return
		}
	}
	return
}

func (p *WeatherMet) getWindDirection(orderedEntries []Product) (value float64) {

	for _, entry := range orderedEntries {
		if entry.Location.WindDirection != nil {
			value, _ = strconv.ParseFloat(entry.Location.WindDirection.Deg, 32)
			value = math.Round(value*100) / 100
			return
		}
	}
	return
}
