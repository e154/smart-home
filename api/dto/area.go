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

package dto

import (
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Area ...
type Area struct{}

// NewAreaDto ...
func NewAreaDto() Area {
	return Area{}
}

// AddArea ...
func (r Area) AddArea(from *api.NewAreaRequest) (area *m.Area) {
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
func (r Area) UpdateArea(from *api.UpdateAreaRequest) (area *m.Area) {
	area = &m.Area{
		Id:          from.Id,
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
func (r Area) ToSearchResult(list []*m.Area) *api.SearchAreaResult {

	items := make([]*api.Area, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToArea(i))
	}

	return &api.SearchAreaResult{
		Items: items,
	}
}

// ToListResult ...
func (r Area) ToListResult(list []*m.Area, total uint64, pagination common.PageParams) *api.GetAreaListResult {

	items := make([]*api.Area, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToArea(i))
	}

	return &api.GetAreaListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// ToArea ...
func (r Area) ToArea(area *m.Area) (obj *api.Area) {
	obj = ToArea(area)
	return
}

// ToArea ...
func ToArea(area *m.Area) (obj *api.Area) {
	if area == nil {
		return
	}
	obj = &api.Area{
		Id:          area.Id,
		Name:        area.Name,
		Description: area.Description,
		Polygon:     make([]*api.AreaLocation, 0, len(area.Polygon)),
		Center: &api.AreaLocation{
			Lat: area.Center.Lat,
			Lon: area.Center.Lon,
		},
		Zoom:       area.Zoom,
		Resolution: area.Resolution,
		CreatedAt:  timestamppb.New(area.CreatedAt),
		UpdatedAt:  timestamppb.New(area.UpdatedAt),
	}
	for _, location := range area.Polygon {
		obj.Polygon = append(obj.Polygon, &api.AreaLocation{
			Lat: location.Lat,
			Lon: location.Lon,
		})
	}
	return
}

func ImportArea(from *api.Area) (*int64, *m.Area) {
	if from == nil {
		return nil, nil
	}
	return common.Int64(from.Id), &m.Area{
		Id:          from.Id,
		Name:        from.Name,
		Description: from.Description,
	}
}
