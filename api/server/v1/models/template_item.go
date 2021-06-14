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
type NewTemplateItem struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	Status      string  `json:"status"`
	Type        string  `json:"type"`
	ParentName  *string `json:"parent"`
}

// swagger:model
type UpdateTemplateItem struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	Status      string  `json:"status"`
	Type        string  `json:"type"`
	ParentName  *string `json:"parent"`
}

// swagger:model
type UpdateTemplateItemStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// swagger:model
type TemplateItem struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	ParentName  *string   `json:"parent"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// swagger:model
type TemplateTreeBranch struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Status      string                `json:"status"`
	Nodes       []*TemplateTreeBranch `json:"nodes"`
}

// swagger:model
type UpdateTemplateTree []*TemplateTreeBranch

// swagger:model
type TemplateTree []*TemplateTreeBranch
