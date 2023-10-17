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
	m "github.com/e154/smart-home/models"
)

// DashboardTab ...
type DashboardTab struct{}

// NewDashboardTabDto ...
func NewDashboardTabDto() DashboardTab {
	return DashboardTab{}
}

func (r DashboardTab) AddDashboardTab(obj *stub.ApiNewDashboardTabRequest) (ver *m.DashboardTab) {
	ver = &m.DashboardTab{
		Name:        obj.Name,
		ColumnWidth: int(obj.ColumnWidth),
		Gap:         obj.Gap,
		Background:  obj.Background,
		Icon:        obj.Icon,
		Enabled:     obj.Enabled,
		Weight:      int(obj.Weight),
		DashboardId: obj.DashboardId,
	}
	return
}

func (r DashboardTab) UpdateDashboardTab(obj *stub.DashboardTabServiceUpdateDashboardTabJSONBody, id int64) (ver *m.DashboardTab) {
	ver = &m.DashboardTab{
		Id:          id,
		Name:        obj.Name,
		Icon:        obj.Icon,
		ColumnWidth: int(obj.ColumnWidth),
		Gap:         obj.Gap,
		Background:  obj.Background,
		Enabled:     obj.Enabled,
		Weight:      int(obj.Weight),
		DashboardId: obj.DashboardId,
	}
	return
}

// ToListResult ...
func (r DashboardTab) ToListResult(list []*m.DashboardTab) []*stub.ApiDashboardTab {

	items := make([]*stub.ApiDashboardTab, 0, len(list))

	for _, i := range list {
		items = append(items, ToDashboardTabShort(i))
	}

	return items
}

// ToDashboardTab ...
func (r DashboardTab) ToDashboardTab(ver *m.DashboardTab) (obj *stub.ApiDashboardTab) {
	obj = ToDashboardTab(ver)
	return
}

// ToDashboardTab ...
func ToDashboardTab(ver *m.DashboardTab) (obj *stub.ApiDashboardTab) {
	if ver == nil {
		return
	}
	obj = &stub.ApiDashboardTab{
		Id:          ver.Id,
		Name:        ver.Name,
		Icon:        ver.Icon,
		ColumnWidth: int32(ver.ColumnWidth),
		Gap:         ver.Gap,
		Background:  ver.Background,
		Enabled:     ver.Enabled,
		Weight:      int32(ver.Weight),
		DashboardId: ver.DashboardId,
		Cards:       make([]stub.ApiDashboardCard, 0, len(ver.Cards)),
		Entities:    make(map[string]stub.ApiEntity),
		CreatedAt:   ver.CreatedAt,
		UpdatedAt:   ver.UpdatedAt,
	}

	// Cards
	for _, card := range ver.Cards {
		obj.Cards = append(obj.Cards, *ToDashboardCard(card))
	}

	// Entities
	for key, entity := range ver.Entities {
		obj.Entities[key.String()] = *ToEntity(entity)
	}

	return
}

// ToDashboardTabShort ...
func ToDashboardTabShort(ver *m.DashboardTab) (obj *stub.ApiDashboardTab) {
	if ver == nil {
		return
	}
	obj = &stub.ApiDashboardTab{
		Id:          ver.Id,
		Name:        ver.Name,
		Icon:        ver.Icon,
		ColumnWidth: int32(ver.ColumnWidth),
		Gap:         ver.Gap,
		Background:  ver.Background,
		Enabled:     ver.Enabled,
		Weight:      int32(ver.Weight),
		DashboardId: ver.DashboardId,
		CreatedAt:   ver.CreatedAt,
		UpdatedAt:   ver.UpdatedAt,
	}
	return
}

func ImportDashboardTab(obj *stub.ApiDashboardTab) (ver *m.DashboardTab) {
	ver = &m.DashboardTab{
		Id:          obj.Id,
		Name:        obj.Name,
		ColumnWidth: int(obj.ColumnWidth),
		Gap:         obj.Gap,
		Background:  obj.Background,
		Icon:        obj.Icon,
		Enabled:     obj.Enabled,
		Weight:      int(obj.Weight),
		DashboardId: obj.DashboardId,
		Cards:       make([]*m.DashboardCard, 0, len(obj.Cards)),
	}

	// cards
	for _, cardObj := range obj.Cards {
		ver.Cards = append(ver.Cards, ImportDashboardCard(&cardObj))
	}

	return
}
