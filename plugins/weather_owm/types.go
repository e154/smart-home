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
	"time"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/system/entity_manager"
)

const (
	// EntityWeatherOwm ...
	EntityWeatherOwm = common.EntityType("weather_owm")
	// DefaultApiUrl ...
	DefaultApiUrl = "https://api.openweathermap.org/data/2.5/onecall"
	// Attribution ...
	Attribution = "Weather forecast from openweathermap api"
)

// GeoPos ...
type GeoPos struct {
	Lat float64 `json:"lat"` // latitude
	Lon float64 `json:"lon"` // longitude
}

// City ...
type City struct {
	Id         int64  `json:"id"`         // City ID
	Name       string `json:"name"`       // City name
	Coord      GeoPos `json:"coord"`      // City geo location
	Country    string `json:"country"`    // Country code (GB, JP etc.)
	Population int64  `json:"population"` //
	Timezone   int64  `json:"timezone"`   // Shift in seconds from UTC
	Sunrise    int64  `json:"sunrise"`    // Sunrise time
	Sunset     int64  `json:"sunset"`     // Sunset time
}

// ProductMain ...
type ProductMain struct {
	Temp      float64 `json:"temp"`       // Temperature. Unit Default: Kelvin, Metric: Celsius, Imperial: Fahrenheit.
	FeelsLike float64 `json:"feels_like"` // This temperature parameter accounts for the human perception of weather. Unit Default: Kelvin, Metric: Celsius, Imperial: Fahrenheit.
	TempMin   float64 `json:"temp_min"`   // Minimum temperature at the moment of calculation. This is minimal forecasted temperature (within large megalopolises and urban areas), use this parameter optionally. Unit Default: Kelvin, Metric: Celsius, Imperial: Fahrenheit.
	TempMax   float64 `json:"temp_max"`   // Maximum temperature at the moment of calculation. This is maximal forecasted temperature (within large megalopolises and urban areas), use this parameter optionally. Unit Default: Kelvin, Metric: Celsius, Imperial: Fahrenheit.
	Pressure  float64 `json:"pressure"`   // Atmospheric pressure on the sea level by default, hPa
	SeaLevel  float64 `json:"sea_level"`  // Atmospheric pressure on the sea level, hPa
	GrndLevel float64 `json:"grnd_level"` // Atmospheric pressure on the ground level, hPa
	Humidity  float64 `json:"humidity"`   // Humidity, %
	TempKf    float64 `json:"temp_kf"`    // Internal parameter
}

// ProductWeather ...
type ProductWeather struct {
	Id          int64  `json:"id"`          // Weather condition id
	Main        string `json:"main"`        // Group of weather parameters (Rain, Snow, Extreme etc.)
	Description string `json:"description"` // Weather condition within the group. You can get the output in your language.
	Icon        string `json:"icon"`        // Weather icon id
}

// ProductClouds ...
type ProductClouds struct {
	All int64 `json:"all"` // Cloudiness, %
}

// ProductWind ...
type ProductWind struct {
	Speed float64 `json:"speed"` // Wind speed. Unit Default: meter/sec, Metric: meter/sec, Imperial: miles/hour.
	Deg   float64 `json:"deg"`   // Wind direction, degrees (meteorological)
	Gust  float64 `json:"gust"`  // Wind gust. Unit Default: meter/sec, Metric: meter/sec, Imperial: miles/hour
}

// ProductSnow ...
type ProductSnow struct {
	Last3Hours float64 `json:"3h"` // Snow volume for last 3 hours
}

// ProductRain ...
type ProductRain struct {
	Last3Hours float64 `json:"3h"` // Rain volume for last 3 hours, mm
}

// ProductSys ...
type ProductSys struct {
	Pod string `json:"pod"` // Part of the day (n - night, d - day)
}

