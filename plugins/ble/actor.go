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

package ble

import (
	"strings"

	"github.com/pkg/errors"
	"tinygo.org/x/bluetooth"

	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	actionPool chan events.EventCallEntityAction
	ble        *Ble
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) (actor *Actor) {

	actor = &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		actionPool: make(chan events.EventCallEntityAction, 1000),
		ble:        NewBle(),
	}

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			a.ScriptEngine.PushFunction("WriteGattChar", GetWriteGattCharBind(actor))
		}
	}

	if actor.ScriptsEngine != nil {
		actor.ScriptsEngine.PushFunction("WriteGattChar", GetWriteGattCharBind(actor))
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	if actor.Actions == nil {
		actor.Actions = NewActions()
	}

	return actor
}

func (e *Actor) Destroy() {

}

func (e *Actor) Spawn() {
	e.BaseActor.Spawn()
	e.Service.ScriptService().PushFunctions("WriteGattChar", GetWriteGattCharBind(e))
}

// SetState ...
func (e *Actor) SetState(params supervisor.EntityStateParams) error {

	e.SetActorState(params.NewState)
	e.DeserializeAttr(params.AttributeValues)
	e.SaveState(false, params.StorageSave)

	return nil
}

func (e *Actor) addAction(event events.EventCallEntityAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallEntityAction) {

	if strings.ToUpper(msg.ActionName) == ActionScan {
		address, err := bluetooth.ParseUUID(e.Setts[AttrAddress].String())
		if err != nil {
			e.ble.Scan(nil)
		} else {
			e.ble.Scan(&address)
		}
	}

	if action, ok := e.Actions[msg.ActionName]; ok {
		if action.ScriptEngine != nil && action.ScriptEngine.Engine() != nil {
			if _, err := action.ScriptEngine.Engine().AssertFunction(FuncEntityAction, e.Id, action.Name, msg.Args); err != nil {
				log.Error(errors.Wrapf(err, "entity id: %s ", e.Id).Error())
			}
			return
		}
	}
	if e.ScriptsEngine != nil && e.ScriptsEngine.Engine() != nil {
		if _, err := e.ScriptsEngine.AssertFunction(FuncEntityAction, e.Id, msg.ActionName, msg.Args); err != nil {
			log.Error(errors.Wrapf(err, "entity id: %s ", e.Id).Error())
		}
	}
}
