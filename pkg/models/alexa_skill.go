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

	"github.com/e154/smart-home/pkg/common"
)

// AlexaSkill ...
type AlexaSkill struct {
	Id          int64             `json:"id"`
	SkillId     string            `json:"skill_id" validate:"required"`
	Description string            `json:"description"`
	Status      common.StatusType `json:"status" validate:"required"`
	Intents     []*AlexaIntent    `json:"intents"`
	Script      *Script           `json:"script"`
	ScriptId    *int64            `json:"script_id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
