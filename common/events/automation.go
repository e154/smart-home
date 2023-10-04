// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package events

import (
	"context"

	"github.com/e154/smart-home/common"
)

// Tasks
// ------------------------------------------------

type EventTaskCompleted struct {
	Id  int64           `json:"id"`
	Ctx context.Context `json:"ctx"`
}

// EventEnableTask ...
type EventEnableTask struct {
	Id int64 `json:"id"`
}

// EventDisableTask ...
type EventDisableTask struct {
	Id int64 `json:"id"`
}

// EventAddedTask ...
type EventAddedTask struct {
	Id int64 `json:"id"`
}

// EventRemoveTask ...
type EventRemoveTask struct {
	Id int64 `json:"id"`
}

// EventUpdateTask ...
type EventUpdateTask struct {
	Id int64 `json:"id"`
}

// EventTaskLoaded ...
type EventTaskLoaded struct {
	Id int64 `json:"id"`
}

// EventTaskUnloaded ...
type EventTaskUnloaded struct {
	Id int64 `json:"id"`
}

// Triggers
// ------------------------------------------------

type EventTriggerCompleted struct {
	Id       int64                  `json:"id"`
	Args     map[string]interface{} `json:"args"`
	EntityId *common.EntityId       `json:"entity_id,omitempty"`
	Ctx      context.Context        `json:"ctx"`
}

// EventEnableTrigger ...
type EventEnableTrigger struct {
	Id int64 `json:"id"`
}

// EventDisableTrigger ...
type EventDisableTrigger struct {
	Id int64 `json:"id"`
}

// EventCallTrigger ...
type EventCallTrigger struct {
	Id  int64           `json:"id"`
	Ctx context.Context `json:"ctx"`
}

// EventAddedTrigger ...
type EventAddedTrigger struct {
	Id int64 `json:"id"`
}

// EventRemovedTrigger ...
type EventRemovedTrigger struct {
	Id int64 `json:"id"`
}

// EventUpdatedTrigger ...
type EventUpdatedTrigger struct {
	Id int64 `json:"id"`
}

// EventTriggerLoaded ...
type EventTriggerLoaded struct {
	Id int64 `json:"id"`
}

// EventTriggerUnloaded ...
type EventTriggerUnloaded struct {
	Id int64 `json:"id"`
}

// Actions
// ------------------------------------------------

type EventActionCompleted struct {
	Id  int64           `json:"id"`
	Ctx context.Context `json:"ctx"`
}

// EventCallTaskAction ...
type EventCallTaskAction struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// EventCallAction ...
type EventCallAction struct {
	Id  int64           `json:"id"`
	Ctx context.Context `json:"ctx"`
}

// EventAddedAction ...
type EventAddedAction struct {
	Id int64 `json:"id"`
}

// EventRemovedAction ...
type EventRemovedAction struct {
	Id int64 `json:"id"`
}

// EventUpdatedAction ...
type EventUpdatedAction struct {
	Id int64 `json:"id"`
}

// EventActionLoaded ...
type EventActionLoaded struct {
	Id int64 `json:"id"`
}

// EventActionUnloaded ...
type EventActionUnloaded struct {
	Id int64 `json:"id"`
}

// Conditions
// ------------------------------------------------

// EventAddedCondition ...
type EventAddedCondition struct {
	Id int64 `json:"id"`
}

// EventRemovedCondition ...
type EventRemovedCondition struct {
	Id int64 `json:"id"`
}

// EventUpdatedCondition ...
type EventUpdatedCondition struct {
	Id int64 `json:"id"`
}
