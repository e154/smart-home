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

package entity_manager

import (
	"time"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
)

// PluginActor ...
type PluginActor interface {

	// Spawn ...
	Spawn() PluginActor

	// Attributes ...
	Attributes() m.Attributes

	// Settings ...
	Settings() m.Attributes

	// Metrics ...
	Metrics() []m.Metric

	// SetState ...
	SetState(EntityStateParams) error

	// Info ...
	Info() ActorInfo
}

// ActorConstructor ...
type ActorConstructor func() PluginActor

// EntityManager ...
type EntityManager interface {

	// SetPluginManager ...
	SetPluginManager(pluginManager common.PluginManager)

	// LoadEntities ...
	LoadEntities()

	// Shutdown ...
	Shutdown()

	// SetMetric ...
	SetMetric(common.EntityId, string, map[string]interface{})

	// SetState ...
	SetState(common.EntityId, EntityStateParams) error

	// GetEntityById ...
	GetEntityById(common.EntityId) (m.EntityShort, error)

	// GetActorById ...
	GetActorById(common.EntityId) (PluginActor, error)

	// List ...
	List() ([]m.EntityShort, error)

	// Spawn ...
	Spawn(ActorConstructor) PluginActor

	// Remove ...
	Remove(common.EntityId)

	// CallAction ...
	CallAction(common.EntityId, string, map[string]interface{})

	// CallScene ...
	CallScene(common.EntityId, map[string]interface{})

	// Add ...
	Add(*m.Entity) error

	// Update ...
	Update(*m.Entity) error
}

// ActorAction ...
type ActorAction struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	ImageUrl     *string         `json:"image_url"`
	Icon         *string         `json:"icon"`
	ScriptEngine *scripts.Engine `json:"-"`
}

// ToEntityActionShort ...
func ToEntityActionShort(from map[string]ActorAction) (to map[string]m.EntityActionShort) {
	to = make(map[string]m.EntityActionShort)
	for k, v := range from {
		to[k] = m.EntityActionShort{
			Name:        v.Name,
			Description: v.Description,
			ImageUrl:    v.ImageUrl,
			Icon:        v.Icon,
		}
	}
	return
}

// ActorState ...
type ActorState struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    *string `json:"image_url"`
	Icon        *string `json:"icon"`
}

// ToEntityStateShort ...
func ToEntityStateShort(from map[string]ActorState) (to map[string]m.EntityStateShort) {
	to = make(map[string]m.EntityStateShort)
	for k, v := range from {
		to[k] = m.EntityStateShort{
			Name:        v.Name,
			Description: v.Description,
			ImageUrl:    v.ImageUrl,
			Icon:        v.Icon,
		}
	}
	return
}

// Copy ...
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

// Message ...
type Message struct {
	From    common.EntityId
	To      common.EntityId
	Payload interface{}
}

const (
	// StateAwait ...
	StateAwait = "await"
	// StateOk ...
	StateOk = "ok"
	// StateError ...
	StateError = "error"
	// StateInProcess ...
	StateInProcess = "in process"
)

type actorInfo struct {
	Actor    PluginActor
	quit     chan struct{}
	OldState event_bus.EventEntityState
}

// ActorInfo ...
type ActorInfo struct {
	Id                common.EntityId        `json:"id"`
	ParentId          *common.EntityId       `json:"parent_id"`
	PluginName        string                 `json:"plugin_name"`
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
	AutoLoad          bool                   `json:"auto_load"`
	Value             interface{}            `json:"value"`
	States            map[string]ActorState  `json:"states"`
	Actions           map[string]ActorAction `json:"actions"`
}
