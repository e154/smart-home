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

// Trigger ...
type Trigger struct{}

// NewTriggerDto ...
func NewTriggerDto() Trigger {
	return Trigger{}
}

// AddTrigger ...
func (r Trigger) AddTrigger(from *stub.ApiNewTriggerRequest) (action *m.Trigger) {
	action = &m.Trigger{
		Name:        from.Name,
		Description: from.Description,
		EntityId:    common.NewEntityIdFromPtr(from.EntityId),
		ScriptId:    from.ScriptId,
		AreaId:      from.AreaId,
		PluginName:  from.PluginName,
		Payload:     AttributeFromApi(from.Attributes),
		Enabled:     from.Enabled,
	}
	return
}

// UpdateTrigger ...
func (r Trigger) UpdateTrigger(from *stub.TriggerServiceUpdateTriggerJSONBody, id int64) (action *m.Trigger) {
	action = &m.Trigger{
		Id:          id,
		Name:        from.Name,
		Description: from.Description,
		EntityId:    common.NewEntityIdFromPtr(from.EntityId),
		ScriptId:    from.ScriptId,
		AreaId:      from.AreaId,
		PluginName:  from.PluginName,
		Payload:     AttributeFromApi(from.Attributes),
		Enabled:     from.Enabled,
	}
	return
}

// ToSearchResult ...
func (r Trigger) ToSearchResult(list []*m.Trigger) *stub.ApiSearchTriggerResult {

	items := make([]stub.ApiTrigger, 0, len(list))

	for _, i := range list {
		items = append(items, *r.ToTrigger(i))
	}

	return &stub.ApiSearchTriggerResult{
		Items: items,
	}
}

// ToListResult ...
func (r Trigger) ToListResult(list []*m.Trigger) []*stub.ApiTrigger {

	items := make([]*stub.ApiTrigger, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToTrigger(i))
	}

	return items
}

// ToTrigger ...
func (r Trigger) ToTrigger(action *m.Trigger) (obj *stub.ApiTrigger) {
	obj = ToTrigger(action)
	return
}

// ToTrigger ...
func ToTrigger(trigger *m.Trigger) (obj *stub.ApiTrigger) {
	if trigger == nil {
		return
	}
	obj = &stub.ApiTrigger{
		Id:          trigger.Id,
		Name:        trigger.Name,
		Description: trigger.Description,
		Entity:      ToEntity(trigger.Entity),
		EntityId:    trigger.EntityId.StringPtr(),
		Script:      GetStubScript(trigger.Script),
		ScriptId:    trigger.ScriptId,
		AreaId:      trigger.AreaId,
		Area:        GetStubArea(trigger.Area),
		PluginName:  trigger.PluginName,
		Enabled:     trigger.Enabled,
		IsLoaded:    common.Bool(trigger.IsLoaded),
		Attributes:  AttributeToApi(trigger.Payload),
		CreatedAt:   trigger.CreatedAt,
		UpdatedAt:   trigger.UpdatedAt,
	}
	return
}
