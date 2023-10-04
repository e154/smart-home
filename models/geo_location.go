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

package models

type Point struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// "status": "success",
// "country": "Russia",
// "countryCode": "RU",
// "region": "NVS",
// "regionName": "Novosibirsk Oblast",
// "city": "Novosibirsk",
// "zip": "630008",
// "lat": 54.9022,
// "lon": 83.0335,
// "timezone": "Asia/Novosibirsk",
// "isp": "Novotelecom Ltd.",
// "org": "Novotelecom ltd.",
// "as": "AS31200 Novotelecom Ltd",
// "query": "xxx.xxx.xxx.xxx"
type GeoLocation struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

// "ip": "xxx.xxx.xxx.xxx",
// "version": "IPv4",
// "city": "Novosibirsk",
// "region": "Novosibirsk Oblast",
// "region_code": "NVS",
// "country": "RU",
// "country_name": "Russia",
// "country_code": "RU",
// "country_code_iso3": "RUS",
// "country_capital": "Moscow",
// "country_tld": ".ru",
// "continent_code": "EU",
// "in_eu": false,
// "postal": "630009",
// "latitude": 54.9022,
// "longitude": 83.0335,
// "timezone": "Asia/Novosibirsk",
// "utc_offset": "+0700",
// "country_calling_code": "+7",
// "currency": "RUB",
// "currency_name": "Ruble",
// "languages": "ru,tt,xal,cau,ady,kv,ce,tyv,cv,udm,tut,mns,bua,myv,mdf,chm,ba,inh,tut,kbd,krc,av,sah,nog",
// "country_area": 17100000.0,
// "country_population": 144478050.0,
// "asn": "AS31200",
// "org": "Novotelecom Ltd"
type RegionInfo struct {
	Ip                 string  `json:"ip"`
	Version            string  `json:"version"`
	City               string  `json:"city"`
	Region             string  `json:"region"`
	RegionCode         string  `json:"region_code"`
	Country            string  `json:"country"`
	CountryName        string  `json:"country_name"`
	CountryCode        string  `json:"country_code"`
	CountryCodeIso3    string  `json:"country_code_iso3"`
	CountryCapital     string  `json:"country_capital"`
	CountryTld         string  `json:"country_tld"`
	ContinentCode      string  `json:"continent_code"`
	InEu               bool    `json:"in_eu"`
	Postal             string  `json:"postal"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Timezone           string  `json:"timezone"`
	UtcOffset          string  `json:"utc_offset"`
	CountryCallingCode string  `json:"country_calling_code"`
	Currency           string  `json:"currency"`
	CurrencyName       string  `json:"currency_name"`
	Languages          string  `json:"languages"`
	CountryArea        float64 `json:"country_area"`
	CountryPopulation  float64 `json:"country_population"`
	Asn                string  `json:"asn"`
	Org                string  `json:"org"`
}
