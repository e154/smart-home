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

package models

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/e154/smart-home/common"
)

// Variable ...
type Variable struct {
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Name      string           `json:"name"`
	Value     string           `json:"value"`
	EntityId  *common.EntityId `json:"entity_id"`
	System    bool             `json:"system"`
	Changed   bool
}

// NewVariable ...
func NewVariable(name string) Variable {
	return Variable{Name: name}
}

// GetObj ...
func (v *Variable) GetObj(obj interface{}) (err error) {
	err = json.Unmarshal([]byte(v.Value), obj)
	return
}

// SetObj ...
func (v *Variable) SetObj(obj interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(obj); err != nil {
		return
	}
	v.Value = string(b)
	return
}

// GetBool ...
func (v *Variable) GetBool() bool {
	b, _ := strconv.ParseBool(v.Value)
	return b
}
