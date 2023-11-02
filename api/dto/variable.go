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

// Variable ...
type Variable struct{}

// NewVariableDto ...
func NewVariableDto() Variable {
	return Variable{}
}

// AddVariable ...
func (r Variable) AddVariable(from *stub.ApiNewVariableRequest) (ver m.Variable) {
	ver = m.Variable{
		Name:  from.Name,
		Value: from.Value,
	}
	return
}

// UpdateVariable ...
func (r Variable) UpdateVariable(obj *stub.VariableServiceUpdateVariableJSONBody, name string) (ver m.Variable) {
	ver = m.Variable{
		Name:  name,
		Value: obj.Value,
	}
	return
}

// ToListResult ...
func (r Variable) ToListResult(list []m.Variable) []*stub.ApiVariable {

	items := make([]*stub.ApiVariable, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToVariable(i))
	}

	return items
}

// ToVariable ...
func (r Variable) ToVariable(ver m.Variable) (obj *stub.ApiVariable) {
	obj = ToVariable(ver)
	return
}

// ToVariable ...
func ToVariable(ver m.Variable) (obj *stub.ApiVariable) {
	if ver.Name == "" {
		return
	}
	obj = &stub.ApiVariable{
		Name:      ver.Name,
		Value:     ver.Value,
		System:    ver.System,
		CreatedAt: ver.CreatedAt,
		UpdatedAt: ver.UpdatedAt,
	}
	return
}
