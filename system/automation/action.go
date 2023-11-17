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

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/system/bus"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
)

// ActionFunc ...
const ActionFunc = "automationAction"

// Action ...
type Action struct {
	model         *m.Action
	scriptService scripts.ScriptService
	eventBus      bus.Bus
	scriptEngine  *scripts.EngineWatcher
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
		if action.scriptEngine, err = scriptService.NewEngineWatcher(model.Script); err != nil {
			return
		}

		action.scriptEngine.Spawn(func(engine *scripts.Engine) {
			if _, err = engine.Do(); err != nil {
				return
			}

			engine.PushStruct("Action", NewActionBind(action))

			if model.EntityId != nil {
				if _, err = engine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", model.EntityId.String())); err != nil {
					log.Error(err.Error())
				}
			}
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
