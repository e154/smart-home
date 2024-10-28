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
	"time"
)

// TriggerPayload ...
type TriggerPayload struct {
	Obj Attributes `json:"obj"`
}

// Trigger ...
type Trigger struct {
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Name        string     `json:"name" validate:"required,lte=255"`
	PluginName  string     `json:"plugin_name" validate:"required,lte=255"`
	Description string     `json:"description"`
	Id          int64      `json:"id"`
	Entities    []*Entity  `json:"entities" validate:"dive"`
	Script      *Script    `json:"script"`
	ScriptId    *int64     `json:"script_id"`
	Payload     Attributes `json:"payload"`
	AreaId      *int64     `json:"area_id"`
	Area        *Area      `json:"area"`
	Enabled     bool       `json:"enabled"`
	IsLoaded    bool       `json:"is_loaded"`
}

// NewTrigger ...
type NewTrigger struct {
	Name        string     `json:"name" validate:"required,lte=255"`
	PluginName  string     `json:"plugin_name" validate:"required,lte=255"`
	Description string     `json:"description"`
	EntityIds   []string   `json:"entity_ids"`
	ScriptId    *int64     `json:"script_id"`
	Payload     Attributes `json:"payload"`
	AreaId      *int64     `json:"area_id"`
	Enabled     bool       `json:"enabled"`
	IsLoaded    bool       `json:"is_loaded"`
}

// UpdateTrigger ...
type UpdateTrigger struct {
	Name        string     `json:"name" validate:"required,lte=255"`
	PluginName  string     `json:"plugin_name" validate:"required,lte=255"`
	Description string     `json:"description"`
	Id          int64      `json:"id"`
	EntityIds   []string   `json:"entity_ids" validate:"dive"`
	ScriptId    *int64     `json:"script_id"`
	Payload     Attributes `json:"payload"`
	AreaId      *int64     `json:"area_id"`
	Enabled     bool       `json:"enabled"`
	IsLoaded    bool       `json:"is_loaded"`
}
