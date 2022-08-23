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
	"encoding/json"
	"sync"

	"go.uber.org/atomic"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/scripts"
)

// Trigger ...
type Trigger struct {
	scriptEngine  *scripts.Engine
	scriptService scripts.ScriptService
	lastStatus    atomic.Bool
	model         *m.Trigger
	name          string
	triggerPlugin triggers.ITrigger
	isStarted     atomic.Bool
	cb            func(entityId *common.EntityId)
	taskName      string
	subscriber    triggers.Subscriber
	sync.Mutex
}

// NewTrigger ...
func NewTrigger(scriptService scripts.ScriptService,
	taskName string,
	model *m.Trigger,
	rawPlugin triggers.IGetTrigger,
	cb func(entityId *common.EntityId)) (tr *Trigger, err error) {

	pluginName := model.PluginName
	if pluginName == "" {
		pluginName = triggers.StateChangeName
	}

	var triggerPlugin triggers.ITrigger
	if triggerPlugin, err = rawPlugin.GetTrigger(pluginName); err != nil {
		log.Error(err.Error())
		return
	}

	tr = &Trigger{
		model:         model,
		name:          model.Name,
		scriptService: scriptService,
		triggerPlugin: triggerPlugin,
		cb:            cb,
		taskName:      taskName,
	}

	tr.subscriber = triggers.Subscriber{
		EntityId: model.EntityId,
		Payload:  model.Payload,
		Handler: func(_ string, msg interface{}) {
			b, _ := json.Marshal(msg)
			obj := &Obj{
				Payload:     string(b),
				TriggerName: tr.model.Name,
				TaskName:    tr.taskName,
				EntityId:    tr.model.EntityId.String(),
			}
			result, err := tr.Check(obj)
			if err != nil || !result {
				return
			}
			tr.cb(tr.EntityId())
			return
		},
	}

	if tr.scriptEngine, err = scriptService.NewEngine(model.Script); err != nil {
		return
	}

	tr.scriptEngine.PushStruct("Trigger", NewTriggerBind(tr))

	if _, err = tr.scriptEngine.Do(); err != nil {
		return
	}

	return
}

// Check ...
func (tr *Trigger) Check(msg interface{}) (state bool, err error) {
	tr.Lock()
	defer tr.Unlock()

	var result string
	if result, err = tr.scriptEngine.AssertFunction(tr.triggerPlugin.FunctionName(), msg); err != nil {
		log.Error(err.Error())
	}

	state = result == "true"
	tr.lastStatus.Store(state)

	return
}

// Id ...
func (tr *Trigger) Id() int64 {
	return tr.model.Id
}

// EntityId ...
func (tr *Trigger) EntityId() *common.EntityId {
	return tr.model.EntityId
}

// Start ...
func (tr *Trigger) Start() (err error) {
	log.Infof("start trigger '%s'", tr.name)
	err = tr.triggerPlugin.Subscribe(tr.subscriber)
	return
}

// Stop ...
func (tr *Trigger) Stop() (err error) {
	log.Infof("stop trigger '%s'", tr.name)
	err = tr.triggerPlugin.Unsubscribe(tr.subscriber)
	return
}

// Call ...
func (tr *Trigger) Call() {
	tr.triggerPlugin.CallManual()
}

type Obj struct {
	Payload     string `json:"payload"`
	TriggerName string `json:"trigger_name"`
	TaskName    string `json:"task_name"`
	EntityId    string `json:"entity_id"`
}
