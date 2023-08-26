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

// Trigger ...
type Trigger struct{}

// NewTriggerDto ...
func NewTriggerDto() Trigger {
	return Trigger{}
}

// AddTrigger ...
func (r Trigger) AddTrigger(from *api.NewTriggerRequest) (action *m.Trigger) {
	action = &m.Trigger{
		Name:       from.Name,
		EntityId:   common.NewEntityIdFromPtr(from.EntityId),
		ScriptId:   from.ScriptId,
		PluginName: from.PluginName,
		Payload:    AttributeFromApi(from.Attributes),
	}
	return
}

// UpdateTrigger ...
func (r Trigger) UpdateTrigger(from *api.UpdateTriggerRequest) (action *m.Trigger) {
	action = &m.Trigger{
		Id:         from.Id,
		Name:       from.Name,
		EntityId:   common.NewEntityIdFromPtr(from.EntityId),
		ScriptId:   from.ScriptId,
		PluginName: from.PluginName,
		Payload:    AttributeFromApi(from.Attributes),
	}
	return
}

// ToSearchResult ...
func (r Trigger) ToSearchResult(list []*m.Trigger) *api.SearchTriggerResult {

	items := make([]*api.Trigger, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToTrigger(i))
	}

	return &api.SearchTriggerResult{
		Items: items,
	}
}

// ToListResult ...
func (r Trigger) ToListResult(list []*m.Trigger, total uint64, pagination common.PageParams) *api.GetTriggerListResult {

	items := make([]*api.Trigger, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToTrigger(i))
	}

	return &api.GetTriggerListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// ToTrigger ...
func (r Trigger) ToTrigger(action *m.Trigger) (obj *api.Trigger) {
	obj = ToTrigger(action)
	return
}

// ToTrigger ...
func ToTrigger(trigger *m.Trigger) (obj *api.Trigger) {
	if trigger == nil {
		return
	}
	obj = &api.Trigger{
		Id:         trigger.Id,
		Name:       trigger.Name,
		Entity:     ToEntity(trigger.Entity),
		EntityId:   trigger.EntityId.StringPtr(),
		Script:     ToScript(trigger.Script),
		ScriptId:   trigger.ScriptId,
		PluginName: trigger.PluginName,
		Attributes: AttributeToApi(trigger.Payload),
		CreatedAt:  timestamppb.New(trigger.CreatedAt),
		UpdatedAt:  timestamppb.New(trigger.UpdatedAt),
	}
	return
}
