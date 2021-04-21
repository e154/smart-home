// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugin_manager"
	"go.uber.org/atomic"
	"math"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	Name          = "met"
	DefaultApiUrl = "https://api.met.no/weatherapi/locationforecast/2.0/classic"
)

var (
	log = common.MustGetLogger("plugins.met")
)

// MET Plugin
type pluginWeatherMet struct {
	adaptors      *adaptors.Adaptors
	entityManager *entity_manager.EntityManager
	isStarted     *atomic.Bool
	quit          chan struct{}
	pause         uint
	service       plugin_manager.IPluginManager
	eventBus      *event_bus.EventBus
	zones         *sync.Map
	lock          *sync.Mutex
	counter       int
}

func Register(manager *plugin_manager.PluginManager,
	entityManager *entity_manager.EntityManager,
	eventBus *event_bus.EventBus,
	adaptors *adaptors.Adaptors) {
	manager.Register(&pluginWeatherMet{
		adaptors:      adaptors,
		entityManager: entityManager,
		isStarted:     atomic.NewBool(false),
		quit:          make(chan struct{}),
		pause:         60,
		eventBus:      eventBus,
		zones:         &sync.Map{},
		lock:          &sync.Mutex{},
	})
}

func (p *pluginWeatherMet) Load(service plugin_manager.IPluginManager, plugins map[string]interface{}) (err error) {

	p.service = service

	if p.isStarted.Load() {
		return
	}

	p.eventBus.Subscribe(event_bus.TopicEntities, p.eventHandler)
	p.eventBus.Subscribe(weather.TopicPluginWeather, p.eventHandler)

	go func() {
		ticker := time.NewTicker(time.Minute * time.Duration(p.pause))
		p.quit = make(chan struct{})
		p.isStarted.Store(true)

		defer func() {
			ticker.Stop()
			close(p.quit)
		}()

		for {
			select {
			case <-p.quit:
				return
			case <-ticker.C:
				if err = p.updateForecastForAll(); err != nil {
					log.Error(err.Error())
				}
			}
		}
	}()

	if err = p.updateForecastForAll(); err != nil {
		log.Error(err.Error())
	}

	return nil
}

func (p *pluginWeatherMet) Unload() (err error) {
	if !p.isStarted.Load() {
		return
	}
	p.isStarted.Store(false)
	p.quit <- struct{}{}
	p.eventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)
	p.eventBus.Unsubscribe(weather.TopicPluginWeather, p.eventHandler)
	return nil
}

func (p pluginWeatherMet) Name() string {
	return Name
}

func (p *pluginWeatherMet) eventHandler(msg interface{}) {
	switch v := msg.(type) {
	case event_bus.EventAddedNewEntity:
		if v.Type != "weather" {
			return
		}

		p.addZone()

	case event_bus.EventSubStateChanged:
		if v.Type != "weather" || v.State != weather.SubStatePositionUpdate {
			return
		}

		p.updateZonesList(v.EntityId, v.Attributes)

	case event_bus.EventRemoveEntity:
		if v.Type != "weather" {
			return
		}
		p.zones.Delete(v.EntityId.Name())

	}
}

func (p *pluginWeatherMet) addZone() {

}

func (p *pluginWeatherMet) updateZonesList(entityId common.EntityId, attrs m.EntityAttributes) {

	zone := weather.Zone{
		Name:      entityId.Name(),
		Lat:       attrs[zone.AttrLat].Float64(),
		Lon:       attrs[zone.AttrLon].Float64(),
		Elevation: attrs[zone.AttrElevation].Float64(),
	}

	var update bool
	if _, ok := p.zones.Load(entityId); !ok {
		update = true
	}
	p.zones.Store(entityId, zone)

	if !update {
		return
	}
	p.updateForecastForAll()
}

func (p *pluginWeatherMet) send(to common.EntityId, payload interface{}) {
	p.entityManager.Send(entity_manager.Message{
		From:    Name,
		To:      to,
		Payload: payload,
	})
}

func (p *pluginWeatherMet) updateForecastForAll() (err error) {

	p.lock.Lock()
	defer p.lock.Unlock()

	p.zones.Range(func(key, value interface{}) bool {
		zone, ok := value.(weather.Zone)
		if !ok {
			return true
		}

		if err = p.updateForecast(zone); err != nil {
			log.Error(err.Error())
		}

		return true
	})

	return nil
}

func (p *pluginWeatherMet) updateForecast(zone weather.Zone) (err error) {

	var forecast m.EntityAttributeValue
	forecast, err = p.GetForecast(zone)
	if err != nil {
		log.Error(err.Error())
		return
	}

	attr := weather.BaseForecast()
	attr.Deserialize(forecast)

	p.send(common.EntityId(fmt.Sprintf("weather.%s", zone.Name)), entity_manager.MessageRequestState{
		Name:       "forecast",
		Attributes: attr,
	})

	return nil
}

