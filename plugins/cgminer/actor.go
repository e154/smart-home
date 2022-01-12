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

package cgminer

import (
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/pkg/errors"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/cgminer/bitmine"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	eventBus   event_bus.EventBus
	miner      IMiner
	actionPool chan event_bus.EventCallAction
}

// NewActor ...
func NewActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	eventBus event_bus.EventBus) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor:  entity_manager.NewBaseActor(entity, scriptService, adaptors),
		eventBus:   eventBus,
		actionPool: make(chan event_bus.EventCallAction, 10),
	}

	if actor.ParentId == nil {
		log.Warnf("entity %s, parent is nil", actor.Id)
	}

	actor.Manager = entityManager

	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	if actor.Actions == nil {
		actor.Actions = NewActions()
	}

	actor.DeserializeAttr(entity.Attributes.Serialize())

	if actor.Setts == nil || actor.Setts[SettingManufacturer] == nil {
		err = errors.Wrap(common.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	switch actor.Setts[SettingManufacturer].Value {
	case bitmine.ManufactureBitmine:
		switch actor.Setts[SettingModel].Value {
		case bitmine.DeviceS9:
		case bitmine.DeviceS7:
		case bitmine.DeviceL3:
		case bitmine.DeviceL3Plus:
		case bitmine.DeviceD3:
		case bitmine.DeviceT9:
		default:
			err = fmt.Errorf("unknown model %s", actor.Setts[SettingModel].Value)
		}
	default:
		err = fmt.Errorf("unknown manufacture %s", actor.Setts[SettingManufacturer].Value)
	}

	if _, ok := actor.Setts[SettingHost]; !ok {
		err = errors.Wrap(common.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingPort]; !ok {
		err = errors.Wrap(common.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingTimeout]; !ok {
		err = errors.Wrap(common.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingModel]; !ok {
		err = errors.Wrap(common.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingHost]; !ok {
		err = errors.Wrap(common.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingUser]; !ok {
		err = errors.Wrap(common.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingPass]; !ok {
		err = errors.Wrap(common.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	host := actor.Setts[SettingHost].String()
	port := int(actor.Setts[SettingPort].Int64())
	timeout := actor.Setts[SettingTimeout].Int64()
	model := actor.Setts[SettingModel].String()
	user := actor.Setts[SettingUser].String()
	pass := actor.Setts[SettingPass].String()
	actor.miner, err = bitmine.NewBitmine(bitmine.NewTransport(host, port, timeout), model, user, pass)

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			// bind
			a.ScriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
			a.ScriptEngine.PushFunction("Miner", actor.miner.Bind())
			a.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
			a.ScriptEngine.Do()
		}
	}

	if actor.ScriptEngine != nil {
		actor.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
		actor.ScriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
	}

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return
}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {
	e.Update()
	return e
}

// SetState ...
func (e *Actor) SetState(params entity_manager.EntityStateParams) error {

	oldState := e.GetEventState(e)

	e.Now(oldState)

	if params.NewState != nil {
		state := e.States[*params.NewState]
		e.State = &state
		e.State.ImageUrl = state.ImageUrl
	}

	e.AttrMu.Lock()
	e.Attrs.Deserialize(params.AttributeValues)
	e.AttrMu.Unlock()

	e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventStateChanged{
		StorageSave: params.StorageSave,
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
	})

	return nil
}

// Update ...
func (e *Actor) Update() {

}

func (e *Actor) addAction(event event_bus.EventCallAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg event_bus.EventCallAction) {
	action, ok := e.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId, action.Name); err != nil {
		log.Error(err.Error())
	}
}
