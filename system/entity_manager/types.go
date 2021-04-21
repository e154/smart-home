// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package entity_manager

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
	"time"
)

type IMessage interface{}

type IActor interface {

	// Spawn ...
	Spawn(system IActorManager) IActor

	// Receive ...
	Receive(message Message)

	// Attributes ...
	Attributes() m.EntityAttributes

	// Destroy ...
	Destroy()

	// Metrics ...
	Metrics() []m.Metric

	// SetState
	SetState(EntityStateParams)

	// Info
	Info() ActorInfo
}

type IConstructor func(system IActorManager) (state IActor)

type IActorManager interface {
	// Spawn ...
	Spawn(constructor IConstructor) IActor

	// Destroy ...
	Destroy(id common.EntityId)

	// Send ...
	Send(Message) error

	// Broadcast ...
	Broadcast(Message)

	// Shutdown ...
	Shutdown()

	// SetMetric ...
	SetMetric(id common.EntityId, name string, value map[string]interface{})
}

type EntityAttribute struct {
	Name  string                     `json:"name"`
	Type  common.EntityAttributeType `json:"type"`
	Value interface{}                `json:"value,omitempty"`
}

type EntityAction struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    *string `json:"image_url"`
	Icon        *string `json:"icon"`
}

type ActorAction struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	ImageUrl     *string         `json:"image_url"`
	Icon         *string         `json:"icon"`
	ScriptEngine *scripts.Engine `json:"-"`
}

type EntityState struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    *string `json:"image_url"`
	Icon        *string `json:"icon"`
}

type ActorState struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    *string `json:"image_url"`
	Icon        *string `json:"icon"`
}

func (a *ActorState) Copy() (state *ActorState) {

	if a == nil {
		return nil
	}

	state = &ActorState{
		Name:        a.Name,
		Description: a.Description,
	}
	if a.ImageUrl != nil {
		state.ImageUrl = common.String(*a.ImageUrl)
	}
	if a.Icon != nil {
		state.Icon = common.String(*a.Icon)
	}
	return
}

type Message struct {
	From    common.EntityId
	To      common.EntityId
	Payload interface{}
}

const (
	StateAwait     = "await"
	StateOk        = "ok"
	StateError     = "error"
	StateInProcess = "in process"
)

type MessageCallAction struct {
	Name string                 `json:"name"`
	Arg  map[string]interface{} `json:"arg"`
}

type MessageCallScene struct {
	Arg map[string]interface{} `json:"arg"`
}

type actorInfo struct {
	Actor    IActor
	Queue    chan Message
	OldState event_bus.EventEntityState
}

type ActorInfo struct {
	Id                common.EntityId        `json:"id"`
	Type              common.EntityType      `json:"type"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Hidde             bool                   `json:"hidde"`
	UnitOfMeasurement string                 `json:"unit_of_measurement"`
	LastChanged       *time.Time             `json:"last_changed"`
	LastUpdated       *time.Time             `json:"last_updated"`
	DependsOn         []string               `json:"depends_on"`
	State             *ActorState            `json:"state"`
	ImageUrl          *string                `json:"image_url"`
	Icon              *common.Icon           `json:"icon"`
	Area              *m.Area                `json:"area"`
	Value             interface{}            `json:"value"`
	States            map[string]ActorState  `json:"states"`
	Actions           map[string]ActorAction `json:"actions"`
}
