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
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/fx"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/event_bus/events"
	"github.com/e154/smart-home/system/scripts"
)

var (
	log = logger.MustGetLogger("entity.manager")
)

type entityManager struct {
	eventBus      event_bus.EventBus
	adaptors      *adaptors.Adaptors
	scripts       scripts.ScriptService
	pluginManager common.PluginManager
	actors        sync.Map
	quit          chan struct{}
}

// NewEntityManager ...
func NewEntityManager(lc fx.Lifecycle,
	eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors,
	scripts scripts.ScriptService) EntityManager {
	manager := &entityManager{
		eventBus: eventBus,
		adaptors: adaptors,
		scripts:  scripts,
		actors:   sync.Map{},
		quit:     make(chan struct{}),
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) (err error) {
			manager.Shutdown()
			return nil
		},
	})

	// script bind
	scripts.PushStruct("entityManager", NewEntityManagerBind(manager))

	return manager
}

// SetPluginManager ...
func (e *entityManager) SetPluginManager(pluginManager common.PluginManager) {
	e.pluginManager = pluginManager

	// event subscribe
	_ = e.eventBus.Subscribe(event_bus.TopicEntities, e.eventHandler)
	_ = e.eventBus.Subscribe(event_bus.TopicPlugins, e.eventHandler)
}

// LoadEntities ...
func (e *entityManager) LoadEntities() {

	var page int64
	var entities []*m.Entity
	const perPage = 100
	var err error

LOOP:
	entities, _, err = e.adaptors.Entity.List(perPage, perPage*page, "", "", true)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// add entities from database
	for _, entity := range entities {
		if err := e.Add(entity); err != nil {
			log.Warnf("%s, %s", entity.Id, err.Error())
		}
	}

	if len(entities) != 0 {
		page++
		goto LOOP
	}
}

// Shutdown ...
func (e *entityManager) Shutdown() {

	_ = e.eventBus.Unsubscribe(event_bus.TopicEntities, e.eventHandler)
	_ = e.eventBus.Unsubscribe(event_bus.TopicPlugins, e.eventHandler)

	e.actors.Range(func(key, value interface{}) bool {
		actor := value.(*actorInfo)
		actor.quit <- struct{}{}
		e.actors.Delete(key)
		return true
	})

	log.Info("Shutdown")
}

// SetMetric ...
func (e *entityManager) SetMetric(id common.EntityId, name string, value map[string]float32) {

	item, ok := e.actors.Load(id)
	if !ok {
		return
	}
	actor := item.(*actorInfo)

	var err error
	for _, metric := range actor.Actor.Metrics() {
		if metric.Name != name {
			continue
		}

		err = e.adaptors.MetricBucket.Add(m.MetricDataItem{
			Value:    value,
			MetricId: metric.Id,
			Time:     time.Now(),
		})

		if err != nil {
			log.Errorf(err.Error())
		}
	}
}

// SetState ...
func (e *entityManager) SetState(id common.EntityId, params EntityStateParams) (err error) {

	item, ok := e.actors.Load(id)
	if !ok {
		err = common.ErrNotFound
		return
	}
	actor := item.(*actorInfo)

	// store old state
	currentState := GetEventState(actor.Actor)
	actor.CurrentState = &currentState

	err = actor.Actor.SetState(params)

	return
}

// GetEntityById ...
func (e *entityManager) GetEntityById(id common.EntityId) (entity m.EntityShort, err error) {

	item, ok := e.actors.Load(id)
	if !ok {
		err = common.ErrNotFound
		return
	}
	actor := item.(*actorInfo)
	entity = NewEntity(actor.Actor)
	return
}

// GetActorById ...
func (e *entityManager) GetActorById(id common.EntityId) (actor PluginActor, err error) {

	item, ok := e.actors.Load(id)
	if !ok {
		err = common.ErrNotFound
		return
	}
	actor = item.(*actorInfo).Actor
	return
}

// List ...
func (e *entityManager) List() (entities []m.EntityShort, err error) {

	// sort index
	var index = make([]string, 0)
	e.actors.Range(func(key, value interface{}) bool {
		actor := value.(*actorInfo)
		info := actor.Actor.Info()
		index = append(index, info.Id.String())
		return true
	})
	sort.Strings(index)

	entities = make([]m.EntityShort, 0, len(index))
	var i int
	for _, n := range index {

		item, ok := e.actors.Load(n)
		if !ok {
			continue
		}
		actor := item.(*actorInfo)
		entities = append(entities, NewEntity(actor.Actor))

		// metric preview
		if len(entities[i].Metrics) > 0 {

			for j, metric := range entities[i].Metrics {
				var optionItems = make([]string, len(metric.Options.Items))
				for i, item := range metric.Options.Items {
					optionItems[i] = item.Name
				}

				if entities[i].Metrics[j].Data, err = e.adaptors.MetricBucket.Simple24HPreview(metric.Id, optionItems); err != nil {
					log.Error(err.Error())
					return
				}
			}
		}
		i++
	}
	return
}

