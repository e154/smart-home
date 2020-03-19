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
	flow          *Flow
	scriptService *scripts.ScriptService
	ScriptEngine  *scripts.Engine
	deviceAction  *m.DeviceAction
	doLock        sync.Mutex
}

func NewAction(device *m.Device,
	deviceAction *m.DeviceAction,
	node *Node,
	flow *Flow,
	scriptService *scripts.ScriptService) (action *Action, err error) {

	action = &Action{
		Device:        device,
		Node:          node,
		flow:          flow,
		scriptService: scriptService,
		deviceAction:  deviceAction,
	}

	err = action.newScript()

	return
}

func (a *Action) Do() (res string, err error) {
	a.doLock.Lock()
	defer a.doLock.Unlock()
	if a.deviceAction.Script == nil {
		return
	}
	res, err = a.ScriptEngine.EvalScript(a.deviceAction.Script)
	return
}

func (a *Action) newScript() (err error) {

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

		a.ScriptEngine.PushStruct("message", NewMessage())
	}

	// bind device
	a.ScriptEngine.PushStruct("Device", &DeviceBind{
		model: a.Device,
		node:  a.Node,
	})

	a.ScriptEngine.PushStruct("Action", &ActionBind{
		Id:          a.deviceAction.Id,
		Name:        a.deviceAction.Name,
		Description: a.deviceAction.Description,
		action:      a,
	})

	return nil
}

func (a *Action) GetDevice() *m.Device {
	return a.Device
}

func (a *Action) GetNode() *Node {
	return a.Node
}
