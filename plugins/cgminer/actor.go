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

package cgminer

import (
	"fmt"

	"github.com/e154/smart-home/system/scripts"
	"github.com/pkg/errors"

	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/cgminer/bitmine"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	miner      IMiner
	actionPool chan events.EventCallEntityAction
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		actionPool: make(chan events.EventCallEntityAction, 1000),
	}

	//if actor.ParentId == nil {
	//	log.Warnf("entity %s, parent is nil", actor.Id)
	//}

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
		err = errors.Wrap(apperr.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
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
		err = errors.Wrap(apperr.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingPort]; !ok {
		err = errors.Wrap(apperr.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingTimeout]; !ok {
		err = errors.Wrap(apperr.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingModel]; !ok {
		err = errors.Wrap(apperr.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingHost]; !ok {
		err = errors.Wrap(apperr.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingUser]; !ok {
		err = errors.Wrap(apperr.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
		return
	}

	if _, ok := actor.Setts[SettingPass]; !ok {
		err = errors.Wrap(apperr.ErrBadSettings, fmt.Sprintf("actor 'cgminer', current settings %+v", actor.Setts))
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
		if a.ScriptEngine.Engine() != nil {
			// bind
			a.ScriptEngine.Engine().PushFunction("Miner", actor.miner.Bind())
			_, _ = a.ScriptEngine.Engine().Do()
		}
	}

	actor.ScriptsEngine.Spawn(func(engine *scripts.Engine) {
		engine.PushFunction("Miner", actor.miner.Bind())
		engine.Do()
	})

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return
}

func (a *Actor) Destroy() {

}

// Spawn ...
func (e *Actor) Spawn() {

}

// SetState ...
func (e *Actor) SetState(params supervisor.EntityStateParams) error {

	e.SetActorState(params.NewState)
	e.DeserializeAttr(params.AttributeValues)
	e.SaveState(false, params.StorageSave)

	return nil
}

// Update ...
func (e *Actor) Update() {

}

func (e *Actor) addAction(event events.EventCallEntityAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallEntityAction) {
	if action, ok := e.Actions[msg.ActionName]; ok {
		if action.ScriptEngine != nil && action.ScriptEngine.Engine() != nil {
			if _, err := action.ScriptEngine.Engine().AssertFunction(FuncEntityAction, msg.EntityId, action.Name, msg.Args); err != nil {
				log.Error(errors.Wrapf(err, "entity id: %s ", e.Id).Error())
			}
			return
		}
	}
	if e.ScriptsEngine != nil && e.ScriptsEngine.Engine() != nil {
		if _, err := e.ScriptsEngine.AssertFunction(FuncEntityAction, msg.EntityId, msg.ActionName, msg.Args); err != nil {
			log.Error(errors.Wrapf(err, "entity id: %s ", e.Id).Error())
		}
	}
}
