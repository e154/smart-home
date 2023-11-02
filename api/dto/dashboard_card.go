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

// DashboardCard ...
type DashboardCard struct{}

// NewDashboardCardDto ...
func NewDashboardCardDto() DashboardCard {
	return DashboardCard{}
}

func (r DashboardCard) AddDashboardCard(obj *stub.ApiNewDashboardCardRequest) (ver *m.DashboardCard) {
	ver = &m.DashboardCard{
		Title:          obj.Title,
		Height:         int(obj.Height),
		Width:          int(obj.Width),
		Background:     obj.Background,
		Weight:         int(obj.Weight),
		Enabled:        obj.Enabled,
		DashboardTabId: obj.DashboardTabId,
		Payload:        obj.Payload,
		Hidden:         obj.Hidden,
	}
	if obj.EntityId != nil && *obj.EntityId != "" {
		ver.EntityId = common.NewEntityId(*obj.EntityId)
	}
	return
}

func (r DashboardCard) UpdateDashboardCard(obj *stub.DashboardCardServiceUpdateDashboardCardJSONBody, id int64) (ver *m.DashboardCard) {
	ver = &m.DashboardCard{
		Id:             id,
		Title:          obj.Title,
		Height:         int(obj.Height),
		Width:          int(obj.Width),
		Background:     obj.Background,
		Weight:         int(obj.Weight),
		Enabled:        obj.Enabled,
		Hidden:         obj.Hidden,
		DashboardTabId: obj.DashboardTabId,
		Payload:        obj.Payload,
		Items:          make([]*m.DashboardCardItem, 0, len(obj.Items)),
	}
	if obj.EntityId != nil && *obj.EntityId != "" {
		ver.EntityId = common.NewEntityId(*obj.EntityId)
	}
	// items
	for _, item := range obj.Items {
		qwe := &m.DashboardCardItem{
			Id:              item.Id,
			Title:           item.Title,
			Type:            item.Type,
			Weight:          int(item.Weight),
			Enabled:         item.Enabled,
			DashboardCardId: id,
			Payload:         item.Payload,
			Hidden:          item.Hidden,
			Frozen:          item.Frozen,
		}
		if item.EntityId != nil {
			qwe.EntityId = common.NewEntityId(*item.EntityId)
		}
		ver.Items = append(ver.Items, qwe)
	}

	return
}

// ToListResult ...
func (r DashboardCard) ToListResult(list []*m.DashboardCard) []*stub.ApiDashboardCard {

	items := make([]*stub.ApiDashboardCard, 0, len(list))

	for _, i := range list {
		items = append(items, ToDashboardCard(i))
	}

	return items
}

// ToDashboardCard ...
func ToDashboardCard(ver *m.DashboardCard) (obj *stub.ApiDashboardCard) {
	if ver == nil {
		return
	}
	obj = &stub.ApiDashboardCard{
		Id:             ver.Id,
		Title:          ver.Title,
		Height:         int32(ver.Height),
		Width:          int32(ver.Width),
		Background:     ver.Background,
		Weight:         int32(ver.Weight),
		Enabled:        ver.Enabled,
		DashboardTabId: ver.DashboardTabId,
		Payload:        ver.Payload,
		Hidden:         ver.Hidden,
		Entities:       make(map[string]stub.ApiEntity),
		CreatedAt:      ver.CreatedAt,
		UpdatedAt:      ver.UpdatedAt,
	}

	if ver.EntityId != nil && *ver.EntityId != "" {
		obj.EntityId = common.String(string(*ver.EntityId))
	}

	// Items
	for _, item := range ver.Items {
		cardItem := ToDashboardCardItem(item)
		obj.Items = append(obj.Items, *cardItem)
	}

	// Entities
	for key, entity := range ver.Entities {
		obj.Entities[key.String()] = *ToEntity(entity)
	}

	return
}

func ImportDashboardCard(obj *stub.ApiDashboardCard) (ver *m.DashboardCard) {
	if obj == nil {
		return
	}
	ver = &m.DashboardCard{
		Id:             obj.Id,
		Title:          obj.Title,
		Height:         int(obj.Height),
		Width:          int(obj.Width),
		Background:     obj.Background,
		Weight:         int(obj.Weight),
		Enabled:        obj.Enabled,
		DashboardTabId: obj.DashboardTabId,
		Payload:        obj.Payload,
		Hidden:         obj.Hidden,
		Items:          make([]*m.DashboardCardItem, 0, len(obj.Items)),
	}

	if obj.EntityId != nil && *obj.EntityId != "" {
		ver.EntityId = common.NewEntityId(*obj.EntityId)
	}

	// items
	for _, itemObj := range obj.Items {
		ver.Items = append(ver.Items, ImportDashboardCardItem(&itemObj))
	}

	return
}
