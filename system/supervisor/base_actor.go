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

package supervisor

import (
	"context"
	"fmt"
	"github.com/e154/smart-home/common/events"
	"runtime/debug"
	"sync"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
)

// BaseActor ...
type BaseActor struct {
	PluginActor
	Id                common.EntityId
	ParentId          *common.EntityId
	Name              string
	Description       string
	EntityType        string
	State             *ActorState
	Area              *m.Area
	Metric            []*m.Metric
	Hidden            bool
	AttrMu            *sync.RWMutex
	Attrs             m.Attributes
	Actions           map[string]ActorAction
	States            map[string]ActorState
	ScriptEngines     []*scripts.EngineWatcher
	Icon              *string
	ImageUrl          *string
	UnitOfMeasurement string
	Scripts           []*m.Script
	Value             *atomic.String
	AutoLoad          bool
	LastChanged       *time.Time
	LastUpdated       *time.Time
	SettingsMu        *sync.RWMutex
	Setts             m.Attributes
	Service           Service
	currentStateMu    *sync.RWMutex
	currentState      *bus.EventEntityState
}

// NewBaseActor ...
func NewBaseActor(entity *m.Entity,
	service Service) BaseActor {
	actor := BaseActor{
		Service:           service,
		Id:                common.EntityId(fmt.Sprintf("%s.%s", entity.PluginName, entity.Id.Name())),
		Name:              entity.Id.Name(),
		Description:       entity.Description,
		EntityType:        entity.PluginName,
		ParentId:          entity.ParentId,
		State:             nil,
		Area:              entity.Area,
		Hidden:            entity.Hidden,
		Actions:           make(map[string]ActorAction),
		States:            make(map[string]ActorState),
		Icon:              entity.Icon,
		ImageUrl:          nil,
		UnitOfMeasurement: "",
		Scripts:           entity.Scripts,
		Value:             atomic.NewString(StateAwait),
		LastChanged:       nil,
		LastUpdated:       nil,
		AutoLoad:          entity.AutoLoad,
		AttrMu:            &sync.RWMutex{},
		Attrs:             entity.Attributes.Copy(),
		SettingsMu:        &sync.RWMutex{},
		Setts:             entity.Settings,
		currentStateMu:    &sync.RWMutex{},
	}

	// Image
	if entity.Image != nil {
		actor.ImageUrl = common.String(entity.Image.Url)
	}

	// Metric
	actor.Metric = make([]*m.Metric, len(entity.Metrics))
	copy(actor.Metric, entity.Metrics)

	// Items
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
			if action.ScriptEngine, err = service.ScriptService().NewEngineWatcher(a.Script); err != nil {
				log.Error(err.Error())
			}
			action.ScriptEngine.Spawn(func(engine *scripts.Engine) {
				if _, err = engine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id)); err != nil {
					log.Error(err.Error())
				}
				if _, err = engine.Do(); err != nil {
					log.Error(err.Error())
				}
			})
		}

		if a.Image != nil {
			action.ImageUrl = common.String(a.Image.Url)
		}
		actor.Actions[a.Name] = action
	}

	if entity.Scripts != nil {
		for _, script := range entity.Scripts {
			var scriptEngine *scripts.EngineWatcher
			if scriptEngine, err = service.ScriptService().NewEngineWatcher(script); err == nil {
				scriptEngine.Spawn(func(engine *scripts.Engine) {
					if _, err = engine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id)); err != nil {
						log.Error(err.Error())
					}
				})
				if _, err = scriptEngine.Engine().Do(); err != nil {
					log.Error(err.Error())
				}
			}
			actor.ScriptEngines = append(actor.ScriptEngines, scriptEngine)
		}
	}

	return actor
}

func (e *BaseActor) StopWatchers() {
	for _, engine := range e.ScriptEngines {
		engine.Stop()
	}
	for _, a := range e.Actions {
		if a.ScriptEngine != nil {
			a.ScriptEngine.Stop()
		}
	}
}

// Metrics ...
func (e *BaseActor) Metrics() []*m.Metric {
	return e.Metric
}

// SetState ...
func (e *BaseActor) SetState(EntityStateParams) error {
	return apperr.ErrUnimplemented
}

// Attributes ...
func (e *BaseActor) Attributes() m.Attributes {
	e.attrLock()
	e.AttrMu.RLock()
	defer e.AttrMu.RUnlock()
	return e.Attrs.Copy()
}

// DeserializeAttr ...
func (e *BaseActor) DeserializeAttr(data m.AttributeValue) {
	e.attrLock()
	e.AttrMu.Lock()
	defer e.AttrMu.Unlock()
	_, _ = e.Attrs.Deserialize(data)
}

