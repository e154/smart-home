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

// DashboardCard ...
type DashboardCard struct{}

// NewDashboardCardDto ...
func NewDashboardCardDto() DashboardCard {
	return DashboardCard{}
}

func (r DashboardCard) AddDashboardCard(obj *api.NewDashboardCardRequest) (ver *m.DashboardCard) {
	ver = &m.DashboardCard{
		Title:          obj.Title,
		Height:         int(obj.Height),
		Width:          int(obj.Width),
		Background:     obj.Background,
		Weight:         int(obj.Weight),
		Enabled:        obj.Enabled,
		DashboardTabId: obj.DashboardTabId,
		Payload:        nil, //todo
	}
	return
}

func (r DashboardCard) UpdateDashboardCard(obj *api.UpdateDashboardCardRequest) (ver *m.DashboardCard) {
	ver = &m.DashboardCard{
		Id:             obj.Id,
		Title:          obj.Title,
		Height:         int(obj.Height),
		Width:          int(obj.Width),
		Background:     obj.Background,
		Weight:         int(obj.Weight),
		Enabled:        obj.Enabled,
		DashboardTabId: obj.DashboardTabId,
		Payload:        nil, //todo
	}
	return
}

// ToListResult ...
func (r DashboardCard) ToListResult(list []*m.DashboardCard, total uint64, pagination common.PageParams) *api.GetDashboardCardListResult {

	items := make([]*api.DashboardCard, 0, len(list))

	for _, i := range list {
		items = append(items, ToDashboardCard(i))
	}

	return &api.GetDashboardCardListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// ToDashboardCard ...
func (r DashboardCard) ToDashboardCard(ver *m.DashboardCard) (obj *api.DashboardCard) {
	obj = ToDashboardCard(ver)
	return
}

// ToDashboardCard ...
func ToDashboardCard(ver *m.DashboardCard) (obj *api.DashboardCard) {
	if ver == nil {
		return
	}
	obj = &api.DashboardCard{
		Id:             ver.Id,
		Title:          ver.Title,
		Height:         int32(ver.Height),
		Width:          int32(ver.Width),
		Background:     ver.Background,
		Weight:         int32(ver.Weight),
		Enabled:        ver.Enabled,
		DashboardTabId: ver.DashboardTabId,
		Payload:        "", //todo
		Entities:       make(map[string]*api.Entity),
		CreatedAt:      timestamppb.New(ver.CreatedAt),
		UpdatedAt:      nil,
	}

	// Items
	for _, item := range ver.Items {
		obj.Items = append(obj.Items, ToDashboardCardItem(item))
	}

	// Entities
	for key, entity := range ver.Entities {
		obj.Entities[key.String()] = ToEntity(entity)
	}

	return
}
