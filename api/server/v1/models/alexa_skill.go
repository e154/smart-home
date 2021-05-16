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
type NewAlexaSkill struct {
	SkillId              string `json:"skill_id"`
	Description          string `json:"description"`
	Status               string `json:"status"`
	OnLaunchScriptId     *int64 `json:"on_launch_script_id,omitempty"`
	OnSessionEndScriptId *int64 `json:"on_session_end_script_id,omitempty"`
}

// swagger:model
type UpdateAlexaSkill struct {
	Id                   int64          `json:"id"`
	SkillId              string         `json:"skill_id"`
	Description          string         `json:"description"`
	Status               string         `json:"status"`
	Intents              []*AlexaIntent `json:"intents,omitempty"`
	OnLaunchScriptId     *int64         `json:"on_launch_script_id,omitempty"`
	OnSessionEndScriptId *int64         `json:"on_session_end_script_id,omitempty"`
}

// swagger:model
type AlexaSkill struct {
	Id                   int64          `json:"id"`
	SkillId              string         `json:"skill_id"`
	Description          string         `json:"description"`
	Status               string         `json:"status"`
	Intents              []*AlexaIntent `json:"intents,omitempty"`
	OnLaunchScript       *Script        `json:"on_launch_script,omitempty"`
	OnLaunchScriptId     *int64         `json:"on_launch_script_id,omitempty"`
	OnSessionEndScript   *Script        `json:"on_session_end_script,omitempty"`
	OnSessionEndScriptId *int64         `json:"on_session_end_script_id,omitempty"`
	CreatedAt            time.Time      `json:"created_at,omitempty"`
	UpdatedAt            time.Time      `json:"updated_at,omitempty"`
}

// swagger:model
type AlexaSkillShort struct {
	Id                   int64     `json:"id"`
	SkillId              string    `json:"skill_id"`
	Description          string    `json:"description"`
	Status               string    `json:"status"`
	OnLaunchScriptId     *int64    `json:"on_launch_script_id,omitempty"`
	OnSessionEndScriptId *int64    `json:"on_session_end_script_id,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
