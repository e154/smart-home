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

// DashboardCardItem ...
type DashboardCardItem struct{}

// NewDashboardCardItemDto ...
func NewDashboardCardItemDto() DashboardCardItem {
	return DashboardCardItem{}
}

func (r DashboardCardItem) AddDashboardCardItem(obj *api.NewDashboardCardItemRequest) (ver *m.DashboardCardItem) {
	ver = &m.DashboardCardItem{
		Title:           obj.Title,
		Type:            obj.Type,
		Weight:          int(obj.Weight),
		Enabled:         obj.Enabled,
		DashboardCardId: obj.DashboardCardId,
		Payload:         obj.Payload,
	}

	if obj.EntityId != nil && *obj.EntityId != "" {
		ver.EntityId = common.NewEntityId(*obj.EntityId)
	}
	return
}

func (r DashboardCardItem) UpdateDashboardCardItem(obj *api.UpdateDashboardCardItemRequest) (ver *m.DashboardCardItem) {
	ver = &m.DashboardCardItem{
		Id:              obj.Id,
		Title:           obj.Title,
		Type:            obj.Type,
		Weight:          int(obj.Weight),
		Enabled:         obj.Enabled,
		DashboardCardId: obj.DashboardCardId,
		Payload:         obj.Payload,
	}

	if obj.EntityId != nil && *obj.EntityId != "" {
		ver.EntityId = common.NewEntityId(*obj.EntityId)
	}
	return
}

// ToListResult ...
func (r DashboardCardItem) ToListResult(list []*m.DashboardCardItem, total uint64, pagination common.PageParams) *api.GetDashboardCardItemListResult {

	items := make([]*api.DashboardCardItem, 0, len(list))

	for _, i := range list {
		items = append(items, ToDashboardCardItem(i))
	}

	return &api.GetDashboardCardItemListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// ToDashboardCardItem ...
func (r DashboardCardItem) ToDashboardCardItem(ver *m.DashboardCardItem) (obj *api.DashboardCardItem) {
	obj = ToDashboardCardItem(ver)
	return
}

// ToDashboardCardItem ...
func ToDashboardCardItem(ver *m.DashboardCardItem) (obj *api.DashboardCardItem) {
	if ver == nil {
		return
	}
	obj = &api.DashboardCardItem{
		Id:              ver.Id,
		Title:           ver.Title,
		Type:            ver.Type,
		Weight:          int32(ver.Weight),
		Enabled:         ver.Enabled,
		DashboardCardId: ver.DashboardCardId,
		Payload:         ver.Payload,
		CreatedAt:       timestamppb.New(ver.CreatedAt),
		UpdatedAt:       timestamppb.New(ver.UpdatedAt),
	}

	if obj.EntityId != nil && *obj.EntityId != "" {
		ver.EntityId = common.NewEntityId(*obj.EntityId)
	}

	return
}