// Product ...
type Product struct {
	Dt         int64          `json:"dt"`         // Time of data forecasted, unix, UTC
	Visibility int64          `json:"visibility"` // Average visibility, metres
	Pop        float64        `json:"pop"`        //  Probability of precipitation
	Main       ProductMain    `json:"main"`
	Weather    ProductWeather `json:"weather"`
	Clouds     *ProductClouds `json:"clouds,omitempty"`
	Wind       *ProductWind   `json:"wind,omitempty"`
	Rain       *ProductRain   `json:"rain,omitempty"`
	Snow       *ProductSnow   `json:"snow,omitempty"`
	Sys        *ProductSys    `json:"sys,omitempty"`
	DtTxt      time.Time      `json:"dt_txt"`
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
	return p[a].Dt < p[b].Dt
}

// WeatherFor5Days ...
type WeatherFor5Days struct {
	Cod      string   `json:"cod"`     //  Internal parameter
	Message  int      `json:"message"` //  Internal parameter
	Cnt      int      `json:"cnt"`     //  A number of timestamps returned in the API response
	Products Products `json:"list"`
	City     City     `json:"city"`
}

// WeatherCurrentRain ...
type WeatherCurrentRain struct {
	LastHour float64 `json:"1h"` // Rain volume for last hour, mm
}

// WeatherCurrentSnow ...
type WeatherCurrentSnow struct {
	LastHour float64 `json:"1h"` // Snow volume for last hour, mm
}

// WeatherCurrent ...
type WeatherCurrent struct {
	Dt         int64               `json:"dt"`                   // Current time, Unix, UTC
	Sunrise    int64               `json:"sunrise"`              // Sunrise time, Unix, UTC
	Sunset     int64               `json:"sunset"`               // Sunset time, Unix, UTC
	Temp       float64             `json:"temp"`                 // Temperature. Units - default: kelvin, metric: Celsius, imperial: Fahrenheit.
	FeelsLike  float64             `json:"feels_like"`           // Temperature. This temperature parameter accounts for the human perception of weather
	Pressure   float64             `json:"pressure"`             // Atmospheric pressure on the sea level, hPa
	Humidity   float64             `json:"humidity"`             // Humidity, %
	DewPoint   float64             `json:"dew_point"`            // Atmospheric temperature (varying according to pressure and humidity) below which water droplets begin to condense and dew can form.
	Uvi        float64             `json:"uvi"`                  // Current UV index
	Clouds     float64             `json:"clouds"`               // Cloudiness, %
	Visibility int64               `json:"visibility"`           // Average visibility, metres
	WindSpeed  float64             `json:"wind_speed,omitempty"` // Wind speed. Wind speed.
	WindDeg    float64             `json:"wind_deg,omitempty"`   // Wind direction, degrees (meteorological)
	WindGust   *float64            `json:"wind_gust,omitempty"`  // (where available) Wind gust.
	Rain       *WeatherCurrentRain `json:"rain"`                 //
	Snow       *WeatherCurrentSnow `json:"snow"`                 //
	Weather    []ProductWeather    `json:"weather"`              //
}

// WeatherHourly ...
type WeatherHourly struct {
	Dt         int64               `json:"dt"`                   // Current time, Unix, UTC
	Temp       float64             `json:"temp"`                 // Temperature. Units - default: kelvin, metric: Celsius, imperial: Fahrenheit.
	FeelsLike  float64             `json:"feels_like"`           // Temperature. This temperature parameter accounts for the human perception of weather
	Pressure   float64             `json:"pressure"`             // Atmospheric pressure on the sea level, hPa
	Humidity   float64             `json:"humidity"`             // Humidity, %
	DewPoint   float64             `json:"dew_point"`            // Atmospheric temperature (varying according to pressure and humidity) below which water droplets begin to condense and dew can form.
	Uvi        float64             `json:"uvi"`                  // Current UV index
	Clouds     float64             `json:"clouds"`               // Cloudiness, %
	Visibility int64               `json:"visibility"`           // Average visibility, metres
	WindSpeed  float64             `json:"wind_speed,omitempty"` // Wind speed. Wind speed.
	WindDeg    float64             `json:"wind_deg,omitempty"`   // Wind direction, degrees (meteorological)
	WindGust   *float64            `json:"wind_gust,omitempty"`  // (where available) Wind gust.
	Rain       *WeatherCurrentRain `json:"rain"`                 //
	Snow       *WeatherCurrentSnow `json:"snow"`                 //
	Weather    []ProductWeather    `json:"weather"`              //
	Pop        float64             `json:"pop"`                  // Probability of precipitation
}

