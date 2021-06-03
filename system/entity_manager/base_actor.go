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
	"errors"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type BaseActor struct {
	PluginActor
	Id                common.EntityId
	ParentId          *common.EntityId
	Name              string
	Description       string
	EntityType        common.EntityType
	Manager           EntityManager
	State             *ActorState
	Area              *m.Area
	Metric            []m.Metric
	Hidden            bool
	AttrMu            *sync.Mutex
	Attrs             m.Attributes
	Actions           map[string]ActorAction
	States            map[string]ActorState
	ScriptEngine      *scripts.Engine
	Icon              *common.Icon
	ImageUrl          *string
	UnitOfMeasurement string
	Sripts            []m.Script
	Value             *atomic.String
	AutoLoad          bool
	LastChanged       *time.Time
	LastUpdated       *time.Time
	adaptors          *adaptors.Adaptors
	SettingsMu        *sync.Mutex
	Setts             m.Attributes
}

func NewBaseActor(entity *m.Entity,
	scriptService scripts.ScriptService,
	adaptors *adaptors.Adaptors) BaseActor {
	actor := BaseActor{
		adaptors:          adaptors,
		Id:                common.EntityId(fmt.Sprintf("%s.%s", entity.Type, entity.Id.Name())),
		Name:              entity.Id.Name(),
		Description:       entity.Description,
		EntityType:        entity.Type,
		ParentId:          entity.ParentId,
		Manager:           nil,
		State:             nil,
		Area:              entity.Area,
		Hidden:            entity.Hidden,
		Actions:           make(map[string]ActorAction),
		States:            make(map[string]ActorState),
		Icon:              entity.Icon,
		ImageUrl:          nil,
		UnitOfMeasurement: "",
		Sripts:            entity.Scripts,
		Value:             nil,
		LastChanged:       nil,
		LastUpdated:       nil,
		AutoLoad:          entity.AutoLoad,
		AttrMu:            &sync.Mutex{},
		Attrs:             entity.Attributes.Copy(),
		SettingsMu:        &sync.Mutex{},
		Setts:             entity.Settings,
	}

	// Image
	if entity.Image != nil {
		actor.ImageUrl = common.String(entity.Image.Url)
	}

	// Metric
	actor.Metric = make([]m.Metric, len(entity.Metrics))
	copy(actor.Metric, entity.Metrics)

	// States
	for _, s := range entity.States {
		state := ActorState{
			Name:        s.Name,
			Description: s.Description,
			Icon:        s.Icon,
		}
		if s.Image != nil {
			state.ImageUrl = common.String(s.Image.Url)
		}
		actor.States[s.Name] = state
	}

	var err error
	// Actions
	for _, a := range entity.Actions {
		action := ActorAction{
			Name:        a.Name,
			Description: a.Description,
			Icon:        a.Icon,
		}

		if a.Script != nil {
			if action.ScriptEngine, err = scriptService.NewEngine(a.Script); err != nil {
				log.Error(err.Error())
			}
		}

		if a.Image != nil {
			action.ImageUrl = common.String(a.Image.Url)
		}
		actor.Actions[a.Name] = action
	}

	// Scripts
	if len(entity.Scripts) != 0 {
		if actor.ScriptEngine, err = scriptService.NewEngine(&entity.Scripts[0]); err != nil {
			log.Error(err.Error())
		}

		actor.ScriptEngine.Do()

	} else {
		if actor.ScriptEngine, err = scriptService.NewEngine(nil); err != nil {
			log.Error(err.Error())
		}
	}

	return actor
}

func (b *BaseActor) GetEventState(actor PluginActor) event_bus.EventEntityState {
	return GetEventState(actor)
}

func (e *BaseActor) Metrics() []m.Metric {
	return e.Metric
}

func (e *BaseActor) SetState(EntityStateParams) error {
	return errors.New("method not implemented")
}

func (e *BaseActor) Attributes() m.Attributes {
	e.attrLock()
	return e.Attrs.Copy()
}

func (e *BaseActor) DeserializeAttr(data m.AttributeValue) {
	e.attrLock()
	e.Attrs.Deserialize(data)
}

func (e *BaseActor) Info() (info ActorInfo) {
	info = ActorInfo{
		Id:                e.Id,
		Type:              e.EntityType,
		Name:              e.Name,
		Description:       e.Description,
		Hidde:             e.Hidden,
		DependsOn:         nil,
		State:             e.State,
		ImageUrl:          e.ImageUrl,
		Icon:              e.Icon,
		Area:              e.Area,
		UnitOfMeasurement: e.UnitOfMeasurement,
		LastChanged:       e.LastChanged,
		LastUpdated:       e.LastUpdated,
		Actions:           e.Actions,
		States:            e.States,
		AutoLoad:          e.AutoLoad,
		ParentId:          e.ParentId,
		//Value:             e.value,
	}
	return
}

func (e *BaseActor) Now(oldState event_bus.EventEntityState) time.Time {
	now := time.Now()
	e.LastUpdated = common.Time(now)

	if oldState.LastUpdated != nil {
		e.LastChanged = common.Time(*oldState.LastUpdated)
	} else {
		e.LastChanged = common.Time(*e.LastUpdated)
	}
	return now
}

func (e *BaseActor) SetMetric(id common.EntityId, name string, value map[string]interface{}) {
	if e.Manager != nil {
		e.Manager.SetMetric(id, name, value)
	}
}

func (e *BaseActor) Settings() m.Attributes {
	e.settingsLock()
	return e.Setts.Copy()
}

func (e *BaseActor) DeserializeSettings(settings m.AttributeValue) {
	if settings == nil {
		return
	}
	e.settingsLock()
	e.Setts.Deserialize(settings)
}

func (e *BaseActor) attrLock() {
	if e.AttrMu == nil {
		e.AttrMu = &sync.Mutex{}
	}
	e.AttrMu.Lock()
	defer e.AttrMu.Unlock()
}

func (e *BaseActor) settingsLock() {
	if e.SettingsMu == nil {
		e.SettingsMu = &sync.Mutex{}
	}
	e.SettingsMu.Lock()
	defer e.SettingsMu.Unlock()
}
