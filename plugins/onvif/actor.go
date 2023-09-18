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

package onvif

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/encryptor"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	eventBus      bus.Bus
	client        *Client
}

// NewActor ...
func NewActor(entity *m.Entity,
	visor supervisor.Supervisor,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	eventBus bus.Bus) (actor *Actor) {

	actor = &Actor{
		BaseActor:     supervisor.NewBaseActor(entity, scriptService, adaptors),
		adaptors:      adaptors,
		scriptService: scriptService,
		eventBus:      eventBus,
	}

	actor.client = NewClient(actor.eventHandler)

	actor.Supervisor = visor

	clientBind := NewClientBind(actor.client)

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			// bind
			a.ScriptEngine.PushStruct("Actor", supervisor.NewScriptBind(actor))
			a.ScriptEngine.PushStruct("Camera", clientBind)
			_, _ = a.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
			_, _ = a.ScriptEngine.Do()
		}
	}

	if actor.ScriptEngine != nil {
		_, _ = actor.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
		actor.ScriptEngine.PushStruct("Actor", supervisor.NewScriptBind(actor))
		actor.ScriptEngine.PushStruct("Camera", clientBind)
	}

	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	if actor.Actions == nil {
		actor.Actions = NewActions()
	}

	return actor
}

func (e *Actor) destroy() {
	e.client.Shutdown()
}

// Spawn ...
func (e *Actor) Spawn() supervisor.PluginActor {
	e.client.Start(e.Setts[AttrUserName].String(), e.Setts[AttrPassword].Decrypt(), e.Setts[AttrAddress].String(), e.Setts[AttrOnvifPort].Int64())
	return e
}

// SetState ...
func (e *Actor) SetState(params supervisor.EntityStateParams) error {

	oldState := e.GetEventState(e)

	e.Now(oldState)

	if params.NewState != nil {
		state := e.States[*params.NewState]
		e.State = &state
		e.State.ImageUrl = state.ImageUrl
	}

	e.AttrMu.Lock()
	_, _ = e.Attrs.Deserialize(params.AttributeValues)
	e.AttrMu.Unlock()

	e.eventBus.Publish("system/entities/"+e.Id.String(), events.EventStateChanged{
		StorageSave: params.StorageSave,
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
	})

	return nil
}

func (e *Actor) addAction(event events.EventCallEntityAction) {
	e.runAction(event)
}

func (e *Actor) runAction(msg events.EventCallEntityAction) {
	action, ok := e.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if action.ScriptEngine == nil {
		return
	}
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId, action.Name, msg.Args); err != nil {
		log.Error(err.Error())
	}
}

func (e *Actor) eventHandler(msg interface{}) {
	switch v := msg.(type) {
	case *StreamList:
		go e.prepareStreamList(v)
	case *ConnectionStatus:
		go e.updateState(v)
	case *MotionAlarm:
		go e.prepareMotionAlarm(v)
	}
}

func (e *Actor) updateState(event *ConnectionStatus) {
	info := e.Info()
	var newStat = AttrOffline
	if event.Connected {
		newStat = AttrConnected
	}
	if info.State != nil && info.State.Name == newStat {
		return
	}
	e.SetState(supervisor.EntityStateParams{
		NewState:    common.String(newStat),
		StorageSave: true,
	})
}

func (e *Actor) prepareMotionAlarm(event *MotionAlarm) {
	e.SetState(supervisor.EntityStateParams{
		NewState: common.String(AttrConnected),
		AttributeValues: m.AttributeValue{
			AttrMotion:     event.State,
			AttrMotionTime: event.Time,
		},
		StorageSave: true,
	})
}

func (e *Actor) prepareStreamList(event *StreamList) {
	var attrs = m.AttributeValue{}
	if len(event.List) > 0 {
		attrs[AttrStreamUri], _ = encryptor.Encrypt(event.List[0])
	}
	if event.SnapshotUri != nil {
		attrs[AttrSnapshotUri], _ = encryptor.Encrypt(common.StringValue(event.SnapshotUri))
	}
	e.SetState(supervisor.EntityStateParams{
		NewState:        common.String(AttrConnected),
		AttributeValues: attrs,
		StorageSave:     true,
	})
}
