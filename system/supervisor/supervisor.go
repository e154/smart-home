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

package supervisor

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"go.uber.org/fx"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/common/web"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
)

var (
	log = logger.MustGetLogger("supervisor")
)

type supervisor struct {
	*pluginManager
	scriptService scripts.ScriptService
	actors        sync.Map
	quit          chan struct{}
}

// NewSupervisor ...
func NewSupervisor(lc fx.Lifecycle,
	adaptors *adaptors.Adaptors,
	bus bus.Bus,
	mqttServ mqtt.MqttServ,
	scriptService scripts.ScriptService,
	appConfig *m.AppConfig,
	gateClient *gate_client.GateClient,
	eventBus bus.Bus,
	scheduler *scheduler.Scheduler,
	crawler web.Crawler) Supervisor {
	s := &supervisor{
		scriptService: scriptService,
		actors:        sync.Map{},
		quit:          make(chan struct{}),
	}
	s.pluginManager = &pluginManager{
		adaptors:       adaptors,
		isStarted:      atomic.NewBool(false),
		eventBus:       eventBus,
		enabledPlugins: make(map[string]bool),
		service: &service{
			bus:           bus,
			supervisor:    s,
			mqttServ:      mqttServ,
			adaptors:      adaptors,
			scriptService: scriptService,
			appConfig:     appConfig,
			gateClient:    gateClient,
			scheduler:     scheduler,
			crawler:       crawler,
		},
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})

	return s
}

func (e *supervisor) Start(ctx context.Context) (err error) {

	// event subscribe
	_ = e.eventBus.Subscribe("system/entities/+", e.eventHandler)
	_ = e.eventBus.Subscribe("system/plugins/+", e.eventHandler)

	e.scriptService.PushStruct("supervisor", NewSupervisorBind(e))
	//DEPRECATED
	e.scriptService.PushStruct("entityManager", NewSupervisorBind(e))

	e.pluginManager.Start(ctx)

	_ = e.eventBus.Subscribe("system/services/scripts", e.handlerSystemScripts)
	e.eventBus.Publish("system/services/supervisor", events.EventServiceStarted{Service: "Supervisor"})

	log.Info("Started")

	return
}

// Shutdown ...
func (e *supervisor) Shutdown(ctx context.Context) (err error) {

	_ = e.eventBus.Unsubscribe("system/services/scripts", e.handlerSystemScripts)

	e.scriptService.PopStruct("supervisor")
	e.scriptService.PopStruct("entityManager")

	e.pluginManager.Shutdown(ctx)

	_ = e.eventBus.Unsubscribe("system/entities/+", e.eventHandler)
	_ = e.eventBus.Unsubscribe("system/plugins/+", e.eventHandler)

	e.actors.Range(func(key, value interface{}) bool {
		actor := value.(*actorInfo)
		close(actor.quit)
		e.actors.Delete(key)
		return true
	})

	e.eventBus.Publish("system/services/supervisor", events.EventServiceStopped{Service: "Supervisor"})

	log.Info("Shutdown")

	return
}

// Restart ...
func (e *supervisor) Restart(ctx context.Context) (err error) {
	if err = e.Shutdown(ctx); err != nil {
		return
	}
	err = e.Start(ctx)
	return
}

func (e *supervisor) handlerSystemScripts(_ string, event interface{}) {

	switch event.(type) {
	case events.EventServiceStarted, events.EventServiceRestarted:
		e.scriptService.PushStruct("supervisor", NewSupervisorBind(e))
		//DEPRECATED
		e.scriptService.PushStruct("entityManager", NewSupervisorBind(e))
	}
}

