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
	"context"
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugin_manager"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/fx"
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
	pluginManager plugin_manager.PluginManager
	lock          *sync.Mutex
	actors        map[common.EntityId]*actorInfo
	quit          chan struct{}
}

func NewEntityManager(lc fx.Lifecycle,
	eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors,
	scripts scripts.ScriptService,
	pluginManager plugin_manager.PluginManager) EntityManager {
	manager := &entityManager{
		eventBus:      eventBus,
		adaptors:      adaptors,
		scripts:       scripts,
		pluginManager: pluginManager,
		lock:          &sync.Mutex{},
		actors:        make(map[common.EntityId]*actorInfo),
		quit:          make(chan struct{}),
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
func (e *entityManager) LoadEntities() {

	var page int64
	var entities []*m.Entity
	const perPage = 1000
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

	return
}

// Shutdown ...
func (e *entityManager) Shutdown() {
	log.Info("Shutdown")

	e.lock.Lock()
	defer e.lock.Unlock()

	for pid, actorInfo := range e.actors {
		close(actorInfo.Queue)
		delete(e.actors, pid)
	}
	return
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
func (e *entityManager) SetState(id common.EntityId, params EntityStateParams) {
	e.lock.Lock()
	defer e.lock.Unlock()

	actorInfo, ok := e.actors[id]
	if !ok {
		return
	}

	// store old state
	actorInfo.OldState = GetEventState(actorInfo.Actor)

	actorInfo.Actor.SetState(params)

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

	entities = make([]m.EntityShort, len(e.actors))
	var i int
	for _, actorInfo := range e.actors {
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

	actor = constructor(e)
	info := actor.Info()

	var entityId = info.Id

	if _, ok := e.actors[entityId]; ok {
		log.Warnf("entityId %v exist", entityId)
		return
	}

	// todo fix
	queue := make(chan Message, 99)

	e.actors[entityId] = &actorInfo{
		Actor:    actor,
		Queue:    queue,
		OldState: GetEventState(actor),
	}

	log.Infof("Loaded %v", entityId)

	//e.metric.Update(metrics.EntityAdd{Num: 1})

	go func() {
		defer func() {

			log.Infof("Unload %v", entityId)

			e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventRemoveEntity{
				Type:     info.Type,
				EntityId: entityId,
			})

			var err error
			var plugin plugin_manager.CrudActor
			if plugin, err = e.getCrudActor(entityId); err != nil {
				return
			}
			err = plugin.RemoveActor(entityId)

			//e.metric.Update(metrics.EntityDelete{Num: 1})
		}()

		for msg := range queue {
			actor.Receive(msg)
		}
	}()

	attr := actor.Attributes()

	e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventAddedNewEntity{
		Type:       info.Type,
		EntityId:   entityId,
		Attributes: attr,
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
	})

	return
}

// Send ...
func (e *entityManager) Send(message Message) error {

	switch v := message.Payload.(type) {
	case MessageRequestState:

		e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventRequestState{
			From:       message.From,
			To:         message.To,
			Attributes: v.Attributes,
		})

	case MessageStateChanged:

		e.lock.Lock()
		defer e.lock.Unlock()

		actorInfo, ok := e.actors[message.From]
		if !ok {
			return nil
		}

		if v.NewState.Compare(v.OldState) {
			return nil
		}

		e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventStateChanged{
			Type:     message.From.Type(),
			EntityId: message.From,
			OldState: actorInfo.OldState,
			NewState: v.NewState,
		})

		e.actors[message.From].OldState = v.NewState

		// store state to db
		var state string
		if v.NewState.State != nil {
			state = v.NewState.State.Name
		}

		if !v.StorageSave {
			return nil
		}
		go e.adaptors.EntityStorage.Add(m.EntityStorage{
			State:      state,
			EntityId:   message.From,
			Attributes: v.NewState.Attributes.Serialize(),
		})

		return nil

	case MessageCallAction:

		e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventCallAction{
			Type:       message.To.Type(),
			EntityId:   message.To,
			ActionName: v.Name,
			Args:       v.Arg,
		})

	case MessageCallScene:

		e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventCallScene{
			Type:     message.To.Type(),
			EntityId: message.To,
			Args:     v.Arg,
		})
	}

	if message.To == "" {
		return nil
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	if actorInfo, ok := e.actors[message.To]; ok {
		actorInfo.Queue <- message
	}
	return nil
}

// Broadcast ...
func (e *entityManager) Broadcast(message Message) {
	e.lock.Lock()
	defer e.lock.Unlock()

	for _, actorInfo := range e.actors {
		actorInfo.Queue <- message
	}
}

// CallAction ...
func (e *entityManager) CallAction(id common.EntityId, action string, arg map[string]interface{}) {

	go e.Send(Message{
		To: id,
		Payload: MessageCallAction{
			Name: action,
			Arg:  arg,
		},
	})
}

// CallScene ...
func (e *entityManager) CallScene(id common.EntityId, arg map[string]interface{}) {

	go e.Send(Message{
		To: id,
		Payload: MessageCallScene{
			Arg: arg,
		},
	})
}

func (e *entityManager) getCrudActor(entityId common.EntityId) (result plugin_manager.CrudActor, err error) {
	var plugin plugin_manager.Plugable
	if plugin, err = e.pluginManager.GetPlugin(entityId.Type().String()); err != nil {
		return
	}

	var ok bool
	if result, ok = plugin.(plugin_manager.CrudActor); ok {
		return
		//...
	} else {
		err = fmt.Errorf("cannot cast to the desired type plugin '%s' to plugin_manager.CrudActor", plugin.Name())
	}
	return
}

// Add ...
func (e *entityManager) Add(entity *m.Entity) (err error) {

	var plugin plugin_manager.CrudActor
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

	if actorInfo, ok := e.actors[id]; ok {
		close(actorInfo.Queue)
	} else {
		return
	}

	delete(e.actors, id)
}

// GetEventState ...
func GetEventState(actor PluginActor) (eventState event_bus.EventEntityState) {

	attrs := actor.Attributes()

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
	}

	if info.LastChanged != nil {
		eventState.LastChanged = common.Time(*info.LastChanged)
	}

	if info.LastUpdated != nil {
		eventState.LastUpdated = common.Time(*info.LastUpdated)
	}

	return
}
