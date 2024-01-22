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
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
)

// BaseActor ...
type BaseActor struct {
	PluginActor
	Id                common.EntityId
	ParentId          *common.EntityId
	Name              string
	Description       string
	EntityType        string
	Area              *m.Area
	Metric            []*m.Metric
	Hidden            bool
	AttrMu            *sync.RWMutex
	Attrs             m.Attributes
	Actions           map[string]ActorAction
	States            map[string]ActorState
	ScriptsEngine     *scripts.EnginesWatcher
	Icon              *string
	ImageUrl          *string
	UnitOfMeasurement string
	Scripts           []*m.Script
	Value             *atomic.String
	AutoLoad          bool
	LastChanged       *time.Time
	LastUpdated       *time.Time
	actorStateMu      *sync.RWMutex
	state             *ActorState
	SettingsMu        *sync.RWMutex
	Setts             m.Attributes
	Service           Service
	stateMu           *sync.RWMutex
	currentState      *bus.EventEntityState
	oldState          *bus.EventEntityState
}

// NewBaseActor ...
func NewBaseActor(entity *m.Entity,
	service Service) *BaseActor {
	actor := &BaseActor{
		Service:           service,
		Id:                common.EntityId(fmt.Sprintf("%s.%s", entity.PluginName, entity.Id.Name())),
		Name:              entity.Id.Name(),
		Description:       entity.Description,
		EntityType:        entity.PluginName,
		ParentId:          entity.ParentId,
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
		actorStateMu:      &sync.RWMutex{},
		state:             nil,
		AttrMu:            &sync.RWMutex{},
		Attrs:             entity.Attributes.Copy(),
		SettingsMu:        &sync.RWMutex{},
		Setts:             entity.Settings,
		stateMu:           &sync.RWMutex{},
	}

	// Image
	if entity.Image != nil {
		actor.ImageUrl = common.String(entity.Image.Url)
	}

	// Metric
	actor.Metric = make([]*m.Metric, len(entity.Metrics))
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
		if actor.ScriptsEngine, err = service.ScriptService().NewEnginesWatcher(entity.Scripts); err != nil {
			log.Error(err.Error())
		}
		actor.ScriptsEngine.Spawn(func(engine *scripts.Engine) {
			if _, err = engine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id)); err != nil {
				log.Error(err.Error())
			}
			if _, err = engine.Do(); err != nil {
				log.Error(err.Error())
			}
		})
		if _, err = actor.ScriptsEngine.Engine().AssertFunction("init"); err != nil {
			log.Error(err.Error())
		}
	}

	// restore state
	actor.RestoreState(entity)

	return actor
}

