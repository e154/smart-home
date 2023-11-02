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

	. "github.com/e154/smart-home/common"
)

// Script ...
type Script struct {
	Id          int64          `json:"id"`
	Lang        ScriptLang     `json:"lang" validate:"required"`
	Name        string         `json:"name" validate:"max=254,required"`
	Source      string         `json:"source"`
	Description string         `json:"description"`
	Compiled    string         `json:"-"`
	Versions    ScriptVersions `json:"versions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Info        *ScriptInfo    `json:"info"`
}

type ScriptInfo struct {
	AlexaIntents         int `json:"alexa_intents"`
	EntityActions        int `json:"entity_actions"`
	EntityScripts        int `json:"entity_scripts"`
	AutomationTriggers   int `json:"automation_triggers"`
	AutomationConditions int `json:"automation_conditions"`
	AutomationActions    int `json:"automation_actions"`
}

type ScriptsStatistic struct {
	Total        int32
	Used         int32
	Unused       int32
	CoffeeScript int32
	TypeScript   int32
	JavaScript   int32
}

type ScriptVersions []*ScriptVersion

func (l ScriptVersions) Len() int      { return len(l) }
func (l ScriptVersions) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l ScriptVersions) Less(i, j int) bool {
	return l[i].CreatedAt.UnixNano() > l[j].CreatedAt.UnixNano()
}

type ScriptVersion struct {
	Id        int64      `json:"id"`
	Lang      ScriptLang `json:"lang" validate:"required"`
	Source    string     `json:"source"`
	CreatedAt time.Time  `json:"created_at"`
}
