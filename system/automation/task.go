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
	"fmt"
	"github.com/e154/smart-home/common"

	"github.com/e154/smart-home/common/apperr"

	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/scripts"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

// Task ...
type Task struct {
	model          *m.Task
	automation     *automation
	conditionGroup *ConditionGroup
	actions        map[string]*Action
	script         *scripts.Engine
	enabled        atomic.Bool
	scriptService  scripts.ScriptService
	entityManager  entity_manager.EntityManager
	triggers       map[string]*Trigger
	rawPlugin      triggers.IGetTrigger
}

// NewTask ...
func NewTask(automation *automation,
	scriptService scripts.ScriptService,
	model *m.Task,
	entityManager entity_manager.EntityManager,
	rawPlugin triggers.IGetTrigger) *Task {
	return &Task{
		model:         model,
		automation:    automation,
		enabled:       atomic.Bool{},
		scriptService: scriptService,
		entityManager: entityManager,
		triggers:      make(map[string]*Trigger),
		actions:       make(map[string]*Action),
		rawPlugin:     rawPlugin,
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

func (t *Task) addTrigger(model *m.Trigger) (err error) {

	if _, ok := t.triggers[model.Name]; ok {
		err = errors.Wrap(apperr.ErrInternal, fmt.Sprintf("trigger %s exist", model.Name))
		return
	}

	var tr *Trigger
	if tr, err = NewTrigger(t.scriptService, t.model.Name, model, t.rawPlugin, t.triggerHandler); err != nil {
		log.Error(err.Error())
		return
	}

	t.triggers[model.Name] = tr

	_ = tr.Start()

	return
}

func (t *Task) triggerHandler(entityId *common.EntityId) {
	result, err := t.conditionGroup.Check(entityId)
	if err != nil || !result {
		return
	}

	for _, acion := range t.actions {
		if _, err = acion.Run(entityId); err != nil {
			log.Error(err.Error())
		}
	}
}

// Start ...
func (t *Task) Start() {
	if t.enabled.Load() {
		return
	}
	t.enabled.Store(true)

	log.Infof("task %d start", t.Id())

	// add actions
	for _, model := range t.model.Actions {
		action, err := NewAction(t.scriptService, t.entityManager, model)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		t.actions[model.Name] = action
	}

	// add condition group
	t.conditionGroup = NewConditionGroup(t.automation, t.model.Condition)

	// add conditions
	for _, model := range t.model.Conditions {
		condition, err := NewCondition(t.scriptService, model, t.entityManager)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		t.conditionGroup.AddCondition(condition)
	}

	// add triggers
	for _, model := range t.model.Triggers {
		if err := t.addTrigger(model); err != nil {
			log.Error(err.Error())
		}
	}
}

// Stop ...
func (t *Task) Stop() {
	if !t.enabled.Load() {
		return
	}
	t.enabled.Store(false)
	log.Infof("task %d stopped", t.Id())
	for _, trigger := range t.triggers {
		_ = trigger.Stop()
	}
	t.triggers = make(map[string]*Trigger)
	t.actions = make(map[string]*Action)
	t.conditionGroup = nil
}

// CallTrigger ...
func (t *Task) CallTrigger(name string) {
	if tr, ok := t.triggers[name]; ok {
		tr.Call()
	}
}

// CallAction ...
func (t *Task) CallAction(name string) {
	if action, ok := t.actions[name]; ok {
		_, _ = action.Run(nil)
	}
}
