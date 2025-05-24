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

	bitmine2 "github.com/e154/smart-home/internal/plugins/cgminer/bitmine"
	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/events"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	miner      IMiner
	actionPool chan events.EventCallEntityAction
}

// NewActor ...
func NewActor(entity *m.Entity,
	service plugins.Service) (actor *Actor, err error) {

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
		err = fmt.Errorf("actor 'cgminer', current settings %+v: %w", actor.Setts, apperr.ErrBadSettings)
		return
	}

	switch actor.Setts[SettingManufacturer].Value {
	case bitmine2.ManufactureBitmine:
		switch actor.Setts[SettingModel].Value {
		case bitmine2.DeviceS9:
		case bitmine2.DeviceS7:
		case bitmine2.DeviceL3:
		case bitmine2.DeviceL3Plus:
		case bitmine2.DeviceD3:
		case bitmine2.DeviceT9:
		default:
			err = fmt.Errorf("unknown model %s", actor.Setts[SettingModel].Value)
		}
	default:
		err = fmt.Errorf("unknown manufacture %s", actor.Setts[SettingManufacturer].Value)
	}

	if _, ok := actor.Setts[SettingHost]; !ok {
		err = fmt.Errorf("actor 'cgminer', current settings %+v: %w", actor.Setts, apperr.ErrBadSettings)
		return
	}

	if _, ok := actor.Setts[SettingPort]; !ok {
		err = fmt.Errorf("actor 'cgminer', current settings %+v: %w", actor.Setts, apperr.ErrBadSettings)
		return
	}

	if _, ok := actor.Setts[SettingTimeout]; !ok {
		err = fmt.Errorf("actor 'cgminer', current settings %+v: %w", actor.Setts, apperr.ErrBadSettings)
		return
	}

	if _, ok := actor.Setts[SettingModel]; !ok {
		err = fmt.Errorf("actor 'cgminer', current settings %+v: %w", actor.Setts, apperr.ErrBadSettings)
		return
	}

	if _, ok := actor.Setts[SettingHost]; !ok {
		err = fmt.Errorf("actor 'cgminer', current settings %+v: %w", actor.Setts, apperr.ErrBadSettings)
		return
	}

	if _, ok := actor.Setts[SettingUser]; !ok {
		err = fmt.Errorf("actor 'cgminer', current settings %+v: %w", actor.Setts, apperr.ErrBadSettings)
		return
	}

	if _, ok := actor.Setts[SettingPass]; !ok {
		err = fmt.Errorf("actor 'cgminer', current settings %+v: %w", actor.Setts, apperr.ErrBadSettings)
		return
	}

	host := actor.Setts[SettingHost].String()
	port := int(actor.Setts[SettingPort].Int64())
	timeout := actor.Setts[SettingTimeout].Int64()
	model := actor.Setts[SettingModel].String()
	user := actor.Setts[SettingUser].String()
	pass := actor.Setts[SettingPass].String()
	actor.miner, err = bitmine2.NewBitmine(bitmine2.NewTransport(host, port, timeout), model, user, pass)

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			a.ScriptEngine.PushFunction("Miner", actor.miner.Bind())
		}
	}

	if actor.ScriptsEngine != nil && actor.ScriptsEngine.Engine() != nil {
		actor.ScriptsEngine.PushFunction("Miner", actor.miner.Bind())
	}

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

// SetState ...
func (e *Actor) SetState(params plugins.EntityStateParams) error {

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
			if _, err := action.ScriptEngine.Engine().AssertFunction(FuncEntityAction, e.Id, action.Name, msg.Args); err != nil {
				log.Error(fmt.Errorf("entity id: %s: %w", e.Id, err).Error())
			}
			return
		}
	}
	if e.ScriptsEngine != nil && e.ScriptsEngine.Engine() != nil {
		if _, err := e.ScriptsEngine.AssertFunction(FuncEntityAction, e.Id, msg.ActionName, msg.Args); err != nil {
			log.Error(fmt.Errorf("entity id: %s: %w", e.Id, err).Error())
		}
	}
}