// WeatherDailyTemp ...
type WeatherDailyTemp struct {
	Day   float64 `json:"day"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"` // Evening temperature.
	Morn  float64 `json:"morn"`
}

// WeatherDailyFeelsLike ...
type WeatherDailyFeelsLike struct {
	Day   float64 `json:"day"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"` // Evening temperature.
	Morn  float64 `json:"morn"`
}

// Alert ...
type Alert struct {
	SenderName  string   `json:"sender_name"` // Name of the alert source
	Event       string   `json:"event"`       // Alert event name
	Start       int64    `json:"start"`       // Date and time of the start of the alert, Unix, UTC
	End         int64    `json:"end"`         // Date and time of the end of the alert, Unix, UTC
	Description string   `json:"description"` // Description of the alert
	Tags        []string `json:"tags"`        // Type of severe weather
}

// WeatherDaily ...
type WeatherDaily struct {
	Dt         int64                 `json:"dt"`                   // Current time, Unix, UTC
	Sunrise    int64                 `json:"sunrise"`              // Sunrise time, Unix, UTC
	Sunset     int64                 `json:"sunset"`               // Sunset time, Unix, UTC
	Moonrise   int64                 `json:"moonrise"`             // The time of when the moon rises for this day, Unix, UTC
	Moonset    int64                 `json:"moonset"`              // The time of when the moon sets for this day, Unix, UTC
	MoonPhase  float64               `json:"moon_phase"`           // Moon phase. 0 and 1 are 'new moon', 0.25 is 'first quarter moon', 0.5 is 'full moon' and 0.75 is 'last quarter moon'. The periods in between are called 'waxing crescent', 'waxing gibous', 'waning gibous', and 'waning crescent', respectively.
	Temp       WeatherDailyTemp      `json:"temp"`                 // Temperature. Units - default: kelvin, metric: Celsius, imperial: Fahrenheit.
	FeelsLike  WeatherDailyFeelsLike `json:"feels_like"`           // Temperature. This temperature parameter accounts for the human perception of weather
	Pressure   int64                 `json:"pressure"`             // Atmospheric pressure on the sea level, hPa
	Humidity   float64               `json:"humidity"`             // Humidity, %
	DewPoint   float64               `json:"dew_point"`            // Atmospheric temperature (varying according to pressure and humidity) below which water droplets begin to condense and dew can form.
	Uvi        float64               `json:"uvi"`                  // Current UV index
	Clouds     float64               `json:"clouds"`               // Cloudiness, %
	Visibility float64               `json:"visibility"`           // Average visibility, metres
	WindSpeed  float64               `json:"wind_speed,omitempty"` // Wind speed. Wind speed.
	WindDeg    float64               `json:"wind_deg,omitempty"`   // Wind direction, degrees (meteorological)
	WindGust   *float64              `json:"wind_gust,omitempty"`  // (where available) Wind gust.
	Rain       float64               `json:"rain"`                 //
	Snow       float64               `json:"snow"`                 //
	Weather    []ProductWeather      `json:"weather"`              //
}

// WeatherFor8Days ...
type WeatherFor8Days struct {
	GeoPos
	Timezone       string          `json:"timezone"`        // Timezone name for the requested location
	TimezoneOffset int64           `json:"timezone_offset"` // Shift in seconds from UTC
	Current        WeatherCurrent  `json:"current"`
	Hourly         []WeatherHourly `json:"hourly"`
	Daily          []WeatherDaily  `json:"daily"`
	Alerts         []Alert         `json:"alerts"`
}

