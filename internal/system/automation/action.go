// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/scripts"

	"github.com/e154/bus"
)

// ActionFunc ...
const ActionFunc = "automationAction"

// Action ...
type Action struct {
	model         *m.Action
	scriptService scripts.ScriptService
	eventBus      bus.Bus
	scriptEngine  scripts.EngineWatcher
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
		if action.scriptEngine, err = scriptService.NewEngineWatcher(model.Script); err != nil {
			return
		}

		action.scriptEngine.PushStruct("Action", NewActionBind(action))
		action.scriptEngine.BeforeSpawn(func(engine scripts.Engine) {
			if model.EntityId != nil {
				if _, err = engine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", model.EntityId.String())); err != nil {
					log.Error(err.Error())
				}
			}
		})
		action.scriptEngine.Spawn(func(engine scripts.Engine) {
			//if _, err = engine.Do(); err != nil {
			//	return
			//}
		})
	}

	_ = eventBus.Subscribe(fmt.Sprintf("system/automation/actions/%d", model.Id), action.actionHandler, false)

	return
}

func (a *Action) Remove() {
	if a.scriptEngine != nil {
		a.scriptEngine.Stop()
	}
	_ = a.eventBus.Unsubscribe(fmt.Sprintf("system/automation/actions/%d", a.model.Id), a.actionHandler)
}

// Run ...
func (a *Action) Run(entityId *common.EntityId) (result string, err error) {
	a.Lock()
	defer a.Unlock()

	//log.Infof("run action")

	if a.scriptEngine != nil {
		if result, err = a.scriptEngine.Engine().AssertFunction(ActionFunc, entityId); err != nil {
			log.Error(err.Error())
		}
	}

	if a.model.EntityId != nil && a.model.EntityActionName != nil {
		id := *a.model.EntityId
		action := *a.model.EntityActionName
		a.eventBus.Publish("system/entities/"+id.String(), events.EventCallEntityAction{
			PluginName: common.String(id.PluginName()),
			EntityId:   id.Ptr(),
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
