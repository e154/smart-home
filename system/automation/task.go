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

package automation

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/atomic"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/telemetry"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
)

// Task ...
type Task struct {
	model          *m.Task
	eventBus       bus.Bus
	conditionGroup *ConditionGroup
	actions        map[int64]*Action
	script         *scripts.Engine
	enabled        atomic.Bool
	scriptService  scripts.ScriptService
	sync.Mutex
	telemetry telemetry.Telemetry
}

// NewTask ...
func NewTask(eventBus bus.Bus,
	scriptService scripts.ScriptService,
	model *m.Task) *Task {
	return &Task{
		model:         model,
		eventBus:      eventBus,
		enabled:       atomic.Bool{},
		scriptService: scriptService,
		actions:       make(map[int64]*Action),
	}
}

// Id ...
func (t *Task) Id() int64 {
	return t.model.Id
}

// Name ...
func (t *Task) Name() string {
	return t.model.Name
}

func (t *Task) triggerHandler(_ string, msg interface{}) {

	if t == nil || !t.enabled.Load() {
		return
	}

	switch v := msg.(type) {
	case events.EventTriggerCompleted:
		taskCtx, taskSpan := telemetry.Start(v.Ctx, "task")
		taskSpan.SetAttributes("id", t.model.Id)

		var actionCtx context.Context
		defer func() {
			taskSpan.End()

			t.Lock()
			t.telemetry = telemetry.Unpack(actionCtx)
			t.Unlock()
		}()

		conditionsCtx, span := telemetry.Start(taskCtx, "conditions")
		result, err := t.conditionGroup.Check(v.EntityId)
		if err != nil || !result {
			if err != nil {
				span.SetStatus(telemetry.Error, err.Error())
			}
			span.End()
			return
		}
		span.End()

		for _, action := range t.actions {
			if actionCtx != nil {
				conditionsCtx = actionCtx
			}
			actionCtx, span = telemetry.Start(conditionsCtx, "actions")
			span.SetAttributes("id", action.model.Id)
			if _, err = action.Run(v.EntityId); err != nil {
				log.Error(err.Error())
				span.SetStatus(telemetry.Error, err.Error())
			}
			span.End()
		}

		//fmt.Println("time spent", timeSpent.Microseconds())
		t.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", t.model.Id), events.EventTaskCompleted{
			Id:  t.model.Id,
			Ctx: actionCtx,
		})
	}
}

// Start ...
func (t *Task) Start() {
	if t.enabled.Load() {
		return
	}
	t.enabled.Store(true)

	//log.Infof("task %d start", t.Id())

	// add actions
	for _, model := range t.model.Actions {
		action, err := NewAction(t.scriptService, t.eventBus, model)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		t.actions[model.Id] = action
	}

	// add condition group
	t.conditionGroup = NewConditionGroup(t.model.Condition)

	// add conditions
	for _, model := range t.model.Conditions {
		condition, err := NewCondition(t.scriptService, model)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		t.conditionGroup.AddCondition(condition)
	}

	// add triggers
	for _, model := range t.model.Triggers {
		t.eventBus.Subscribe(fmt.Sprintf("system/automation/triggers/%d", model.Id), t.triggerHandler)
	}

	t.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", t.model.Id), events.EventTaskLoaded{
		Id: t.model.Id,
	})
}

// Stop ...
func (t *Task) Stop() {
	if !t.enabled.Load() {
		return
	}
	t.enabled.Store(false)
	//log.Infof("task %d stopped", t.Id())
	for _, model := range t.model.Triggers {
		t.eventBus.Unsubscribe(fmt.Sprintf("system/automation/triggers/%d", model.Id), t.triggerHandler)
	}
	for _, action := range t.actions {
		action.Remove()
	}
	t.actions = make(map[int64]*Action)
	t.conditionGroup = nil

	t.eventBus.Publish(fmt.Sprintf("system/automation/tasks/%d", t.model.Id), events.EventTaskUnloaded{
		Id: t.model.Id,
	})
}

func (t *Task) Telemetry() telemetry.Telemetry {
	t.Lock()
	defer t.Unlock()
	return t.telemetry
}