// Spawn ...
func (e *entityManager) Spawn(constructor ActorConstructor) (actor PluginActor) {

	actor = constructor()
	info := actor.Info()

	defer func(entityId common.EntityId) {
		log.Infof("loaded %v", entityId)
	}(info.Id)

	var entityId = info.Id

	item, ok := e.actors.Load(entityId)
	if ok {
		log.Warnf("entityId %v exist", entityId)
		actor = item.(PluginActor)
		return
	}

	currentState := GetEventState(actor)
	actorInfo := &actorInfo{
		Actor:        actor,
		quit:         make(chan struct{}),
		CurrentState: &currentState,
	}
	e.actors.Store(entityId, actorInfo)

	//e.metric.Update(metrics.EntityAdd{Num: 1})

	go func() {
		defer func() {

			log.Infof("unload %v", entityId)

			e.eventBus.Publish(event_bus.TopicEntities, events.EventRemoveActor{
				PluginName: info.PluginName,
				EntityId:   entityId,
			})

			var err error
			var plugin CrudActor
			if plugin, err = e.getCrudActor(entityId); err != nil {
				return
			}
			_ = plugin.RemoveActor(entityId)

			//e.metric.Update(metrics.EntityDelete{Num: 1})
		}()

		<-actorInfo.quit
	}()

	attr := actor.Attributes()
	settings := actor.Settings()

	e.eventBus.Publish(event_bus.TopicEntities, events.EventAddedActor{
		PluginName: info.PluginName,
		EntityId:   entityId,
		Attributes: attr,
		Settings:   settings,
	})

	_ = e.adaptors.Entity.Add(&m.Entity{
		Id:          entityId,
		Description: info.Description,
		PluginName:  info.PluginName,
		Icon:        info.Icon,
		Area:        info.Area,
		Hidden:      info.Hidde,
		AutoLoad:    info.AutoLoad,
		ParentId:    info.ParentId,
		Attributes:  attr.Signature(),
		Settings:    settings,
	})

	return
}

// eventHandler ...
func (e *entityManager) eventHandler(_ string, message interface{}) {

	switch msg := message.(type) {
	case events.EventStateChanged:
		go e.eventStateChangedHandler(msg)
	case events.EventLoadedPlugin:
		go func() { _ = e.eventLoadedPlugin(msg) }()
	case events.EventUnloadedPlugin:
		go e.eventUnloadedPlugin(msg)
	case events.EventCreatedEntity:
		go e.eventCreatedEntity(msg)
	case events.EventUpdatedEntity:
		go e.eventUpdatedEntity(msg)
	case events.EventDeletedEntity:
		go e.eventDeletedEntity(msg)
	case events.EventEntitySetState:
		go e.eventEntitySetState(msg)
	case events.EventGetLastState:
		go e.eventLastState(msg)
	}
}

func (e *entityManager) eventStateChangedHandler(msg events.EventStateChanged) {

	item, ok := e.actors.Load(msg.EntityId)
	if !ok {
		return
	}
	actor := item.(*actorInfo)

	if msg.NewState.Compare(msg.OldState) {
		return
	}

	if actor.CurrentState != nil {
		if actor.CurrentState.Compare(msg.NewState) {
			return
		}
	}

	actor.CurrentState = &msg.NewState

	// store state to db
	var state string
	if msg.NewState.State != nil {
		state = msg.NewState.State.Name
	}

	if !msg.StorageSave {
		return
	}

	go func() {
		_, err := e.adaptors.EntityStorage.Add(&m.EntityStorage{
			State:      state,
			EntityId:   msg.EntityId,
			Attributes: msg.NewState.Attributes.Serialize(),
		})
		if err != nil {
			log.Error(err.Error())
		}
	}()
}

func (e *entityManager) eventLastState(msg events.EventGetLastState) {

	item, ok := e.actors.Load(msg.EntityId)
	if !ok {
		return
	}
	actor := item.(*actorInfo)

	if actor.CurrentState == nil {
		currentState := GetEventState(actor.Actor)
		actor.CurrentState = &currentState
	}

	info := actor.Actor.Info()

	e.eventBus.Publish(event_bus.TopicEntities, events.EventStateChanged{
		StorageSave: false,
		PluginName:  info.PluginName,
		EntityId:    info.Id,
		OldState:    *actor.CurrentState,
		NewState:    *actor.CurrentState,
	})
}

