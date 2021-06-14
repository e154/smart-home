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

package models

import "time"

// swagger:model
type NewNode struct {
	Port           int64  `json:"port"`
	Status         string `json:"status"`
	Name           string `json:"name"`
	IP             string `json:"ip"`
	Description    string `json:"description"`
	Login          string `json:"login"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

// swagger:model
type UpdateNode struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Port           int64  `json:"port"`
	Status         string `json:"status"`
	IP             string `json:"ip"`
	Description    string `json:"description"`
	Login          string `json:"login"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

// swagger:model
type Node struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Ip          string    `json:"ip"`
	Port        int       `json:"port"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Login       string    `json:"login"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
