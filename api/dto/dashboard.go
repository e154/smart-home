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

// Dashboard ...
type Dashboard struct{}

// NewDashboardDto ...
func NewDashboardDto() Dashboard {
	return Dashboard{}
}

func (r Dashboard) AddDashboard(obj *api.NewDashboardRequest) (ver *m.Dashboard) {
	ver = &m.Dashboard{
		Name:        obj.Name,
		Description: obj.Description,
		Enabled:     obj.Enabled,
		AreaId:      obj.AreaId,
	}
	return
}

func (r Dashboard) UpdateDashboard(obj *api.UpdateDashboardRequest) (ver *m.Dashboard) {
	ver = &m.Dashboard{
		Id:          obj.Id,
		Name:        obj.Name,
		Description: obj.Description,
		Enabled:     obj.Enabled,
		AreaId:      obj.AreaId,
	}
	return
}

// ToListResult ...
func (r Dashboard) ToListResult(list []*m.Dashboard, total uint64, pagination common.PageParams) *api.GetDashboardListResult {

	items := make([]*api.DashboardShort, 0, len(list))

	for _, i := range list {
		items = append(items, ToDashboardShort(i))
	}

	return &api.GetDashboardListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// ToDashboard ...
func (r Dashboard) ToDashboard(ver *m.Dashboard) (obj *api.Dashboard) {
	obj = ToDashboard(ver)
	return
}

// ToSearchResult ...
func (r Dashboard) ToSearchResult(list []*m.Dashboard) *api.SearchDashboardResult {

	items := make([]*api.Dashboard, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToDashboard(i))
	}

	return &api.SearchDashboardResult{
		Items: items,
	}
}

// ToDashboard ...
func ToDashboard(ver *m.Dashboard) (obj *api.Dashboard) {
	if ver == nil {
		return
	}
	obj = &api.Dashboard{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		Enabled:     ver.Enabled,
		AreaId:      ver.AreaId,
		Area:        ToArea(ver.Area),
		Tabs:        make([]*api.DashboardTab, 0, len(ver.Tabs)),
		Entities:    make(map[string]*api.Entity),
		CreatedAt:   timestamppb.New(ver.CreatedAt),
		UpdatedAt:   timestamppb.New(ver.UpdatedAt),
	}

	// Tabs
	for _, tab := range ver.Tabs {
		obj.Tabs = append(obj.Tabs, ToDashboardTab(tab))
	}

	// Entities
	for key, entity := range ver.Entities {
		obj.Entities[key.String()] = ToEntity(entity)
	}

	return
}

// ToDashboardShort ...
func ToDashboardShort(ver *m.Dashboard) (obj *api.DashboardShort) {
	if ver == nil {
		return
	}
	obj = &api.DashboardShort{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		Enabled:     ver.Enabled,
		AreaId:      ver.AreaId,
		CreatedAt:   timestamppb.New(ver.CreatedAt),
		UpdatedAt:   timestamppb.New(ver.UpdatedAt),
	}
	return
}

// ImportDashboard ...
func ImportDashboard(obj *api.Dashboard) (ver *m.Dashboard) {
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
		ver.Tabs = append(ver.Tabs, ImportDashboardTab(tabObj))
	}

	return
}
