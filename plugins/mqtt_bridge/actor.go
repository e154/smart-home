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

package mqtt_bridge

import (
	"context"
	"strings"

	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	actionPool chan events.EventCallEntityAction
	bridge     *MqttBridge
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) (actor *Actor) {

	actor = &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		actionPool: make(chan events.EventCallEntityAction, 1000),
	}

	var topics = strings.Split(entity.Settings[AttrTopics].String(), ",")

	config := &Config{
		KeepAlive:      int(entity.Settings[AttrKeepAlive].Int64()),
		PingTimeout:    int(entity.Settings[AttrPingTimeout].Int64()),
		Broker:         entity.Settings[AttrBroker].String(),
		ClientID:       entity.Settings[AttrClientID].String(),
		ConnectTimeout: int(entity.Settings[AttrConnectTimeout].Int64()),
		CleanSession:   entity.Settings[AttrCleanSession].Bool(),
		Username:       entity.Settings[AttrUsername].String(),
		Password:       entity.Settings[AttrPassword].Decrypt(),
		Qos:            byte(entity.Settings[AttrQos].Int64()),
		Direction:      Direction(entity.Settings[AttrDirection].String()),
		Topics:         topics,
	}

	var err error
	if actor.bridge, err = NewMqttBridge(config, service.MqttServ(), actor); err != nil {
		log.Error(err.Error())
	}

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return actor
}

func (e *Actor) Destroy() {
	close(e.actionPool)
	if err := e.bridge.Shutdown(context.Background()); err != nil {
		log.Error(err.Error())
	}
}

// Spawn ...
func (e *Actor) Spawn() {
	if err := e.bridge.Start(context.Background()); err != nil {
		log.Error(err.Error())
	}
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
	action, ok := e.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if action.ScriptEngine.Engine() == nil {
		return
	}
	if _, err := action.ScriptEngine.Engine().AssertFunction(FuncEntityAction, msg.EntityId, action.Name, msg.Args); err != nil {
		log.Error(err.Error())
	}
}