func (e *supervisor) LoadEntities() {

	var page int64
	var entities []*m.Entity
	const perPage = 500
	var err error

LOOP:
	entities, _, err = e.adaptors.Entity.List(context.Background(), perPage, perPage*page, "", "", true, nil, nil, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// add entities from database
	for _, entity := range entities {
		if err = e.AddEntity(entity); err != nil {
			log.Warnf("%s, %s", entity.Id, err.Error())
		}
	}

	if len(entities) != 0 {
		page++
		goto LOOP
	}
}

func (e *supervisor) updateMetric(actor *actorInfo, state bus.EventEntityState) {
	metrics := actor.Actor.Metrics()
	if metrics == nil {
		return
	}

	var data = make(map[string]float32)
	var name string

	for _, metric := range metrics {
		for _, prop := range metric.Options.Items {
			if value, ok := state.Attributes[prop.Name]; ok {
				name = metric.Name
				switch value.Type {
				case common.AttributeInt:
					data[prop.Name] = float32(value.Int64())
				case common.AttributeFloat:
					data[prop.Name] = common.Rounding32(value.Float64(), 2)
				}
			}
		}
	}

	if len(data) == 0 || name == "" {
		return
	}

	e.SetMetric(state.EntityId, name, data)

}

// SetMetric ...
func (e *supervisor) SetMetric(id common.EntityId, name string, value map[string]float32) {

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

		err = e.adaptors.MetricBucket.Add(context.Background(), &m.MetricDataItem{
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
func (e *supervisor) SetState(id common.EntityId, params EntityStateParams) (err error) {

	item, ok := e.actors.Load(id)
	if !ok {
		err = apperr.ErrNotFound
		return
	}
	actor := item.(*actorInfo)

	// store old state
	currentState := GetEventState(actor.Actor)
	actor.CurrentState = &currentState

	err = actor.Actor.SetState(params)

	return
}

// EntityIsLoaded ...
func (e *supervisor) EntityIsLoaded(id common.EntityId) (loaded bool) {
	_, loaded = e.actors.Load(id)
	return
}

// GetEntityById ...
func (e *supervisor) GetEntityById(id common.EntityId) (entity m.EntityShort, err error) {

	item, ok := e.actors.Load(id)
	if !ok {
		err = apperr.ErrNotFound
		return
	}
	actor := item.(*actorInfo)
	entity = NewEntity(actor.Actor)
	return
}

// GetActorById ...
func (e *supervisor) GetActorById(id common.EntityId) (actor PluginActor, err error) {

	item, ok := e.actors.Load(id)
	if !ok {
		err = apperr.ErrNotFound
		return
	}
	actor = item.(*actorInfo).Actor
	return
}

// List ...
func (e *supervisor) List() (entities []m.EntityShort, err error) {

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

				if entities[i].Metrics[j].Data, err = e.adaptors.MetricBucket.Simple24HPreview(context.Background(), metric.Id, optionItems); err != nil {
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
func (e *supervisor) Spawn(constructor ActorConstructor) (actor PluginActor) {

	actor = constructor()
	info := actor.Info()

	defer func(entityId common.EntityId) {
		log.Infof("loaded entity '%v'", entityId)
	}(info.Id)

	var entityId = info.Id

	item, ok := e.actors.Load(entityId)
	if ok && item != nil {
		log.Warnf("entityId '%v' exist", entityId)
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

	go func() {
		defer func() {

			log.Infof("unload entity %v", entityId)

			var err error
			var plugin CrudActor
			if plugin, err = e.getCrudActor(entityId); err != nil {
				return
			}
			_ = plugin.RemoveActor(entityId)

			e.eventBus.Publish("system/entities/"+entityId.String(), events.EventEntityUnloaded{
				PluginName: info.PluginName,
				EntityId:   entityId,
			})
		}()

		<-actorInfo.quit
	}()

	attr := actor.Attributes()
	settings := actor.Settings()

	e.eventBus.Publish("system/entities/"+entityId.String(), events.EventAddedActor{
		PluginName: info.PluginName,
		EntityId:   entityId,
		Attributes: attr,
		Settings:   settings,
	})

	_ = e.adaptors.Entity.Add(context.Background(), &m.Entity{
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
func (e *supervisor) eventHandler(_ string, message interface{}) {

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
	case events.CommandUnloadEntity:
		go e.eventUnloadEntity(msg)
	case events.CommandLoadEntity:
		go e.eventLoadEntity(msg)
	case events.EventEntitySetState:
		go e.eventEntitySetState(msg)
	case events.EventGetLastState:
		go e.eventLastState(msg)
	}
}

func (e *supervisor) eventStateChangedHandler(msg events.EventStateChanged) {

	item, ok := e.actors.Load(msg.EntityId)
	if !ok {
		return
	}
	actor := item.(*actorInfo)

	go e.updateMetric(actor, msg.NewState)

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
		_, err := e.adaptors.EntityStorage.Add(context.Background(), &m.EntityStorage{
			State:      state,
			EntityId:   msg.EntityId,
			Attributes: msg.NewState.Attributes.Serialize(),
		})
		if err != nil {
			log.Error(err.Error())
		}
	}()
}

func (e *supervisor) eventLastState(msg events.EventGetLastState) {

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

	if actor.CurrentState.LastChanged == nil && actor.CurrentState.LastUpdated == nil {
		entity, _ := e.adaptors.Entity.GetById(context.Background(), msg.EntityId)
		actor.CurrentState.Attributes = entity.Attributes
	}

	e.eventBus.Publish("system/entities/"+msg.EntityId.String(), events.EventLastStateChanged{
		PluginName: info.PluginName,
		EntityId:   info.Id,
		OldState:   *actor.CurrentState,
		NewState:   *actor.CurrentState,
	})
}

func (e *supervisor) eventLoadedPlugin(msg events.EventLoadedPlugin) (err error) {

	log.Infof("Load plugin '%s' entities", msg.PluginName)

	var page int64
	var entities []*m.Entity
	const perPage = 500

LOOP:

	if entities, err = e.adaptors.Entity.GetByType(context.Background(), msg.PluginName, perPage, perPage*page); err != nil {
		log.Error(err.Error())
		return
	}

	for _, entity := range entities {
		if err := e.AddEntity(entity); err != nil {
			log.Warnf("%s, %s", entity.Id, err.Error())
		}
	}

	if len(entities) != 0 {
		page++
		goto LOOP
	}

	return
}

func (e *supervisor) eventUnloadedPlugin(msg events.EventUnloadedPlugin) {

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

func (e *supervisor) eventCreatedEntity(msg events.EventCreatedEntity) {

	entity, err := e.adaptors.Entity.GetById(context.Background(), msg.EntityId)
	if err != nil {
		return
	}

	if !entity.AutoLoad {
		return
	}

	if err = e.AddEntity(entity); err != nil {
		log.Error(err.Error())
	}
}

func (e *supervisor) eventUpdatedEntity(msg events.EventUpdatedEntity) {

	entity, err := e.adaptors.Entity.GetById(context.Background(), msg.EntityId)
	if err != nil {
		return
	}

	if err = e.UpdateEntity(entity); err != nil {
		log.Error(err.Error())
	}
}

func (e *supervisor) eventUnloadEntity(msg events.CommandUnloadEntity) {

	e.Remove(msg.EntityId)
}

func (e *supervisor) eventLoadEntity(msg events.CommandLoadEntity) {
	entity, _ := e.adaptors.Entity.GetById(context.Background(), msg.EntityId)
	if err := e.AddEntity(entity); err != nil {
		log.Warnf("%s, %s", entity.Id, err.Error())
	}
}

func (e *supervisor) eventEntitySetState(msg events.EventEntitySetState) {

	_ = e.SetState(msg.EntityId, EntityStateParams{
		NewState:        msg.NewState,
		AttributeValues: msg.AttributeValues,
		SettingsValue:   msg.SettingsValue,
		StorageSave:     msg.StorageSave,
	})
}

// CallAction ...
func (e *supervisor) CallAction(id common.EntityId, action string, arg map[string]interface{}) {
	e.eventBus.Publish("system/entities/"+id.String(), events.EventCallEntityAction{
		PluginName: id.PluginName(),
		EntityId:   id,
		ActionName: action,
		Args:       arg,
	})
}

// CallScene ...
func (e *supervisor) CallScene(id common.EntityId, arg map[string]interface{}) {
	e.eventBus.Publish("system/entities/"+id.String(), events.EventCallScene{
		PluginName: id.PluginName(),
		EntityId:   id,
		Args:       arg,
	})
}

func (e *supervisor) getCrudActor(entityId common.EntityId) (result CrudActor, err error) {
	var plugin interface{}
	if plugin, err = e.getPlugin(entityId.PluginName()); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	var ok bool
	if result, ok = plugin.(CrudActor); ok {
		return
		//...
	} else {
		err = errors.Wrap(apperr.ErrInternal, fmt.Sprintf("can`t static cast '%s' to plugins.CrudActor", entityId.PluginName()))
	}
	return
}

// AddEntity ...
func (e *supervisor) AddEntity(entity *m.Entity) (err error) {

	if _, ok := e.enabledPlugins[entity.PluginName]; !ok {
		return
	}

	var crudActor CrudActor
	if crudActor, err = e.getCrudActor(entity.Id); err != nil {
		return
	}

	if err = crudActor.AddOrUpdateActor(entity); err != nil {
		return
	}

	e.eventBus.Publish("system/entities/"+entity.Id.String(), events.EventEntityLoaded{
		EntityId: entity.Id,
	})

	return
}

// UpdateEntity ...
func (e *supervisor) UpdateEntity(entity *m.Entity) (err error) {

	e.unsafeRemove(entity.Id)

	//todo fix
	time.Sleep(time.Millisecond * 1000)

	_ = e.AddEntity(entity)

	return
}

// Remove ...
func (e *supervisor) Remove(id common.EntityId) {

	e.unsafeRemove(id)
}

func (e *supervisor) unsafeRemove(id common.EntityId) {

	item, ok := e.actors.Load(id)
	if !ok {
		return
	}
	actor := item.(*actorInfo)
	close(actor.quit)
	e.actors.Delete(id)
}

// GetEventState ...
func GetEventState(actor PluginActor) (eventState bus.EventEntityState) {

	attrs := actor.Attributes()
	setts := actor.Settings()

	var state *bus.EntityState

	info := actor.Info()
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
