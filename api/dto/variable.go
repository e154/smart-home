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

// Variable ...
type Variable struct{}

// NewVariableDto ...
func NewVariableDto() Variable {
	return Variable{}
}

// AddVariable ...
func (r Variable) AddVariable(from *api.NewVariableRequest) (ver m.Variable) {
	ver = m.Variable{
		Name:  from.Name,
		Value: from.Value,
	}
	return
}

// UpdateVariable ...
func (r Variable) UpdateVariable(obj *api.UpdateVariableRequest) (ver m.Variable) {
	ver = m.Variable{
		Name:  obj.Name,
		Value: obj.Value,
	}
	return
}

// ToListResult ...
func (r Variable) ToListResult(list []m.Variable, total uint64, pagination common.PageParams) *api.GetVariableListResult {

	items := make([]*api.Variable, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToVariable(i))
	}

	return &api.GetVariableListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// ToVariable ...
func (r Variable) ToVariable(ver m.Variable) (obj *api.Variable) {
	obj = ToVariable(ver)
	return
}

// ToVariable ...
func ToVariable(ver m.Variable) (obj *api.Variable) {
	if ver.Name == "" {
		return
	}
	obj = &api.Variable{
		Name:      ver.Name,
		Value:     ver.Value,
		System:    ver.System,
		CreatedAt: timestamppb.New(ver.CreatedAt),
		UpdatedAt: timestamppb.New(ver.UpdatedAt),
	}
	return
}
