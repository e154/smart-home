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
	}
	return
}

// UpdateArea ...
func (r Area) UpdateArea(obj *api.UpdateAreaRequest) (area *m.Area) {
	area = &m.Area{
		Id:          obj.Id,
		Name:        obj.Name,
		Description: obj.Description,
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
		Meta: &api.GetAreaListResult_Meta{
			Limit:        uint64(pagination.Limit),
			ObjectsCount: total,
			Offset:       uint64(pagination.Offset),
		},
	}
}

// ToArea ...
func (r Area) ToArea(area *m.Area) (obj *api.Area) {
	if area == nil {
		return
	}
	obj = &api.Area{
		Id:          area.Id,
		Name:        area.Name,
		Description: area.Description,
	}
	return
}