// Zone ...
type Zone struct {
	Name        string           `json:"name"`
	Lat         float64          `json:"lat"`
	Lon         float64          `json:"lon"`
	Weatherdata *WeatherFor8Days `json:"weatherdata"`
	LoadedAt    *time.Time       `json:"loaded_at"`
}

const (
	// AttrAppid ...
	AttrAppid = "appid"
	// AttrUnits ...
	AttrUnits = "units"
	// AttrLang ...
	AttrLang = "lang"
)

// NewSettings ...
func NewSettings() map[string]*m.Attribute {
	return map[string]*m.Attribute{
		AttrAppid: {
			Name: AttrAppid,
			Type: common.AttributeString,
		},
		AttrUnits: {
			Name: AttrUnits,
			Type: common.AttributeString,
		},
		AttrLang: {
			Name: AttrLang,
			Type: common.AttributeString,
		},
	}
}

//  "id": 500,
//  "main": "Rain",
//  "description": "небольшой дождь",
//  "icon": "10d"
func WeatherCondition(w ProductWeather) (state entity_manager.ActorState) {

	//fmt.Println("------")
	//debug.Println(w)

	var n, winter bool
	n = string(w.Icon[len(w.Icon)-1]) == "n"

	switch w.Id {

	// Thunderstorm

	//thunderstorm with light rain
	case 200:
		state = weather.GetActorState(weather.StateLightRainAndThunder, n, winter)
	//thunderstorm with rain
	case 201:
		state = weather.GetActorState(weather.StateRainAndThunder, n, winter)
	//thunderstorm with heavy rain
	case 202:
		state = weather.GetActorState(weather.StateHeavyRainAndThunder, n, winter)
	//light thunderstorm
	case 210:
		state = weather.GetActorState(weather.StateHeavyRainShowersAndThunder, n, winter)
	//thunderstorm
	case 211:
		state = weather.GetActorState(weather.StateHeavyRainShowersAndThunder, n, winter)
	//heavy thunderstorm
	case 212:
		state = weather.GetActorState(weather.StateHeavyRainShowersAndThunder, n, winter)
	//ragged thunderstorm
	case 221:
		state = weather.GetActorState(weather.StateHeavyRainShowersAndThunder, n, winter)
	//thunderstorm with light drizzle
	case 230:
		state = weather.GetActorState(weather.StateHeavyRainShowersAndThunder, n, winter)
	//thunderstorm with drizzle
	case 231:
		state = weather.GetActorState(weather.StateHeavyRainShowersAndThunder, n, winter)
	//thunderstorm with heavy drizzle
	case 232:
		state = weather.GetActorState(weather.StateHeavyRainShowersAndThunder, n, winter)

	// Drizzle

	//light intensity drizzle
	case 300:
		state = weather.GetActorState(weather.StateLightRainShowers, n, winter)
	//drizzle
	case 301:
		state = weather.GetActorState(weather.StateLightRainShowers, n, winter)
	//heavy intensity drizzle
	case 302:
		state = weather.GetActorState(weather.StateLightRainShowers, n, winter)
	//light intensity drizzle rain
	case 310:
		state = weather.GetActorState(weather.StateLightRainShowers, n, winter)
	//drizzle rain
	case 311:
		state = weather.GetActorState(weather.StateLightRainShowers, n, winter)
	//heavy intensity drizzle rain
	case 312:
		state = weather.GetActorState(weather.StateLightRainShowers, n, winter)
	//shower rain and drizzle
	case 313:
		state = weather.GetActorState(weather.StateLightRainShowers, n, winter)
	//heavy shower rain and drizzle
	case 314:
		state = weather.GetActorState(weather.StateLightRainShowers, n, winter)
	//shower drizzle
	case 321:
		state = weather.GetActorState(weather.StateLightRainShowers, n, winter)

	// Rain

	//light rain
	case 500:
		state = weather.GetActorState(weather.StateLightRain, n, winter)
	//moderate rain
	case 501:
		state = weather.GetActorState(weather.StateRain, n, winter)
	//heavy intensity rain
	case 502:
		state = weather.GetActorState(weather.StateHeavyRain, n, winter)
	//very heavy rain
	case 503:
		state = weather.GetActorState(weather.StateHeavySnowShowers, n, winter)
	//extreme rain
	case 504:
		state = weather.GetActorState(weather.StateHeavySnowShowers, n, winter)
	//freezing rain
	case 511:
		state = weather.GetActorState(weather.StateHeavySleetShowers, n, winter)
	//light intensity shower rain
	case 520:
		// todo fix
		state = weather.GetActorState(weather.StateHeavySnowShowers, n, winter)
	//shower rain
	case 521:
		// todo fix
		state = weather.GetActorState(weather.StateHeavySnowShowers, n, winter)
	//heavy intensity shower rain
	case 522:
		// todo fix
		state = weather.GetActorState(weather.StateHeavySnowShowers, n, winter)
	//ragged shower rain
	case 531:
		// todo fix
		state = weather.GetActorState(weather.StateHeavySnowShowers, n, winter)

	// Snow

	//light snow
	case 600:
		state = weather.GetActorState(weather.StateLightSnow, n, winter)
	//Snow
	case 601:
		state = weather.GetActorState(weather.StateSnow, n, winter)
	//Heavy snow
	case 602:
		state = weather.GetActorState(weather.StateHeavySnow, n, winter)
	//Sleet
	case 611:
		state = weather.GetActorState(weather.StateLightSleetShowers, n, winter)
	//Light shower sleet
	case 612:
		state = weather.GetActorState(weather.StateSleetShowers, n, winter)
	//Shower sleet
	case 613:
		state = weather.GetActorState(weather.StateHeavySleetShowers, n, winter)
	//Light rain and snow
	case 615:
		state = weather.GetActorState(weather.StateLightSleetShowers, n, winter)
	//Rain and snow
	case 616:
		state = weather.GetActorState(weather.StateLightSleetShowers, n, winter)
	//Light shower snow
	case 620:
		state = weather.GetActorState(weather.StateLightSnowShowers, n, winter)
	//Shower snow
	case 621:
		state = weather.GetActorState(weather.StateSnowShowers, n, winter)
	//Heavy shower snow
	case 622:
		state = weather.GetActorState(weather.StateHeavySnowShowers, n, winter)

	// Atmosphere

	//Mist
	case 701:
		state = weather.GetActorState(weather.StateFog, n, winter)
	//Smoke
	case 711:
		state = weather.GetActorState(weather.StateFog, n, winter)
	//Haze
	case 721:
		state = weather.GetActorState(weather.StateFog, n, winter)
	//Dust	sand/dust whirls
	case 731:
		state = weather.GetActorState(weather.StateFog, n, winter)
	//Fog
	case 741:
		state = weather.GetActorState(weather.StateFog, n, winter)
	//Sand
	case 751:
		state = weather.GetActorState(weather.StateFog, n, winter)
	//Dust
	case 761:
		state = weather.GetActorState(weather.StateFog, n, winter)
	//Ash volcanic ash
	case 762:
		state = weather.GetActorState(weather.StateFog, n, winter)
	//Squall
	case 771:
		state = weather.GetActorState(weather.StateFog, n, winter)
	//Tornado
	case 781:
		state = weather.GetActorState(weather.StateFog, n, winter)

	// Clear

	//clear sky
	case 800:
		state = weather.GetActorState(weather.StateClearSky, n, winter)

	// Clouds

	// few clouds: 11-25%
	case 801:
		state = weather.GetActorState(weather.StateFair, n, winter)
	// scattered clouds: 25-50%
	case 802:
		state = weather.GetActorState(weather.StateFair, n, winter)
	// broken clouds: 51-84%
	case 803:
		state = weather.GetActorState(weather.StatePartlyCloudy, n, winter)
	// overcast clouds: 85-100%
	case 804:
		state = weather.GetActorState(weather.StateCloudy, n, winter)

	default:
		log.Errorf("unknown weather id %d", w.Id)
	}

	//debug.Println(state)

	return
}
