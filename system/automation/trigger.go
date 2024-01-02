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
	"go.uber.org/atomic"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/telemetry"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
)

// Trigger ...
type Trigger struct {
	scriptService      scripts.ScriptService
	lastStatus         *atomic.Bool
	model              *m.Trigger
	name               string
	triggerPlugin      triggers.ITrigger
	taskName           string
	triggerSubscribers []*TriggerSubscriber
	eventBus           bus.Bus
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
		model:              model,
		name:               model.Name,
		scriptService:      scriptService,
		triggerPlugin:      triggerPlugin,
		eventBus:           eventBus,
		lastStatus:         atomic.NewBool(false),
		triggerSubscribers: make([]*TriggerSubscriber, 0),
	}

	return
}

// Start ...
func (tr *Trigger) Start() {
	log.Infof("start trigger '%s'", tr.name)

	if len(tr.model.Entities) > 0 {
		for _, entity := range tr.model.Entities {
			engine := tr.genScriptEngine(&entity.Id)
			tr.triggerSubscribers = append(tr.triggerSubscribers, &TriggerSubscriber{
				Engine:     engine,
				Subscriber: tr.genSubscriber(&entity.Id, tr.genCheckFunc(engine)),
			})
		}

	} else {
		engine := tr.genScriptEngine(nil)
		tr.triggerSubscribers = append(tr.triggerSubscribers, &TriggerSubscriber{
			Engine:     engine,
			Subscriber: tr.genSubscriber(nil, tr.genCheckFunc(engine)),
		})
	}

	for _, sub := range tr.triggerSubscribers {
		_ = tr.triggerPlugin.Subscribe(sub.Subscriber)
	}
	_ = tr.eventBus.Subscribe(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), tr.eventHandler, false)
	tr.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), events.EventTriggerLoaded{
		Id: tr.model.Id,
	})
}

// Stop ...
func (tr *Trigger) Stop() {
	log.Infof("stop trigger '%s'", tr.name)
	for _, sub := range tr.triggerSubscribers {
		if sub.Engine != nil {
			sub.Engine.Stop()
		}
	}
	_ = tr.eventBus.Unsubscribe(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), tr.eventHandler)
	for _, sub := range tr.triggerSubscribers {
		_ = tr.triggerPlugin.Unsubscribe(sub.Subscriber)
	}
	tr.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), events.EventTriggerUnloaded{
		Id: tr.model.Id,
	})
}

func (tr *Trigger) genCheckFunc(scriptEngine *scripts.EngineWatcher) func(msg interface{}) (state bool, err error) {
	return func(msg interface{}) (state bool, err error) {

		if scriptEngine != nil {
			var result string
			if result, err = scriptEngine.Engine().AssertFunction(tr.triggerPlugin.FunctionName(), msg); err != nil {
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
}

func (tr *Trigger) genSubscriber(entityId *common.EntityId, check func(msg interface{}) (state bool, err error)) triggers.Subscriber {

	return triggers.Subscriber{
		EntityId: entityId,
		Payload:  tr.model.Payload,
		Handler: func(_ string, msg interface{}) {
			triggerCtx, span := telemetry.Start(context.Background(), "trigger")
			span.SetAttributes("id", tr.model.Id)
			b, _ := json.Marshal(msg)
			args := map[string]interface{}{
				"payload":      string(b),
				"trigger_name": tr.model.Name,
				"entity_id":    entityId,
			}
			result, err := check(args)
			span.End()
			if err != nil || !result {
				return
			}
			//fmt.Println("call trigger", tr.model.Name, tr.triggerPlugin.Name())
			tr.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), events.EventTriggerCompleted{
				Id:       tr.model.Id,
				Args:     args,
				EntityId: entityId,
				Ctx:      triggerCtx,
			})
		},
	}
}

func (tr *Trigger) genScriptEngine(entityId *common.EntityId) (scriptEngine *scripts.EngineWatcher) {
	if tr.model.Script != nil {

		var err error
		if scriptEngine, err = tr.scriptService.NewEngineWatcher(tr.model.Script); err != nil {
			return
		}
		scriptEngine.Spawn(func(engine *scripts.Engine) {

			engine.PushStruct("Trigger", NewTriggerBind(tr))

			if entityId != nil {
				if _, err = engine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entityId.String())); err != nil {
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

func (tr *Trigger) eventHandler(_ string, msg interface{}) {

	var entityId *common.EntityId

	if len(tr.model.Entities) == 1 {
		entityId = &tr.model.Entities[0].Id
	}

	switch v := msg.(type) {
	case events.EventCallTrigger:
		triggerCtx, span := telemetry.Start(v.Ctx, "trigger")
		span.SetAttributes("id", tr.model.Id)
		span.End()
		tr.eventBus.Publish(fmt.Sprintf("system/automation/triggers/%d", tr.model.Id), events.EventTriggerCompleted{
			Id:       tr.model.Id,
			Args:     nil,
			EntityId: entityId,
			Ctx:      triggerCtx,
		})
	}
}
