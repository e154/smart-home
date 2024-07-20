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
			a.ScriptEngine.PushFunction("ReadGattChar", GetReadGattCharBind(actor))
			a.ScriptEngine.PushFunction("SubscribeGatt", GetSubscribeGattBind(actor))
			a.ScriptEngine.PushFunction("DisconnectGatt", GetDisconnectBind(actor))
		}
	}

	if actor.ScriptsEngine != nil {
		actor.ScriptsEngine.PushFunction("WriteGattChar", GetWriteGattCharBind(actor))
		actor.ScriptsEngine.PushFunction("ReadGattChar", GetReadGattCharBind(actor))
		actor.ScriptsEngine.PushFunction("SubscribeGatt", GetSubscribeGattBind(actor))
		actor.ScriptsEngine.PushFunction("DisconnectGatt", GetDisconnectBind(actor))
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
	e.ble.Disconnect()
}

func (e *Actor) Spawn() {

	var timeout, connectionTimeout int64 = 5, 5
	if e.Setts[AttrConnectionTimeoutSec] != nil {
		connectionTimeout = e.Setts[AttrConnectionTimeoutSec].Int64()
	}

	if e.Setts[AttrTimeoutSec] != nil {
		connectionTimeout = e.Setts[AttrTimeoutSec].Int64()
	}

	var address string
	if e.Setts[AttrAddress] != nil {
		address = e.Setts[AttrAddress].String()
	}

	var debug bool
	if e.Setts[AttrDebug] != nil {
		debug = e.Setts[AttrDebug].Bool()
	}

	e.ble = NewBle(address, timeout, connectionTimeout, debug)

	e.BaseActor.Spawn()
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
