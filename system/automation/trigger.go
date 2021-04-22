// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/orgo
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
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
	"sync"
)

// Trigger ...
type Trigger struct {
	scriptEngine  *scripts.Engine
	scriptService scripts.ScriptService
	lastStatus    atomic.Bool
	model         *m.Trigger
	queue         chan interface{}
	triggerPlugin triggers.ITrigger
	isStarted     atomic.Bool
	sync.Mutex
}

// NewTrigger ...
func NewTrigger(scriptService scripts.ScriptService,
	model *m.Trigger,
	triggerPLugin triggers.ITrigger) (tr *Trigger, err error) {

	tr = &Trigger{
		model:         model,
		queue:         make(chan interface{}, taskMsgBuffer),
		scriptService: scriptService,
		triggerPlugin: triggerPLugin,
	}

	if tr.scriptEngine, err = scriptService.NewEngine(model.Script); err != nil {
		return
	}

	if _, err = tr.scriptEngine.Do(); err != nil {
		return
	}

	tr.scriptEngine.PushStruct("Trigger", NewTriggerBind(tr))

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

	return
}

// Stop ...
func (tr *Trigger) Stop() (err error) {

	return
}
