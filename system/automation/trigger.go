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
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/e154/smart-home/common/telemetry"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
)

// Trigger ...
type Trigger struct {
	scriptEngine  *scripts.EngineWatcher
	scriptService scripts.ScriptService
	lastStatus    atomic.Bool
	model         *m.Trigger
	name          string
	triggerPlugin triggers.ITrigger
	isStarted     atomic.Bool
	taskName      string
	subscriber    triggers.Subscriber
	eventBus      bus.Bus
	sync.Mutex
}

// NewTrigger ...
func NewTrigger(
	eventBus bus.Bus,
	scriptService scripts.ScriptService,
	model *m.Trigger,
	rawPlugin triggers.IGetTrigger) (tr *Trigger, err error) {

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
		eventBus:      eventBus,
	}

	tr.subscriber = triggers.Subscriber{
		EntityId: model.EntityId,
		Payload:  model.Payload,
		Handler: func(_ string, msg interface{}) {
			triggerCtx, span := telemetry.Start(context.Background(), "trigger")
			span.SetAttributes("id", tr.model.Id)
			b, _ := json.Marshal(msg)
			args := map[string]interface{}{
				"payload":      string(b),
				"trigger_name": tr.model.Name,
				"entity_id":    tr.EntityId(),
			}
			result, err := tr.Check(args)
			if err != nil || !result {
				span.End()
				return
			}
			span.End()
			//fmt.Println("call trigger", tr.model.Name, tr.triggerPlugin.Name())
			eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), events.EventTriggerCompleted{
				Id:       model.Id,
				Args:     args,
				EntityId: tr.EntityId(),
				Ctx:      triggerCtx,
			})
		},
	}

	if model.Script != nil {

		if tr.scriptEngine, err = scriptService.NewEngineWatcher(model.Script); err != nil {
			return
		}
		tr.scriptEngine.Spawn(func(engine *scripts.Engine) {

			engine.PushStruct("Trigger", NewTriggerBind(tr))

			if model.EntityId != nil {
				if _, err = engine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", model.EntityId.String())); err != nil {
					log.Error(err.Error())
				}
			}

			if _, err = engine.Do(); err != nil {
				return
			}
		})

	}

	return
}

// Check ...
func (tr *Trigger) Check(msg interface{}) (state bool, err error) {
	tr.Lock()
	defer tr.Unlock()

	if tr.scriptEngine != nil {
		var result string
		if result, err = tr.scriptEngine.Engine().AssertFunction(tr.triggerPlugin.FunctionName(), msg); err != nil {
			log.Error(err.Error())
		}

		state = result == "true"
		tr.lastStatus.Store(state)
		return
	}

	state = true

	tr.lastStatus.Store(state)

	return
}

// Id ...
func (tr *Trigger) Id() int64 {
	return tr.model.Id
}

// EntityId ...
func (tr *Trigger) EntityId() *common.EntityId {
	if tr.model == nil {
		return nil
	}
	return tr.model.EntityId
}

// Start ...
func (tr *Trigger) Start() {
	log.Infof("start trigger '%s'", tr.name)
	_ = tr.triggerPlugin.Subscribe(tr.subscriber)
	_ = tr.eventBus.Subscribe(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), tr.eventHandler, false)
	tr.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), events.EventTriggerLoaded{
		Id: tr.model.Id,
	})
}

// Stop ...
func (tr *Trigger) Stop() {
	log.Infof("stop trigger '%s'", tr.name)
	if tr.scriptEngine != nil {
		tr.scriptEngine.Stop()
	}
	_ = tr.eventBus.Unsubscribe(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), tr.eventHandler)
	_ = tr.triggerPlugin.Unsubscribe(tr.subscriber)
	tr.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), events.EventTriggerUnloaded{
		Id: tr.model.Id,
	})
}

func (tr *Trigger) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventCallTrigger:
		triggerCtx, span := telemetry.Start(v.Ctx, "trigger")
		span.SetAttributes("id", tr.model.Id)
		span.End()
		tr.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), events.EventTriggerCompleted{
			Id:       tr.model.Id,
			Args:     nil,
			EntityId: tr.EntityId(),
			Ctx:      triggerCtx,
		})
	}
}
