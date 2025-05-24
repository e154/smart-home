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
	"fmt"
	"strings"

	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	actionPool chan events.EventCallEntityAction
	ble        Bluetooth
}

// NewActor ...
func NewActor(entity *m.Entity,
	service plugins.Service) (actor *Actor) {

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
			a.ScriptEngine.PushFunction("BleWrite", GetWriteBind(actor))
			a.ScriptEngine.PushFunction("BleRead", GetReadBind(actor))
			//a.ScriptEngine.PushFunction("BleSubscribe", GetSubscribeBind(actor))
			//a.ScriptEngine.PushFunction("BleDisconnect", GetDisconnectBind(actor))
		}
	}

	if actor.ScriptsEngine != nil {
		actor.ScriptsEngine.PushFunction("BleWrite", GetWriteBind(actor))
		actor.ScriptsEngine.PushFunction("BleRead", GetReadBind(actor))
		//actor.ScriptsEngine.PushFunction("BleSubscribe", GetSubscribeBind(actor))
		//actor.ScriptsEngine.PushFunction("BleDisconnect", GetDisconnectBind(actor))
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

	var timeout, connectionTimeout = DefaultTimeout, DefaultConnectionTimeout
	if e.Setts[AttrConnectionTimeoutSec] != nil {
		connectionTimeout = e.Setts[AttrConnectionTimeoutSec].Int64()
	}

	if e.Setts[AttrTimeoutSec] != nil {
		timeout = e.Setts[AttrTimeoutSec].Int64()
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
func (e *Actor) SetState(params plugins.EntityStateParams) error {

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
		if address, ok := e.Setts[AttrAddress]; ok {
			e.ble.Scan(common.String(address.String()))
		} else {
			e.ble.Scan(nil)
		}
	}

	if action, ok := e.Actions[msg.ActionName]; ok {
		if action.ScriptEngine != nil && action.ScriptEngine.Engine() != nil {
			if _, err := action.ScriptEngine.Engine().AssertFunction(FuncEntityAction, e.Id, action.Name, msg.Args); err != nil {
				log.Error(fmt.Errorf("entity id: %s : %w", e.Id, err).Error())
			}
			return
		}
	}
	if e.ScriptsEngine != nil && e.ScriptsEngine.Engine() != nil {
		if _, err := e.ScriptsEngine.AssertFunction(FuncEntityAction, e.Id, msg.ActionName, msg.Args); err != nil {
			log.Error(fmt.Errorf("entity id: %s : %w", e.Id, err).Error())
		}
	}
}
