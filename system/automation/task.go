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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
)

const (
	taskMsgBuffer = 50
)

// Task ...
type Task struct {
	model          *m.Task
	automation     *automation
	conditionGroup *ConditionGroup
	actions        []*Action
	script         *scripts.Engine
	enabled        atomic.Bool
	scriptService  scripts.ScriptService
	entityManager  entity_manager.EntityManager
	triggers       map[*Trigger]struct{}
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
		triggers:      make(map[*Trigger]struct{}),
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

	pluginName := model.PluginName
	if pluginName == "" {
		pluginName = triggers.StateChangeName
	}

	var triggerPLugin triggers.ITrigger
	if triggerPLugin, err = t.rawPlugin.GetTrigger(pluginName); err != nil {
		log.Error(err.Error())
		return
	}

	var tr *Trigger
	if tr, err = NewTrigger(t.scriptService, model, triggerPLugin); err != nil {
		log.Error(err.Error())
		return
	}

	tr.Start()

	queue := make(chan interface{}, taskMsgBuffer)
	var handler = func(_ string, msg interface{}) {
		queue <- map[string]interface{}{
			"payload":      msg,
			"trigger_name": model.Name,
			"task_name":    t.model.Name,
			"entity_id":    model.EntityId.String(),
		}
	}

	if err = triggerPLugin.Subscribe(triggers.Subscriber{
		EntityId: tr.EntityId(),
		Handler:  handler,
		Payload:  tr.model.Payload,
	}); err != nil {
		log.Error(err.Error())
		return
	}

	go func(entityId *common.EntityId) {
		var result bool
		var err error

		defer func() {
			triggerPLugin.Unsubscribe(triggers.Subscriber{
				EntityId: entityId,
				Handler:  handler,
				Payload:  tr.model.Payload,
			})
			close(queue)
		}()

		for msg := range queue {

			result, err = tr.Check(msg)
			if err != nil || !result {
				continue
			}

			result, err = t.conditionGroup.Check(entityId)
			if err != nil || !result {
				continue
			}

			for _, acion := range t.actions {
				if _, err = acion.Run(entityId); err != nil {
					log.Error(err.Error())
				}
			}
		}
	}(tr.EntityId())

	return
}

// Start ...
func (t *Task) Start() {
	if t.enabled.Load() {
		return
	}
	t.enabled.Store(true)

	// add actions
	for _, model := range t.model.Actions {
		action, err := NewAction(t.scriptService, t.entityManager, model)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		t.actions = append(t.actions, action)
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
	for trigger := range t.triggers {
		trigger.Stop()
	}
	t.triggers = make(map[*Trigger]struct{})
	t.conditionGroup = nil
	t.actions = make([]*Action, 0)
}
