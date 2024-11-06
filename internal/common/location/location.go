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

package location

import (
	"encoding/json"
	"fmt"

	web2 "github.com/e154/smart-home/internal/system/web"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/web"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

const (
	// IpApi ...
	IpApi = "http://ip-api.com/json"
	// IPAPI ...
	IPAPI = "https://ipapi.co/json/"
)

// GeoLocationFromIP ...
func GeoLocationFromIP(ip string) (location m.GeoLocation, err error) {

	crawler := web2.New()

	var body []byte
	if _, body, err = crawler.Probe(web.Request{Method: "GET", Url: fmt.Sprintf("%s/%s", IpApi, ip)}); err != nil {
		return
	}
	location = m.GeoLocation{}
	err = json.Unmarshal(body, &location)

	return
}

// GetRegionInfo ...
func GetRegionInfo() (info m.RegionInfo, err error) {

	crawler := web2.New()

	var body []byte
	if _, body, err = crawler.Probe(web.Request{Method: "GET", Url: IPAPI}); err != nil {
		return
	}
	info = m.RegionInfo{}
	err = json.Unmarshal(body, &info)

	return
}

const (
	earthRadiusKm    = 6371.0
	earthRadiusMiles = 3959.0
)

const (
	Miles      = string("M")
	Kilometers = string("K")
)

// GetDistanceBetweenPoints ...
func GetDistanceBetweenPoints(point1, point2 m.Point, unit ...string) (distance float64) {

	point := s2.PointFromLatLng(s2.LatLngFromDegrees(point1.Lat, point1.Lon))
	minDistance := point.Distance(
		s2.PointFromLatLng(
			s2.LatLngFromDegrees(point2.Lat, point2.Lon),
		),
	)

	distance = angleToDistance(minDistance, unit...)

	return
}

// GetDistanceToPolygon ...
func GetDistanceToPolygon(point1 m.Point, polygon1 []m.Point, unit ...string) (distance float64) {

	point := s2.PointFromLatLng(s2.LatLngFromDegrees(point1.Lat, point1.Lon))

	var points []s2.Point
	for _, point := range polygon1 {
		points = append(points, s2.PointFromLatLng(s2.LatLngFromDegrees(point.Lat, point.Lon)))
	}

	loop := s2.LoopFromPoints(points)

	var minDistance s1.Angle
	minDistanceSet := false
	for _, vertex := range loop.Vertices() {
		distance := point.Distance(s2.PointFromLatLng(s2.LatLngFromPoint(vertex)))
		if !minDistanceSet || distance < minDistance {
			minDistance = distance
			minDistanceSet = true
		}
	}

	distance = angleToDistance(minDistance, unit...)

	return
}

func PointInsidePolygon(point1 m.Point, polygon1 []m.Point) bool {

	point := s2.PointFromLatLng(s2.LatLngFromDegrees(point1.Lat, point1.Lon))

	var points []s2.Point
	for _, point := range polygon1 {
		points = append(points, s2.PointFromLatLng(s2.LatLngFromDegrees(point.Lat, point.Lon)))
	}

	loop := s2.LoopFromPoints(points)

	return !loop.ContainsPoint(point)
}

// Преобразовать угол в расстояние с учетом радиуса Земли.
func angleToDistance(angle s1.Angle, unit ...string) (distance float64) {
	if len(unit) > 0 {
		switch unit[0] {
		case Miles:
			distance = float64(angle.Radians()) * earthRadiusMiles
		case Kilometers:
			distance = float64(angle.Radians()) * earthRadiusKm
		}
	} else {
		distance = float64(angle.Radians()) * earthRadiusKm
	}
	return
}
