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

// DashboardCardItem ...
type DashboardCardItem struct{}

// NewDashboardCardItemDto ...
func NewDashboardCardItemDto() DashboardCardItem {
	return DashboardCardItem{}
}

func (r DashboardCardItem) AddDashboardCardItem(obj *stub.ApiNewDashboardCardItemRequest) (ver *m.DashboardCardItem) {
	ver = &m.DashboardCardItem{
		Title:           obj.Title,
		Type:            obj.Type,
		Weight:          int(obj.Weight),
		Enabled:         obj.Enabled,
		DashboardCardId: obj.DashboardCardId,
		Payload:         obj.Payload,
		Hidden:          obj.Hidden,
		Frozen:          obj.Frozen,
	}

	if obj.EntityId != nil && *obj.EntityId != "" {
		ver.EntityId = common.NewEntityId(*obj.EntityId)
	}
	return
}

func (r DashboardCardItem) UpdateDashboardCardItem(obj *stub.DashboardCardItemServiceUpdateDashboardCardItemJSONBody, id int64) (ver *m.DashboardCardItem) {
	ver = &m.DashboardCardItem{
		Id:              id,
		Title:           obj.Title,
		Type:            obj.Type,
		Weight:          int(obj.Weight),
		Enabled:         obj.Enabled,
		DashboardCardId: obj.DashboardCardId,
		Payload:         obj.Payload,
		Hidden:          obj.Hidden,
		Frozen:          obj.Frozen,
	}

	if obj.EntityId != nil && *obj.EntityId != "" {
		ver.EntityId = common.NewEntityId(*obj.EntityId)
	}
	return
}

// ToListResult ...
func (r DashboardCardItem) ToListResult(list []*m.DashboardCardItem) []*stub.ApiDashboardCardItem {

	items := make([]*stub.ApiDashboardCardItem, 0, len(list))

	for _, i := range list {
		items = append(items, ToDashboardCardItem(i))
	}

	return items
}

// ToDashboardCardItem ...
func (r DashboardCardItem) ToDashboardCardItem(ver *m.DashboardCardItem) (obj *stub.ApiDashboardCardItem) {
	obj = ToDashboardCardItem(ver)
	return
}

// ToDashboardCardItem ...
func ToDashboardCardItem(ver *m.DashboardCardItem) (obj *stub.ApiDashboardCardItem) {
	if ver == nil {
		return
	}
	obj = &stub.ApiDashboardCardItem{
		Id:              ver.Id,
		Title:           ver.Title,
		Type:            ver.Type,
		Weight:          int32(ver.Weight),
		Enabled:         ver.Enabled,
		DashboardCardId: ver.DashboardCardId,
		Payload:         ver.Payload,
		Hidden:          ver.Hidden,
		Frozen:          ver.Frozen,
		CreatedAt:       ver.CreatedAt,
		UpdatedAt:       ver.UpdatedAt,
	}

	if ver.EntityId != nil && *ver.EntityId != "" {
		obj.EntityId = common.String(string(*ver.EntityId))
	}

	return
}

func ImportDashboardCardItem(obj *stub.ApiDashboardCardItem) (ver *m.DashboardCardItem) {
	ver = &m.DashboardCardItem{
		Id:              obj.Id,
		Title:           obj.Title,
		Type:            obj.Type,
		Weight:          int(obj.Weight),
		Enabled:         obj.Enabled,
		DashboardCardId: obj.DashboardCardId,
		Payload:         obj.Payload,
		Hidden:          obj.Hidden,
		Frozen:          obj.Frozen,
	}
	if obj.EntityId != nil {
		ver.EntityId = common.NewEntityId(*obj.EntityId)
	}
	return
}
