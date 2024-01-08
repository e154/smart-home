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
	m "github.com/e154/smart-home/models"
)

// Condition ...
type Condition struct{}

// NewConditionDto ...
func NewConditionDto() Condition {
	return Condition{}
}

// AddCondition ...
func (r Condition) AddCondition(from *stub.ApiNewConditionRequest) (action *m.Condition) {
	action = &m.Condition{
		Name:        from.Name,
		Description: from.Description,
		ScriptId:    from.ScriptId,
		AreaId:      from.AreaId,
	}
	return
}

// UpdateCondition ...
func (r Condition) UpdateCondition(from *stub.ConditionServiceUpdateConditionJSONBody, id int64) (action *m.Condition) {
	action = &m.Condition{
		Id:          id,
		Name:        from.Name,
		Description: from.Description,
		ScriptId:    from.ScriptId,
		AreaId:      from.AreaId,
	}
	return
}

// ToSearchResult ...
func (r Condition) ToSearchResult(list []*m.Condition) *stub.ApiSearchConditionResult {

	items := make([]stub.ApiCondition, 0, len(list))

	for _, i := range list {
		condition := ToCondition(i)
		items = append(items, *condition)
	}

	return &stub.ApiSearchConditionResult{
		Items: items,
	}
}

// ToListResult ...
func (r Condition) ToListResult(list []*m.Condition) []*stub.ApiCondition {

	items := make([]*stub.ApiCondition, 0, len(list))

	for _, cond := range list {
		items = append(items, &stub.ApiCondition{
			Id:          cond.Id,
			Name:        cond.Name,
			Description: cond.Description,
			ScriptId:    cond.ScriptId,
			AreaId:      cond.AreaId,
			Area:        GetStubArea(cond.Area),
			CreatedAt:   cond.CreatedAt,
			UpdatedAt:   cond.UpdatedAt,
		})
	}

	return items
}

// ToCondition ...
func ToCondition(cond *m.Condition) (obj *stub.ApiCondition) {
	if cond == nil {
		return
	}
	obj = &stub.ApiCondition{
		Id:          cond.Id,
		Name:        cond.Name,
		Description: cond.Description,
		ScriptId:    cond.ScriptId,
		AreaId:      cond.AreaId,
		Script:      GetStubScript(cond.Script),
		Area:        GetStubArea(cond.Area),
		CreatedAt:   cond.CreatedAt,
		UpdatedAt:   cond.UpdatedAt,
	}

	return
}
