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

// Dashboard ...
type Dashboard struct{}

// NewDashboardDto ...
func NewDashboardDto() Dashboard {
	return Dashboard{}
}

func (r Dashboard) AddDashboard(obj *stub.ApiNewDashboardRequest) (ver *m.Dashboard) {
	ver = &m.Dashboard{
		Name:        obj.Name,
		Description: obj.Description,
		Enabled:     obj.Enabled,
		AreaId:      obj.AreaId,
	}
	return
}

func (r Dashboard) UpdateDashboard(obj *stub.DashboardServiceUpdateDashboardJSONBody, id int64) (ver *m.Dashboard) {
	ver = &m.Dashboard{
		Id:          id,
		Name:        obj.Name,
		Description: obj.Description,
		Enabled:     obj.Enabled,
		AreaId:      obj.AreaId,
	}
	return
}

// ToListResult ...
func (r Dashboard) ToListResult(list []*m.Dashboard) []*stub.ApiDashboardShort {

	items := make([]*stub.ApiDashboardShort, 0, len(list))

	for _, i := range list {
		items = append(items, ToDashboardShort(i))
	}

	return items
}

// ToDashboard ...
func (r Dashboard) ToDashboard(ver *m.Dashboard) (obj *stub.ApiDashboard) {
	obj = ToDashboard(ver)
	return
}

// ToSearchResult ...
func (r Dashboard) ToSearchResult(list []*m.Dashboard) *stub.ApiSearchDashboardResult {

	items := make([]stub.ApiDashboard, 0, len(list))

	for _, i := range list {
		item := r.ToDashboard(i)
		items = append(items, *item)
	}

	return &stub.ApiSearchDashboardResult{
		Items: items,
	}
}

// ToDashboard ...
func ToDashboard(ver *m.Dashboard) (obj *stub.ApiDashboard) {
	if ver == nil {
		return
	}
	obj = &stub.ApiDashboard{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		Enabled:     ver.Enabled,
		AreaId:      ver.AreaId,
		Area:        GetStubArea(ver.Area),
		Tabs:        make([]stub.ApiDashboardTab, 0, len(ver.Tabs)),
		Entities:    make(map[string]stub.ApiEntity),
		CreatedAt:   ver.CreatedAt,
		UpdatedAt:   ver.UpdatedAt,
	}

	// Tabs
	for _, tab := range ver.Tabs {
		obj.Tabs = append(obj.Tabs, *ToDashboardTab(tab))
	}

	// Entities
	for key, entity := range ver.Entities {
		obj.Entities[key.String()] = *ToEntity(entity)
	}

	return
}

// ToDashboardShort ...
func ToDashboardShort(ver *m.Dashboard) (obj *stub.ApiDashboardShort) {
	if ver == nil {
		return
	}
	obj = &stub.ApiDashboardShort{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		Enabled:     ver.Enabled,
		AreaId:      ver.AreaId,
		Area:        GetStubArea(ver.Area),
		CreatedAt:   ver.CreatedAt,
		UpdatedAt:   ver.UpdatedAt,
	}
	return
}

// ImportDashboard ...
func ImportDashboard(obj *stub.ApiDashboard) (ver *m.Dashboard) {
	ver = &m.Dashboard{
		Id:          obj.Id,
		Name:        obj.Name,
		Description: obj.Description,
		Enabled:     obj.Enabled,
		AreaId:      obj.AreaId,
		Tabs:        make([]*m.DashboardTab, 0, len(obj.Tabs)),
	}

	// tabs
	for _, tabObj := range obj.Tabs {
		ver.Tabs = append(ver.Tabs, ImportDashboardTab(&tabObj))
	}

	return
}