func (e *BaseActor) StopWatchers() {
	e.ScriptsEngine.Stop()
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

// SetActorState ...
func (e *BaseActor) SetActorState(name *string) {
	if name == nil {
		return
	}
	e.actorStateLock()
	e.actorStateMu.Lock()
	defer e.actorStateMu.Unlock()
	if state, ok := e.States[*name]; ok {
		e.state = &state
	}
}

// SetActorStateImage ...
func (e *BaseActor) SetActorStateImage(imageUrl, icon *string) {
	e.actorStateLock()
	e.actorStateMu.Lock()
	defer e.actorStateMu.Unlock()
	if e.state == nil {
		return
	}
	e.state.ImageUrl = imageUrl
	e.state.Icon = icon
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
		State:             e.state,
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

func (e *BaseActor) actorStateLock() {
	if e.actorStateMu == nil {
		e.actorStateMu = &sync.RWMutex{}
	}
}

func (e *BaseActor) settingsLock() {
	if e.SettingsMu == nil { //todo: check race condition
		e.SettingsMu = &sync.RWMutex{}
	}
}

func (e *BaseActor) GetCurrentState() *bus.EventEntityState {
	e.stateMu.RLock()
	defer e.stateMu.RUnlock()
	if e.currentState != nil {
		return e.currentState
	}
	currentState := e.GetEventState()
	e.currentState = &currentState
	return e.currentState
}

func (e *BaseActor) GetOldState() *bus.EventEntityState {
	e.stateMu.RLock()
	defer e.stateMu.RUnlock()
	return e.oldState
}

func (e *BaseActor) SetCurrentState(state bus.EventEntityState) {
	e.stateMu.Lock()
	e.currentState = &state
	e.stateMu.Unlock()
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
		eventState.LastChanged = common.Time(*e.LastChanged)
	}

	if info.LastUpdated != nil {
		eventState.LastUpdated = common.Time(*e.LastUpdated)
	}

	return
}

func (e *BaseActor) restoreStore(entity *m.Entity, store *m.EntityStorage, state *bus.EventEntityState) {
	if store.State != "" {
		for _, s := range entity.States {
			if store.State == s.Name {
				es := &bus.EntityState{
					Name:        s.Name,
					Description: s.Description,
					Icon:        s.Icon,
				}
				if s.Image != nil {
					es.ImageUrl = common.String(s.Image.Url)
				}
				state.State = es
			}
		}
	}
	// Attributes
	state.Attributes = entity.Attributes.Copy()
	_, _ = state.Attributes.Deserialize(store.Attributes)
	// Settings
	state.Settings = e.Settings()
	// LastChanged
	state.LastChanged = common.Time(store.CreatedAt)
	// LastUpdated
	state.LastUpdated = common.Time(store.CreatedAt)
}

func (e *BaseActor) RestoreState(entity *m.Entity) {
	e.stateMu.RLock()
	defer e.stateMu.RUnlock()
	if len(entity.Storage) > 0 {
		e.currentState = &bus.EventEntityState{
			EntityId: entity.Id,
		}
		var store = entity.Storage[0]
		e.LastUpdated = common.Time(store.CreatedAt)
		e.restoreStore(entity, store, e.currentState)
	}
	if len(entity.Storage) > 1 {
		e.oldState = &bus.EventEntityState{
			EntityId: entity.Id,
		}
		var store = entity.Storage[1]
		e.LastChanged = common.Time(store.CreatedAt)
		e.currentState.LastChanged = common.Time(store.CreatedAt)
		e.restoreStore(entity, store, e.oldState)
	}
}

func (e *BaseActor) SaveState(doNotSaveMetric, storageSave bool) {

	e.stateMu.RLock()
	defer e.stateMu.RUnlock()

	newState := e.GetEventState()

	if !doNotSaveMetric {
		go e.updateMetric(newState)
	}

	if e.currentState != nil && e.currentState.Compare(newState) {
		return
	}

	newState.LastUpdated = common.Time(time.Now())
	e.LastUpdated = common.Time(*newState.LastUpdated)

	if e.currentState != nil {
		if e.currentState.LastUpdated != nil {
			newState.LastChanged = common.Time(*e.currentState.LastUpdated)
			e.LastChanged = common.Time(*e.currentState.LastUpdated)
		}
		if e.oldState != nil {
			*e.oldState = *e.currentState
		} else {
			e.oldState = e.currentState
		}
	}

	e.currentState = &newState

	msg := events.EventStateChanged{
		StorageSave:     storageSave,
		PluginName:      e.Id.PluginName(),
		EntityId:        e.Id,
		NewState:        newState,
		DoNotSaveMetric: doNotSaveMetric,
	}
	if e.oldState != nil {
		msg.OldState = *e.oldState
	}
	e.Service.EventBus().Publish("system/entities/"+e.Id.String(), msg)

	if !storageSave {
		return
	}

	go func() {

		var state string
		if newState.State != nil {
			state = newState.State.Name
		}

		_, err := e.Service.Adaptors().EntityStorage.Add(context.Background(), &m.EntityStorage{
			State:      state,
			EntityId:   e.Id,
			Attributes: newState.Attributes.Serialize(),
			CreatedAt:  *newState.LastUpdated,
		})
		if err != nil {
			//log.Error(err.Error())
		}
	}()
}

func (e *BaseActor) updateMetric(state bus.EventEntityState) {

	if e.Metric == nil {
		return
	}

	var data = make(map[string]interface{})
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
				case common.AttributeString:
					data[prop.Name] = value.String()
				case common.AttributePoint:
					data[prop.Name] = value.Point()
				case common.AttributeMap:
					data[prop.Name] = value.Map()
				default:
					log.Warnf("unimplemented type %s", value.Type)
				}
			}
		}
	}

	if len(data) == 0 || name == "" {
		return
	}

	e.AddMetric(name, data)

}

func (e *BaseActor) AddMetric(name string, value map[string]interface{}) {

	if e.Metric == nil {
		return
	}

	var updated bool

	var err error
	for _, metric := range e.Metric {
		if metric.Name != name {
			continue
		}

		if metric.Id == 0 {
			log.Debugf("check metric for %s", e.Id.String())
			return
		}

		err = e.Service.Adaptors().MetricBucket.Add(context.Background(), &m.MetricDataItem{
			Value:    value,
			MetricId: metric.Id,
			Time:     time.Now(),
		})

		if err != nil {
			log.Errorf(err.Error(), value, metric.Id)
		}

		updated = true
	}

	if !updated {
		return
	}

	e.Service.EventBus().Publish("system/entities/%s"+e.Id.String(), events.EventUpdatedMetric{
		EntityId: e.Id,
	})
}