func (p *pluginWeatherMet) GetForecast(params weather.Zone) (forecast m.EntityAttributeValue, err error) {

	var zone Zone
	if zone, err = p.fetchData(params.Name, params.Lat, params.Lon, params.Elevation); err != nil {
		log.Error(err.Error())
		return
	}

	if forecast, err = p.getCurrentWeather(zone); err != nil {
		return
	}

	now := time.Now()
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

func (p *pluginWeatherMet) fetchData(name string, lat, lon, elevation float64) (zone Zone, err error) {

	if lat == 0 || lon == 0 {
		err = fmt.Errorf("zero positions")
		return
	}

	// get from storage
	if zone, err = p.fetchFromLocalStorage(name); err == nil {
		if zone.LoadetAt != nil && time.Now().Sub(common.TimeValue(zone.LoadetAt)).Minutes() < 60 {
			return
		}
	}

	// fetch from server
	var body []byte
	if body, err = p.fetchFromServer(lat, lon, elevation); err != nil {
		return
	}

	zone.Zone = weather.Zone{
		Name:      name,
		Lat:       lat,
		Lon:       lon,
		Elevation: elevation,
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

func (p *pluginWeatherMet) fetchFromLocalStorage(name string) (zone Zone, err error) {

	log.Debugf("fetch from local storage")

	var variable m.Variable
	if variable, err = p.adaptors.Variable.GetByName(fmt.Sprintf("weather.%s", name)); err != nil {
		return
	}

	zone = Zone{}
	err = json.Unmarshal([]byte(variable.Value), &zone)

	return
}

func (p *pluginWeatherMet) saveToLocalStorage(zone Zone) (err error) {

	log.Debugf("save to local storage %s", zone.Name)

	var b []byte
	if b, err = json.Marshal(zone); err != nil {
		return
	}

	err = p.adaptors.Variable.Update(m.Variable{
		Name:  fmt.Sprintf("weather.%s", zone.Name),
		Value: string(b),
	})

	return
}

func (p *pluginWeatherMet) fetchFromServer(lat, lon, elevation float64) (body []byte, err error) {

	uri, _ := url.Parse(DefaultApiUrl)
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", lat))
	params.Add("lon", fmt.Sprintf("%f", lon))
	//if elevation != 0 {
	//	params.Add("msl", fmt.Sprintf("%f", elevation))
	//}
	uri.RawQuery = params.Encode()

	log.Debugf("fetch from server %s", uri.String())

	body, err = web.Crawler(uri.String())

	return
}

func (p *pluginWeatherMet) getCurrentWeather(zone Zone) (w6 m.EntityAttributeValue, err error) {

	now := time.Now()

	if w6, err = p.getWeather(zone, now, 6); err != nil {
		return
	}

	var w24 m.EntityAttributeValue
	if w24, err = p.getWeather(zone, now, 24); err != nil {
		return
	}

	w6[weather.AttrWeatherMaxTemperature] = w24[weather.AttrWeatherMaxTemperature]
	w6[weather.AttrWeatherMinTemperature] = w24[weather.AttrWeatherMinTemperature]

	return
}

func (p *pluginWeatherMet) getWeather(zone Zone, t time.Time, maxHour float64) (w m.EntityAttributeValue, err error) {

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

	w = m.EntityAttributeValue{
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

func (p *pluginWeatherMet) getTemperature(orderedEntries []Product) (value float64) {

	for _, entry := range orderedEntries {
		if entry.Location.Temperature != nil {
			value, _ = strconv.ParseFloat(entry.Location.Temperature.Value, 32)
			value = math.Round(value*100) / 100
			return
		}
	}
	return
}

func (p *pluginWeatherMet) getMinTemperature(orderedEntries []Product) (value float64) {

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

func (p *pluginWeatherMet) getMaxTemperature(orderedEntries []Product) (value float64) {

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

func (p *pluginWeatherMet) getCondition(orderedEntries []Product) (value string) {

	for _, entry := range orderedEntries {
		if entry.Location.Symbol != nil {
			num, _ := strconv.Atoi(entry.Location.Symbol.Number)
			if num < len(conditions) && conditions[num] > value {
				value = conditions[num]
			}
		}
	}
	return
}

func (p *pluginWeatherMet) getPressure(orderedEntries []Product) (value float64) {

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

func (p *pluginWeatherMet) getHumidity(orderedEntries []Product) (value float64) {

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

func (p *pluginWeatherMet) getDewpointTemperature(orderedEntries []Product) (value float64) {

	for _, entry := range orderedEntries {
		if entry.Location.DewpointTemperature != nil {
			value, _ = strconv.ParseFloat(entry.Location.DewpointTemperature.Value, 32)
			value = math.Round(value*100) / 100
			return
		}
	}
	return
}

func (p *pluginWeatherMet) getWindSpeed(orderedEntries []Product) (value float64) {

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

func (p *pluginWeatherMet) getWindDirection(orderedEntries []Product) (value float64) {

	for _, entry := range orderedEntries {
		if entry.Location.WindDirection != nil {
			value, _ = strconv.ParseFloat(entry.Location.WindDirection.Deg, 32)
			value = math.Round(value*100) / 100
			return
		}
	}
	return
}

func (p *pluginWeatherMet) Type() plugin_manager.PlugableType {
	return plugin_manager.PlugableBuiltIn
}

func (p *pluginWeatherMet) Depends() []string {
	return []string{weather.Name}
}
