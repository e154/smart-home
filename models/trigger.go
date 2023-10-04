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

	"github.com/e154/smart-home/common"
)

// TriggerPayload ...
type TriggerPayload struct {
	Obj Attributes `json:"obj"`
}

// Trigger ...
type Trigger struct {
	Id         int64            `json:"id"`
	Name       string           `json:"name" validate:"required,lte=255"`
	Entity     *Entity          `json:"entity"`
	EntityId   *common.EntityId `json:"entity_id"`
	Script     *Script          `json:"script"`
	ScriptId   *int64           `json:"script_id"`
	PluginName string           `json:"plugin_name" validate:"required,lte=255"`
	Payload    Attributes       `json:"payload"`
	Enabled    bool             `json:"enabled"`
	IsLoaded   bool             `json:"is_loaded"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
