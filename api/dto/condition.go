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

// Condition ...
type Condition struct{}

// NewConditionDto ...
func NewConditionDto() Condition {
	return Condition{}
}

// AddCondition ...
func (r Condition) AddCondition(from *api.NewConditionRequest) (action *m.Condition) {
	action = &m.Condition{
		Name:     from.Name,
		ScriptId: from.ScriptId,
	}
	return
}

// UpdateCondition ...
func (r Condition) UpdateCondition(from *api.UpdateConditionRequest) (action *m.Condition) {
	action = &m.Condition{
		Id:       from.Id,
		Name:     from.Name,
		ScriptId: from.ScriptId,
	}
	return
}

// ToSearchResult ...
func (r Condition) ToSearchResult(list []*m.Condition) *api.SearchConditionResult {

	items := make([]*api.Condition, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToCondition(i))
	}

	return &api.SearchConditionResult{
		Items: items,
	}
}

// ToListResult ...
func (r Condition) ToListResult(list []*m.Condition, total uint64, pagination common.PageParams) *api.GetConditionListResult {

	items := make([]*api.Condition, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToCondition(i))
	}

	return &api.GetConditionListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// ToCondition ...
func (r Condition) ToCondition(action *m.Condition) (obj *api.Condition) {
	obj = ToCondition(action)
	return
}

// ToCondition ...
func ToCondition(cond *m.Condition) (obj *api.Condition) {
	if cond == nil {
		return
	}
	obj = &api.Condition{
		Id:        cond.Id,
		Name:      cond.Name,
		ScriptId:  cond.ScriptId,
		Script:    ToScript(cond.Script),
		CreatedAt: timestamppb.New(cond.CreatedAt),
		UpdatedAt: timestamppb.New(cond.UpdatedAt),
	}

	return
}
