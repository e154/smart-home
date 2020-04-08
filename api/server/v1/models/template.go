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
type NewTemplate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	Status      string  `json:"status"`
	Type        string  `json:"type"`
	ParentName  *string `json:"parent"`
}

// swagger:model
type UpdateTemplate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	Status      string  `json:"status"`
	Type        string  `json:"type"`
	ParentName  *string `json:"parent"`
}

// swagger:model
type Template struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	ParentName  *string   `json:"parent"`
	Markers     []string  `json:"markers"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TemplateField ...
type TemplateField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// swagger:model
type TemplateContent struct {
	Items  []string         `json:"items"`
	Title  string           `json:"title"`
	Fields []*TemplateField `json:"fields"`
}
