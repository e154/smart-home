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

package core

import (
	. "github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	"sync"
)

type Action struct {
	Device        *m.Device
	Node          *Node
	ScriptEngine  *scripts.Engine
	flow          *Flow
	scriptService *scripts.ScriptService
	script        *m.Script
	doLock        sync.Mutex
}

func NewAction(device *m.Device,
	script *m.Script,
	node *Node,
	flow *Flow,
	scriptService *scripts.ScriptService) (action *Action, err error) {

	action = &Action{
		Device:        device,
		Node:          node,
		flow:          flow,
		scriptService: scriptService,
		script:        script,
	}

	err = action.NewScript()

	return
}

func (a *Action) Do() (res string, err error) {
	a.doLock.Lock()
	err = a.ScriptEngine.EvalString(a.script.Compiled)
	a.doLock.Unlock()
	return
}

func (a *Action) NewScript() (err error) {

	if a.flow != nil {
		if a.ScriptEngine, err = a.flow.NewScript(); err != nil {
			return
		}

	} else {
		model := &m.Script{
			Lang: ScriptLangJavascript,
		}
		if a.ScriptEngine, err = a.scriptService.NewEngine(model); err != nil {
			return
		}

		a.ScriptEngine.PushGlobalProxy("message", NewMessage())
	}

	// bind device
	a.ScriptEngine.PushGlobalProxy("device", &DeviceBind{
		model: a.Device,
		node:  a.Node,
	})

	// bind
	javascript := a.ScriptEngine.Get().(*scripts.Javascript)
	ctx := javascript.Ctx()
	if b := ctx.GetGlobalString("IC"); !b {
		return
	}
	ctx.PushObject()
	ctx.PushGoFunction(func() *ActionBind {
		return &ActionBind{action: a}
	})
	ctx.PutPropString(-3, "Action")
	ctx.Pop()

	return nil
}

func (a *Action) GetDevice() *m.Device {
	return a.Device
}

func (a *Action) GetNode() *Node {
	return a.Node
}
