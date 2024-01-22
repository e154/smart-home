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

package onvif

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/media"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	client      *Client
	snapshotUri *string
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) (actor *Actor) {

	actor = &Actor{
		BaseActor: supervisor.NewBaseActor(entity, service),
	}

	actor.client = NewClient(actor.eventHandler)

	clientBind := NewClientBind(actor.client)

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine.Engine() != nil {
			a.ScriptEngine.Spawn(func(engine *scripts.Engine) {
				engine.PushStruct("Camera", clientBind)
				_, _ = engine.Do()
			})
		}
	}

	actor.ScriptsEngine.Spawn(func(engine *scripts.Engine) {
		engine.PushStruct("Camera", clientBind)
		engine.Do()
	})

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

func (a *Actor) Destroy() {
	a.Service.EventBus().Publish("system/media", media.EventRemoveList{Name: a.Id.String()})
	go a.client.Shutdown()
}

// Spawn ...
func (a *Actor) Spawn() {
	a.client.Start(a.Setts[AttrUserName].String(),
		a.Setts[AttrPassword].Decrypt(),
		a.Setts[AttrAddress].String(),
		a.Setts[AttrOnvifPort].Int64(),
		a.Setts[AttrRequireAuthorization].Bool())
}

// SetState ...
func (a *Actor) SetState(params supervisor.EntityStateParams) error {

	a.SetActorState(params.NewState)
	a.DeserializeAttr(params.AttributeValues)
	a.SaveState(false, params.StorageSave)

	return nil
}

func (a *Actor) addAction(event events.EventCallEntityAction) {
	a.runAction(event)
}

func (a *Actor) runAction(msg events.EventCallEntityAction) {
	action, ok := a.Actions[msg.ActionName]
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

func (a *Actor) eventHandler(msg interface{}) {
	switch v := msg.(type) {
	case *StreamList:
		go a.prepareStreamList(v)
	case *ConnectionStatus:
		go a.updateState(v)
	case *MotionAlarm:
		go a.prepareMotionAlarm(v)
	}
}

func (a *Actor) updateState(event *ConnectionStatus) {
	info := a.Info()
	var newStat = AttrOffline
	if event.Connected {
		newStat = AttrConnected
	}
	if info.State != nil && info.State.Name == newStat {
		return
	}
	a.SetState(supervisor.EntityStateParams{
		NewState:    common.String(newStat),
		StorageSave: true,
	})
}

func (a *Actor) prepareMotionAlarm(event *MotionAlarm) {
	a.SetState(supervisor.EntityStateParams{
		NewState: common.String(AttrConnected),
		AttributeValues: m.AttributeValue{
			AttrMotion:     event.State,
			AttrMotionTime: event.Time,
		},
		StorageSave: true,
	})
}

func (a *Actor) prepareStreamList(event *StreamList) {
	a.snapshotUri = event.SnapshotUri
	a.Service.EventBus().Publish("system/media", media.EventUpdateList{
		Name:     a.Id.String(),
		Channels: event.List,
	})
}

func (a *Actor) GetSnapshotUri() string {
	return common.StringValue(a.snapshotUri)
}
