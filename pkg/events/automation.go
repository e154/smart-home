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

	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
)

// Tasks
// ------------------------------------------------

type EventTaskCompleted struct {
	Id  int64           `json:"id"`
	Ctx context.Context `json:"ctx"`
}

// CommandEnableTask ...
type CommandEnableTask struct {
	Id int64 `json:"id"`
}

// CommandDisableTask ...
type CommandDisableTask struct {
	Id int64 `json:"id"`
}

// EventCreatedTaskModel ...
type EventCreatedTaskModel struct {
	Id int64 `json:"id"`
}

// EventRemovedTaskModel ...
type EventRemovedTaskModel struct {
	Id int64 `json:"id"`
}

// EventUpdatedTaskModel ...
type EventUpdatedTaskModel struct {
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

type TriggerMessage struct {
	Payload     interface{}      `json:"payload"`
	TriggerName string           `json:"trigger_name"`
	EntityId    *common.EntityId `json:"entity_id"`
}

type EventTriggerCompleted struct {
	Id       int64            `json:"id"`
	Args     *TriggerMessage  `json:"args"`
	EntityId *common.EntityId `json:"entity_id,omitempty"`
	Ctx      context.Context  `json:"ctx"`
}

// CommandEnableTrigger ...
type CommandEnableTrigger struct {
	Id int64 `json:"id"`
}

// CommandDisableTrigger ...
type CommandDisableTrigger struct {
	Id int64 `json:"id"`
}

// EventCallTrigger ...
type EventCallTrigger struct {
	Id  int64           `json:"id"`
	Ctx context.Context `json:"ctx"`
}

// EventCreatedTriggerModel ...
type EventCreatedTriggerModel struct {
	Id int64 `json:"id"`
}

// EventRemovedTriggerModel ...
type EventRemovedTriggerModel struct {
	Id int64 `json:"id"`
}

// EventUpdatedTriggerModel ...
type EventUpdatedTriggerModel struct {
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

// EventAddedActionModel ...
type EventAddedActionModel struct {
	Id int64 `json:"id"`
}

// EventRemovedActionModel ...
type EventRemovedActionModel struct {
	Id int64 `json:"id"`
}

// EventUpdatedActionModel ...
type EventUpdatedActionModel struct {
	Id     int64     `json:"id"`
	Action *m.Action `json:"action"`
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

// EventAddedConditionModel ...
type EventAddedConditionModel struct {
	Id int64 `json:"id"`
}

// EventRemovedConditionModel ...
type EventRemovedConditionModel struct {
	Id int64 `json:"id"`
}

// EventUpdatedConditionModel ...
type EventUpdatedConditionModel struct {
	Id int64 `json:"id"`
}
