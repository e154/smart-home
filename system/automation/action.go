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
	"sync"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/system/bus"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
)

// Action ...
type Action struct {
	model         *m.Action
	scriptService scripts.ScriptService
	eventBus      bus.Bus
	scriptEngine  *scripts.Engine
	inProcess     atomic.Bool
	sync.Mutex
}

// NewAction ...
func NewAction(scriptService scripts.ScriptService,
	eventBus bus.Bus,
	model *m.Action) (action *Action, err error) {

	action = &Action{
		scriptService: scriptService,
		eventBus:      eventBus,
		model:         model,
	}

	if model.Script != nil {
		if action.scriptEngine, err = scriptService.NewEngine(model.Script); err != nil {
			return
		}

		if _, err = action.scriptEngine.Do(); err != nil {
			return
		}

		action.scriptEngine.PushStruct("Action", NewActionBind(action))
	}

	eventBus.Subscribe(fmt.Sprintf("system/automation/actions/%d", model.Id), action.actionHandler)

	return
}

func (a *Action) Remove() {
	a.eventBus.Unsubscribe(fmt.Sprintf("system/automation/actions/%d", a.model.Id), a.actionHandler)
}

// Run ...
func (a *Action) Run(entityId *common.EntityId) (result string, err error) {
	a.Lock()
	defer a.Unlock()

	//log.Infof("run action")

	if a.scriptEngine != nil {
		if result, err = a.scriptEngine.AssertFunction(ActionFunc, entityId); err != nil {
			log.Error(err.Error())
		}
	}

	if a.model.EntityId != nil && a.model.EntityActionName != nil {
		id := *a.model.EntityId
		action := *a.model.EntityActionName
		a.eventBus.Publish("system/entities/"+id.String(), events.EventCallEntityAction{
			PluginName: id.PluginName(),
			EntityId:   id,
			ActionName: action,
		})
	}

	a.eventBus.Publish(fmt.Sprintf("system/automation/actions/%d", a.model.Id), events.EventActionCompleted{
		Id: a.model.Id,
	})

	return
}

func (a *Action) actionHandler(_ string, msg interface{}) {
	switch msg.(type) {
	case events.EventCallAction:
		a.Run(nil)
	}
}
