// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

import "time"

// swagger:model
type NewRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parent      struct {
		Name string `json:"name"`
	} `json:"parent"`
}

// swagger:model
type UpdateRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parent      struct {
		Name string `json:"name"`
	} `json:"parent"`
}

// swagger:model
type Role struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parent      struct {
		Name string `json:"name"`
	} `json:"parent"`
	Children   []*Role             `json:"children"`
	AccessList map[string][]string `json:"access_list"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

// AccessItem ...
type AccessItem struct {
	Actions     []string `json:"actions"`
	Method      string   `json:"method"`
	Description string   `json:"description"`
	RoleName    string   `json:"role_name"`
}

// AccessLevels ...
type AccessLevels map[string]AccessItem

// swagger:model
type AccessList map[string]AccessLevels

// swagger:model
type AccessListDiff map[string]map[string]bool