// Info ...
func (e *BaseActor) Info() (info ActorInfo) {
	info = ActorInfo{
		Id:                e.Id,
		PluginName:        e.EntityType,
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

// Now ...
func (e *BaseActor) Now(oldState bus.EventEntityState) time.Time {
	now := time.Now()
	e.LastUpdated = common.Time(now)

	if oldState.LastUpdated != nil {
		e.LastChanged = common.Time(*oldState.LastUpdated)
	} else {
		e.LastChanged = common.Time(*e.LastUpdated)
	}
	return now
}

// Settings ...
func (e *BaseActor) Settings() m.Attributes {
	e.settingsLock()
	e.SettingsMu.RLock()
	defer e.SettingsMu.RUnlock()
	return e.Setts.Copy()
}

// DeserializeSettings ...
func (e *BaseActor) DeserializeSettings(settings m.AttributeValue) {
	if settings == nil {
		return
	}
	e.settingsLock()
	e.SettingsMu.Lock()
	defer e.SettingsMu.Unlock()
	_, _ = e.Setts.Deserialize(settings)
}

func (e *BaseActor) attrLock() {
	if e.AttrMu == nil {
		e.AttrMu = &sync.RWMutex{}
	}
}

func (e *BaseActor) settingsLock() {
	if e.SettingsMu == nil { //todo: check race condition
		e.SettingsMu = &sync.RWMutex{}
	}
}

func (e *BaseActor) GetCurrentState() *bus.EventEntityState {
	e.currentStateMu.RLock()
	defer e.currentStateMu.RUnlock()
	return e.currentState
}

func (e *BaseActor) SetCurrentState(state bus.EventEntityState) {
	e.currentStateMu.Lock()
	e.currentState = &state
	e.currentStateMu.Unlock()
}

func (e *BaseActor) GetEventState() (eventState bus.EventEntityState) {

	attrs := e.Attributes()
	setts := e.Settings()

	var state *bus.EntityState

	info := e.Info()
	if info.State != nil {
		state = &bus.EntityState{
			Name:        info.State.Name,
			Description: info.State.Description,
			ImageUrl:    info.State.ImageUrl,
			Icon:        info.State.Icon,
		}
	}

	eventState = bus.EventEntityState{
		EntityId:   info.Id,
		Value:      info.Value,
		State:      state,
		Attributes: attrs,
		Settings:   setts,
	}

	if info.LastChanged != nil {
		eventState.LastChanged = common.Time(*info.LastChanged)
	}

	if info.LastUpdated != nil {
		eventState.LastUpdated = common.Time(*info.LastUpdated)
	}

	return
}

func (e *BaseActor) SaveState(msg events.EventStateChanged) {

	go e.updateMetric(msg.NewState)

	if msg.NewState.Compare(msg.OldState) {
		return
	}

	currentState := e.GetCurrentState()
	if currentState != nil && currentState.Compare(msg.NewState) {
		return
	}

	e.SetCurrentState(msg.NewState)

	// store state to db
	var state string
	if msg.NewState.State != nil {
		state = msg.NewState.State.Name
	}

	go e.Service.EventBus().Publish("system/entities/"+msg.EntityId.String(), msg)

	if !msg.StorageSave {
		return
	}

	go func() {
		_, err := e.Service.Adaptors().EntityStorage.Add(context.Background(), &m.EntityStorage{
			State:      state,
			EntityId:   msg.EntityId,
			Attributes: msg.NewState.Attributes.Serialize(),
		})
		if err != nil {
			log.Error(err.Error())
		}
	}()
}

func (e *BaseActor) updateMetric(state bus.EventEntityState) {

	if e.Metric == nil {
		return
	}

	var data = make(map[string]float32)
	var name string

	for _, metric := range e.Metric {
		for _, prop := range metric.Options.Items {
			if value, ok := state.Attributes[prop.Name]; ok {
				name = metric.Name
				switch value.Type {
				case common.AttributeInt:
					data[prop.Name] = float32(value.Int64())
				case common.AttributeFloat:
					data[prop.Name] = common.Rounding32(value.Float64(), 2)
				//case common.AttributePoint:
				//	data[prop.Name] = value.Point()
				}
			}
		}
	}

	if len(data) == 0 || name == "" {
		return
	}

	e.AddMetric(name, data)

}

func (e *BaseActor) AddMetric(name string, value map[string]float32) {

	if e.Metric == nil {
		return
	}

	var err error
	for _, metric := range e.Metric {
		if metric.Name != name {
			continue
		}

		if metric.Id == 0 {
			fmt.Printf("check metric for %s", e.Id.String())
			return
		}

		err = e.Service.Adaptors().MetricBucket.Add(context.Background(), &m.MetricDataItem{
			Value:    value,
			MetricId: metric.Id,
			Time:     time.Now(),
		})

		if err != nil {
			log.Errorf(err.Error(), value, metric.Id)
			debug.PrintStack()
		}
	}
}
