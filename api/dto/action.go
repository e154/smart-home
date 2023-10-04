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
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Action ...
type Action struct{}

// NewActionDto ...
func NewActionDto() Action {
	return Action{}
}

// AddAction ...
func (r Action) AddAction(from *api.NewActionRequest) (action *m.Action) {
	action = &m.Action{
		Name:             from.Name,
		ScriptId:         from.ScriptId,
		EntityId:         common.NewEntityIdFromPtr(from.EntityId),
		EntityActionName: from.EntityActionName,
	}
	return
}

// UpdateAction ...
func (r Action) UpdateAction(from *api.UpdateActionRequest) (action *m.Action) {
	action = &m.Action{
		Id:               from.Id,
		Name:             from.Name,
		ScriptId:         from.ScriptId,
		EntityId:         common.NewEntityIdFromPtr(from.EntityId),
		EntityActionName: from.EntityActionName,
	}
	return
}

// ToSearchResult ...
func (r Action) ToSearchResult(list []*m.Action) *api.SearchActionResult {

	items := make([]*api.Action, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToAction(i))
	}

	return &api.SearchActionResult{
		Items: items,
	}
}

// ToListResult ...
func (r Action) ToListResult(list []*m.Action, total uint64, pagination common.PageParams) *api.GetActionListResult {

	items := make([]*api.Action, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToAction(i))
	}

	return &api.GetActionListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// ToAction ...
func (r Action) ToAction(action *m.Action) (obj *api.Action) {
	obj = ToAction(action)
	return
}

// ToAction ...
func ToAction(action *m.Action) (obj *api.Action) {
	if action == nil {
		return
	}
	obj = &api.Action{
		Id:               action.Id,
		Name:             action.Name,
		ScriptId:         action.ScriptId,
		Script:           ToScript(action.Script),
		EntityId:         action.EntityId.StringPtr(),
		Entity:           ToEntity(action.Entity),
		EntityActionName: action.EntityActionName,
		CreatedAt:        timestamppb.New(action.CreatedAt),
		UpdatedAt:        timestamppb.New(action.UpdatedAt),
	}
	return
}
