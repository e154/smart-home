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
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"
)

// Area ...
type Area struct{}

// NewAreaDto ...
func NewAreaDto() Area {
	return Area{}
}

// AddArea ...
func (r Area) AddArea(from *stub.ApiNewAreaRequest) (area *models.Area) {
	area = &models.Area{
		Name:        from.Name,
		Description: from.Description,
		Polygon:     make([]models.Point, 0),
		Zoom:        from.Zoom,

		Resolution: from.Resolution,
	}
	if from.Center != nil {
		area.Center = models.Point{
			Lon: from.Center.Lon,
			Lat: from.Center.Lat,
		}
	}
	for _, point := range from.Polygon {
		area.Polygon = append(area.Polygon, models.Point{
			Lon: point.Lon,
			Lat: point.Lat,
		})
	}
	return
}

// UpdateArea ...
func (r Area) UpdateArea(from *stub.AreaServiceUpdateAreaJSONBody, id int64) (area *models.Area) {
	area = &models.Area{
		Id:          id,
		Name:        from.Name,
		Description: from.Description,
		Polygon:     make([]models.Point, 0),
		Zoom:        from.Zoom,
		Resolution:  from.Resolution,
	}
	if from.Center != nil {
		area.Center = models.Point{
			Lon: from.Center.Lon,
			Lat: from.Center.Lat,
		}
	}
	for _, point := range from.Polygon {
		area.Polygon = append(area.Polygon, models.Point{
			Lon: point.Lon,
			Lat: point.Lat,
		})
	}
	return
}

// ToSearchResult ...
func (r Area) ToSearchResult(list []*models.Area) stub.ApiSearchAreaResult {

	items := make([]stub.ApiArea, 0, len(list))

	for _, i := range list {
		items = append(items, *GetStubArea(i))
	}

	return stub.ApiSearchAreaResult{
		Items: items,
	}
}

// ToListResult ...
func (r Area) ToListResult(list []*models.Area) []stub.ApiArea {

	items := make([]stub.ApiArea, 0, len(list))

	for _, i := range list {
		items = append(items, *GetStubArea(i))
	}

	return items
}

// GetStubArea ...
func GetStubArea(area *models.Area) (obj *stub.ApiArea) {
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

func ImportArea(from *stub.ApiArea) (*int64, *models.Area) {
	if from == nil {
		return nil, nil
	}
	return common.Int64(from.Id), &models.Area{
		Id:          from.Id,
		Name:        from.Name,
		Description: from.Description,
	}
}
