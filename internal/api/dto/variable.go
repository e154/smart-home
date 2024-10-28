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
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/e154/smart-home/pkg/models"
)

// Variable ...
type Variable struct{}

// NewVariableDto ...
func NewVariableDto() Variable {
	return Variable{}
}

// AddVariable ...
func (r Variable) AddVariable(from *stub.ApiNewVariableRequest) (ver models.Variable) {
	ver = models.Variable{
		Name:  from.Name,
		Value: from.Value,
	}
	// tags
	for _, name := range from.Tags {
		ver.Tags = append(ver.Tags, &models.Tag{
			Name: name,
		})
	}
	return
}

// UpdateVariable ...
func (r Variable) UpdateVariable(obj *stub.VariableServiceUpdateVariableJSONBody, name string) (ver models.Variable) {
	ver = models.Variable{
		Name:  name,
		Value: obj.Value,
	}
	// tags
	for _, name := range obj.Tags {
		ver.Tags = append(ver.Tags, &models.Tag{
			Name: name,
		})
	}
	return
}

// ToListResult ...
func (r Variable) ToListResult(list []models.Variable) []*stub.ApiVariable {

	items := make([]*stub.ApiVariable, 0, len(list))

	for _, i := range list {
		items = append(items, r.ToVariable(i))
	}

	return items
}

// ToVariable ...
func (r Variable) ToVariable(ver models.Variable) (obj *stub.ApiVariable) {
	obj = ToVariable(ver)
	return
}

// ToSearchResult ...
func (r Variable) ToSearchResult(list []models.Variable) *stub.ApiSearchVariableResult {

	items := make([]stub.ApiVariable, 0, len(list))

	for _, v := range list {
		items = append(items, stub.ApiVariable{
			Name:      v.Name,
			System:    v.System,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return &stub.ApiSearchVariableResult{
		Items: items,
	}
}

// ToVariable ...
func ToVariable(ver models.Variable) (obj *stub.ApiVariable) {
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
	//tags
	for _, tag := range ver.Tags {
		obj.Tags = append(obj.Tags, tag.Name)
	}
	return
}
