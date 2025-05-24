// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Point struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Polygon struct {
	Points []Point
}

func (p *Polygon) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	points, _ := formatPoints(p.Points)
	return clause.Expr{
		SQL:  "?::polygon",
		Vars: []interface{}{fmt.Sprintf("(%s)", points)},
	}
}

func (p *Polygon) Scan(src any) (err error) {
	value := fmt.Sprintf("%v", src)
	value = strings.ReplaceAll(value, "(", "[")
	value = strings.ReplaceAll(value, ")", "]")

	data := [][]float64{}
	if err = json.Unmarshal([]byte(value), &data); err != nil {
		return
	}
	for _, point := range data {
		p.Points = append(p.Points, Point{
			Lon: point[0],
			Lat: point[1],
		})
	}
	return nil
}

func (Polygon) GormDataType() string {
	return "polygon"
}

func formatPoints(polygonPoints []Point) (string, error) {
	var point []Point
	point = append(point, polygonPoints...)

	if len(point) > 0 &&
		point[0].Lon == point[len(point)-1].Lon &&
		point[0].Lat == point[len(point)-1].Lat {
		point = point[:len(point)-1]
	}

	if len(point) < 3 {
		return "", errors.New("polygon must have at least 3 unique points")
	}

	var points string
	for _, loc := range point {
		points += fmt.Sprintf("(%f, %f),", loc.Lon, loc.Lat)
	}

	points += fmt.Sprintf("(%f, %f)", point[0].Lon, point[0].Lat)

	return points, nil
}
