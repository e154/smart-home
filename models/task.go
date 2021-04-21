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

import (
	"github.com/e154/smart-home/common"
	"time"
)

// Task ...
type Task struct {
	Id          int64                `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Enabled     bool                 `json:"enabled"`
	Condition   common.ConditionType `json:"condition"`
	Triggers    []*Trigger           `json:"triggers"`
	Conditions  []*Condition         `json:"conditions"`
	Actions     []*Action            `json:"actions"`
	Area        *Area                `json:"area"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// AddTrigger ...
func (t *Task) AddTrigger(tr *Trigger) {
	t.Triggers = append(t.Triggers, tr)
}

// AddCondition ...
func (t *Task) AddCondition(c *Condition) {
	t.Conditions = append(t.Conditions, c)
}

// AddAction ...
func (t *Task) AddAction(a *Action) {
	t.Actions = append(t.Actions, a)
}