func (e *entityManager) eventLoadedPlugin(msg events.EventLoadedPlugin) (err error) {

	log.Infof("Load plugin '%s' entities", msg.PluginName)

	var entities []*m.Entity
	if entities, err = e.adaptors.Entity.GetByType(msg.PluginName, 1000, 0); err != nil {
		log.Error(err.Error())
		return
	}

	for _, entity := range entities {
		if err := e.Add(entity); err != nil {
			log.Warnf("%s, %s", entity.Id, err.Error())
		}
	}
	return
}

func (e *entityManager) eventUnloadedPlugin(msg events.EventUnloadedPlugin) {

	log.Infof("Unload plugin '%s' entities", msg.PluginName)

	e.actors.Range(func(key, value interface{}) bool {
		entityId := key.(common.EntityId)
		if entityId.PluginName() != msg.PluginName {
			return true
		}
		e.unsafeRemove(entityId)
		return true
	})
}

func (e *entityManager) eventCreatedEntity(msg events.EventCreatedEntity) {

	entity, err := e.adaptors.Entity.GetById(msg.Id)
	if err != nil {
		return
	}

	if err = e.Add(entity); err != nil {
		log.Error(err.Error())
	}
}

func (e *entityManager) eventUpdatedEntity(msg events.EventUpdatedEntity) {

	entity, err := e.adaptors.Entity.GetById(msg.Id)
	if err != nil {
		return
	}

	if err = e.Update(entity); err != nil {
		log.Error(err.Error())
	}
}

func (e *entityManager) eventDeletedEntity(msg events.EventDeletedEntity) {

	e.Remove(msg.Id)
}

func (e *entityManager) eventEntitySetState(msg events.EventEntitySetState) {

	_ = e.SetState(msg.Id, EntityStateParams{
		NewState:        msg.NewState,
		AttributeValues: msg.AttributeValues,
		SettingsValue:   msg.SettingsValue,
		StorageSave:     msg.StorageSave,
	})
}

// CallAction ...
func (e *entityManager) CallAction(id common.EntityId, action string, arg map[string]interface{}) {
	e.eventBus.Publish(event_bus.TopicEntities, events.EventCallAction{
		PluginName: id.PluginName(),
		EntityId:   id,
		ActionName: action,
		Args:       arg,
	})
}

// CallScene ...
func (e *entityManager) CallScene(id common.EntityId, arg map[string]interface{}) {
	e.eventBus.Publish(event_bus.TopicEntities, events.EventCallScene{
		PluginName: id.PluginName(),
		EntityId:   id,
		Args:       arg,
	})
}

func (e *entityManager) getCrudActor(entityId common.EntityId) (result CrudActor, err error) {
	var plugin interface{}
	if plugin, err = e.pluginManager.GetPlugin(entityId.PluginName()); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	var ok bool
	if result, ok = plugin.(CrudActor); ok {
		return
		//...
	} else {
		err = errors.Wrap(common.ErrInternal, fmt.Sprintf("can`t static cast '%s' to plugins.CrudActor", entityId.PluginName()))
	}
	return
}

// Add ...
func (e *entityManager) Add(entity *m.Entity) (err error) {

	var plugin m.Plugin
	if plugin, err = e.adaptors.Plugin.GetByName(entity.PluginName); err != nil {
		return
	}

	if !plugin.Enabled {
		return
	}

	var creudActor CrudActor
	if creudActor, err = e.getCrudActor(entity.Id); err != nil {
		return
	}

	err = creudActor.AddOrUpdateActor(entity)

	return
}

// Update ...
func (e *entityManager) Update(entity *m.Entity) (err error) {

	e.unsafeRemove(entity.Id)

	//todo fix
	time.Sleep(time.Millisecond * 1000)

	_ = e.Add(entity)

	return
}

// Remove ...
func (e *entityManager) Remove(id common.EntityId) {

	e.unsafeRemove(id)
}

func (e *entityManager) unsafeRemove(id common.EntityId) {

	item, ok := e.actors.Load(id)
	if !ok {
		return
	}
	actor := item.(*actorInfo)
	actor.quit <- struct{}{}
	e.actors.Delete(id)
}

// GetEventState ...
func GetEventState(actor PluginActor) (eventState event_bus.EventEntityState) {

	attrs := actor.Attributes()
	setts := actor.Settings()

	var state *event_bus.EntityState

	info := actor.Info()
	if info.State != nil {
		state = &event_bus.EntityState{
			Name:        info.State.Name,
			Description: info.State.Description,
			ImageUrl:    info.State.ImageUrl,
			Icon:        info.State.Icon,
		}
	}

	eventState = event_bus.EventEntityState{
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
