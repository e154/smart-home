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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/fx"
	"sort"
	"sync"
	"time"
)

var (
	log = common.MustGetLogger("entity.manager")
)

type entityManager struct {
	eventBus      event_bus.EventBus
	adaptors      *adaptors.Adaptors
	scripts       scripts.ScriptService
	pluginManager common.PluginManager
	lock          *sync.Mutex
	actors        map[common.EntityId]*actorInfo
	quit          chan struct{}
}

func NewEntityManager(lc fx.Lifecycle,
	eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors,
	scripts scripts.ScriptService) EntityManager {
	manager := &entityManager{
		eventBus: eventBus,
		adaptors: adaptors,
		scripts:  scripts,
		lock:     &sync.Mutex{},
		actors:   make(map[common.EntityId]*actorInfo),
		quit:     make(chan struct{}),
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) (err error) {
			manager.Shutdown()
			return nil
		},
	})

	return manager
}

// LoadEntities ...
func (e *entityManager) LoadEntities(pluginManager common.PluginManager) {

	e.pluginManager = pluginManager

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
			log.Warn(err.Error())
		}
	}

	if len(entities) != 0 {
		page++
		goto LOOP
	}

	// scripts
	e.scripts.PushStruct("entityManager", NewEntityManagerBind(e))

	// event subscribe
	e.eventBus.Subscribe(event_bus.TopicEntities, e.eventHandler)

	return
}

// Shutdown ...
func (e *entityManager) Shutdown() {

	e.lock.Lock()
	defer e.lock.Unlock()

	e.eventBus.Unsubscribe(event_bus.TopicEntities, e.eventHandler)

	for id, actor := range e.actors {
		actor.quit <- struct{}{}
		delete(e.actors, id)
	}

	log.Info("Shutdown")
}

// SetMetric ...
func (e *entityManager) SetMetric(id common.EntityId, name string, value map[string]interface{}) {

	e.lock.Lock()
	defer e.lock.Unlock()

	actorInfo, ok := e.actors[id]
	if !ok {
		return
	}

	var err error
	for _, metric := range actorInfo.Actor.Metrics() {
		if metric.Name != name {
			continue
		}

		var b []byte
		if b, err = json.Marshal(value); err != nil {
			log.Error(err.Error(), "value", value)
			return
		}

		err = e.adaptors.MetricBucket.Add(m.MetricDataItem{
			Value:    b,
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
	e.lock.Lock()
	defer e.lock.Unlock()

	actorInfo, ok := e.actors[id]
	if !ok {
		err = errors.New("not found")
		return
	}

	// store old state
	actorInfo.OldState = GetEventState(actorInfo.Actor)

	err = actorInfo.Actor.SetState(params)

	return
}

// GetEntityById ...
func (e *entityManager) GetEntityById(id common.EntityId) (entity m.EntityShort, err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if actorInfo, ok := e.actors[id]; ok {
		entity = NewEntity(actorInfo.Actor)
	} else {
		err = fmt.Errorf("actor '%v' not found", id)
	}
	return
}

// GetActorById ...
func (e *entityManager) GetActorById(id common.EntityId) (actor PluginActor, err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if actorInfo, ok := e.actors[id]; ok {
		actor = actorInfo.Actor
	} else {
		err = fmt.Errorf("actor '%v' not found", id)
	}
	return
}

// List ...
func (e *entityManager) List() (entities []m.EntityShort, err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	// sort index
	var index = make([]string, 0, len(e.actors))
	for _, actor := range e.actors {
		info := actor.Actor.Info()
		index = append(index, info.Id.String())
	}
	sort.Strings(index)

	entities = make([]m.EntityShort, len(e.actors))
	var i int
	for _, n := range index {

		actorInfo := e.actors[common.EntityId(n)]

		entities[i] = NewEntity(actorInfo.Actor)

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

	e.lock.Lock()
	defer e.lock.Unlock()

	actor = constructor()
	info := actor.Info()

	defer func(entityId common.EntityId) {
		log.Infof("loaded %v", entityId)
	}(info.Id)

	var entityId = info.Id

	if _, ok := e.actors[entityId]; ok {
		log.Warnf("entityId %v exist", entityId)
		return
	}

	e.actors[entityId] = &actorInfo{
		Actor:    actor,
		quit:     make(chan struct{}),
		OldState: GetEventState(actor),
	}

	//e.metric.Update(metrics.EntityAdd{Num: 1})

	go func() {
		defer func() {

			log.Infof("unload %v", entityId)

			e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventRemoveActor{
				Type:     info.Type,
				EntityId: entityId,
			})

			var err error
			var plugin CrudActor
			if plugin, err = e.getCrudActor(entityId); err != nil {
				return
			}
			err = plugin.RemoveActor(entityId)

			//e.metric.Update(metrics.EntityDelete{Num: 1})
		}()

		<-e.actors[entityId].quit
	}()

	attr := actor.Attributes()
	settings := actor.Settings()

	e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventAddedActor{
		Type:       info.Type,
		EntityId:   entityId,
		Attributes: attr,
		Settings:   settings,
	})

	e.adaptors.Entity.Add(&m.Entity{
		Id:          entityId,
		Description: info.Description,
		Type:        info.Type,
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

	switch v := message.(type) {
	case event_bus.EventStateChanged:

		e.lock.Lock()
		defer e.lock.Unlock()

		if _, ok := e.actors[v.EntityId]; !ok {
			return
		}

		if v.NewState.Compare(v.OldState) {
			return
		}

		e.actors[v.EntityId].OldState = v.NewState

		// store state to db
		var state string
		if v.NewState.State != nil {
			state = v.NewState.State.Name
		}

		if !v.StorageSave {
			return
		}

		go e.adaptors.EntityStorage.Add(m.EntityStorage{
			State:      state,
			EntityId:   v.EntityId,
			Attributes: v.NewState.Attributes.Serialize(),
		})

	}
}

// CallAction ...
func (e *entityManager) CallAction(id common.EntityId, action string, arg map[string]interface{}) {
	e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventCallAction{
		Type:       id.Type(),
		EntityId:   id,
		ActionName: action,
		Args:       arg,
	})
}

// CallScene ...
func (e *entityManager) CallScene(id common.EntityId, arg map[string]interface{}) {
	e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventCallScene{
		Type:     id.Type(),
		EntityId: id,
		Args:     arg,
	})
}

func (e *entityManager) getCrudActor(entityId common.EntityId) (result CrudActor, err error) {
	var plugin interface{}
	if plugin, err = e.pluginManager.GetPlugin(entityId.Type().String()); err != nil {
		return
	}

	var ok bool
	if result, ok = plugin.(CrudActor); ok {
		return
		//...
	} else {
		err = fmt.Errorf("cannot cast to the desired type plugin '%s' to plugins.CrudActor", entityId.Type().String())
	}
	return
}

// Add ...
func (e *entityManager) Add(entity *m.Entity) (err error) {

	var plugin CrudActor
	if plugin, err = e.getCrudActor(entity.Id); err != nil {
		return
	}

	err = plugin.AddOrUpdateActor(entity)

	return
}

// Update ...
func (e *entityManager) Update(entity *m.Entity) (err error) {

	e.unsafeRemove(entity.Id)

	//todo fix
	time.Sleep(time.Millisecond * 500)

	e.Add(entity)

	return
}

// Remove ...
func (e *entityManager) Remove(id common.EntityId) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.unsafeRemove(id)
}

func (e *entityManager) unsafeRemove(id common.EntityId) {

	if actor, ok := e.actors[id]; ok {
		actor.quit <- struct{}{}
	} else {
		return
	}

	delete(e.actors, id)
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
