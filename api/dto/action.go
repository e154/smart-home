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

// Action ...
type Action struct{}

// NewActionDto ...
func NewActionDto() Action {
	return Action{}
}

// Add ...
func (r Action) Add(from *stub.ApiNewActionRequest) (action *m.Action) {
	action = &m.Action{
		Name:             from.Name,
		ScriptId:         from.ScriptId,
		EntityId:         common.NewEntityIdFromPtr(from.EntityId),
		EntityActionName: from.EntityActionName,
	}
	return
}

// Update ...
func (r Action) Update(from *stub.ActionServiceUpdateActionJSONBody, id int64) (action *m.Action) {
	action = &m.Action{
		Id:               id,
		Name:             from.Name,
		ScriptId:         from.ScriptId,
		EntityId:         common.NewEntityIdFromPtr(from.EntityId),
		EntityActionName: from.EntityActionName,
	}
	return
}

// ToSearchResult ...
func (r Action) ToSearchResult(list []*m.Action) *stub.ApiSearchActionResult {

	items := make([]stub.ApiAction, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToAction(i))
	}

	return &stub.ApiSearchActionResult{
		Items: items,
	}
}

// ToListResult ...
func (r Action) ToListResult(list []*m.Action) []stub.ApiAction {

	items := make([]stub.ApiAction, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToAction(i))
	}

	return items
}

// ToAction ...
func (r Action) ToAction(action *m.Action) (obj stub.ApiAction) {
	obj = ToAction(action)
	return
}

// ToAction ...
func ToAction(action *m.Action) (obj stub.ApiAction) {
	if action == nil {
		return
	}
	obj = stub.ApiAction{
		Id:               action.Id,
		Name:             action.Name,
		ScriptId:         action.ScriptId,
		Script:           GetStubScript(action.Script),
		EntityId:         action.EntityId.StringPtr(),
		Entity:           ToEntity(action.Entity),
		EntityActionName: action.EntityActionName,
		CreatedAt:        action.CreatedAt,
		UpdatedAt:        action.UpdatedAt,
	}
	return
}
