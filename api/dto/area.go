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

package dto

import (
	stub "github.com/e154/smart-home/api/stub"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// Area ...
type Area struct{}

// NewAreaDto ...
func NewAreaDto() Area {
	return Area{}
}

// AddArea ...
func (r Area) AddArea(from *stub.ApiNewAreaRequest) (area *m.Area) {
	area = &m.Area{
		Name:        from.Name,
		Description: from.Description,
		Polygon:     make([]m.Point, 0),
		Zoom:        from.Zoom,

		Resolution: from.Resolution,
	}
	if from.Center != nil {
		area.Center = m.Point{
			Lon: from.Center.Lon,
			Lat: from.Center.Lat,
		}
	}
	for _, point := range from.Polygon {
		area.Polygon = append(area.Polygon, m.Point{
			Lon: point.Lon,
			Lat: point.Lat,
		})
	}
	return
}

// UpdateArea ...
func (r Area) UpdateArea(from *stub.AreaServiceUpdateAreaJSONBody, id int64) (area *m.Area) {
	area = &m.Area{
		Id:          id,
		Name:        from.Name,
		Description: from.Description,
		Polygon:     make([]m.Point, 0),
		Zoom:        from.Zoom,
		Resolution:  from.Resolution,
	}
	if from.Center != nil {
		area.Center = m.Point{
			Lon: from.Center.Lon,
			Lat: from.Center.Lat,
		}
	}
	for _, point := range from.Polygon {
		area.Polygon = append(area.Polygon, m.Point{
			Lon: point.Lon,
			Lat: point.Lat,
		})
	}
	return
}

// ToSearchResult ...
func (r Area) ToSearchResult(list []*m.Area) stub.ApiSearchAreaResult {

	items := make([]stub.ApiArea, 0, len(list))

	for _, i := range list {
		items = append(items, *ToArea(i))
	}

	return stub.ApiSearchAreaResult{
		Items: items,
	}
}

// ToListResult ...
func (r Area) ToListResult(list []*m.Area) []stub.ApiArea {

	items := make([]stub.ApiArea, 0, len(list))

	for _, i := range list {
		items = append(items, *ToArea(i))
	}

	return items
}

// ToArea ...
func ToArea(area *m.Area) (obj *stub.ApiArea) {
	if area == nil {
		return
	}
	obj = &stub.ApiArea{
		Id:          area.Id,
		Name:        area.Name,
		Description: area.Description,
		Polygon:     make([]stub.ApiAreaLocation, 0, len(area.Polygon)),
		Center: &stub.ApiAreaLocation{
			Lat: area.Center.Lat,
			Lon: area.Center.Lon,
		},
		Zoom:       area.Zoom,
		Resolution: area.Resolution,
		CreatedAt:  area.CreatedAt,
		UpdatedAt:  area.UpdatedAt,
	}
	for _, location := range area.Polygon {
		obj.Polygon = append(obj.Polygon, stub.ApiAreaLocation{
			Lat: location.Lat,
			Lon: location.Lon,
		})
	}
	return
}

func ImportArea(from *stub.ApiArea) (*int64, *m.Area) {
	if from == nil {
		return nil, nil
	}
	return common.Int64(from.Id), &m.Area{
		Id:          from.Id,
		Name:        from.Name,
		Description: from.Description,
	}
}
