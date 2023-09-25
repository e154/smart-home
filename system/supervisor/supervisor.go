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
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"sync"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
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
	quit          chan struct{}
	entitiesWg    *sync.WaitGroup
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
		quit:          make(chan struct{}),
		entitiesWg:    &sync.WaitGroup{},
	}
	s.pluginManager = &pluginManager{
		adaptors:       adaptors,
		isStarted:      atomic.NewBool(false),
		eventBus:       eventBus,
		enabledPlugins: sync.Map{},
		pluginsWg:      &sync.WaitGroup{},
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

	e.scriptService.PushFunctions("GetEntity", GetEntityBind(e))
	e.scriptService.PushFunctions("SetState", SetStateBind(e))
	e.scriptService.PushFunctions("SetStateName", SetStateNameBind(e))
	e.scriptService.PushFunctions("GetState", GetStateBind(e))
	e.scriptService.PushFunctions("SetAttributes", SetAttributesBind(e))
	e.scriptService.PushFunctions("GetAttributes", GetAttributesBind(e))
	e.scriptService.PushFunctions("GetSettings", GetSettingsBind(e))
	e.scriptService.PushFunctions("SetMetric", SetMetricBind(e))
	e.scriptService.PushFunctions("CallAction", CallActionBind(e))
	e.scriptService.PushFunctions("CallScene", CallSceneBind(e))

	e.pluginManager.Start(ctx)

	_ = e.eventBus.Subscribe("system/services/scripts", e.handlerSystemScripts)
	e.eventBus.Publish("system/services/supervisor", events.EventServiceStarted{Service: "Supervisor"})

	log.Info("Started")

	return
}

// Shutdown ...
func (e *supervisor) Shutdown(ctx context.Context) (err error) {

	e.pluginManager.Shutdown(ctx)
	e.entitiesWg.Wait()

	e.scriptService.PopFunction("GetEntity")
	e.scriptService.PopFunction("SetState")
	e.scriptService.PopFunction("SetStateName")
	e.scriptService.PopFunction("SetAttributes")
	e.scriptService.PopFunction("GetAttributes")
	e.scriptService.PopFunction("GetSettings")
	e.scriptService.PopFunction("SetMetric")
	e.scriptService.PopFunction("CallAction")
	e.scriptService.PopFunction("CallScene")

	_ = e.eventBus.Unsubscribe("system/services/scripts", e.handlerSystemScripts)
	_ = e.eventBus.Unsubscribe("system/entities/+", e.eventHandler)
	_ = e.eventBus.Unsubscribe("system/plugins/+", e.eventHandler)

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
		e.scriptService.PushFunctions("GetEntity", GetEntityBind(e))
		e.scriptService.PushFunctions("SetState", SetStateBind(e))
		e.scriptService.PushFunctions("SetStateName", SetStateBind(e))
		e.scriptService.PushFunctions("GetState", GetStateBind(e))
		e.scriptService.PushFunctions("SetAttributes", SetAttributesBind(e))
		e.scriptService.PushFunctions("GetAttributes", GetAttributesBind(e))
		e.scriptService.PushFunctions("GetSettings", GetSettingsBind(e))
		e.scriptService.PushFunctions("SetMetric", SetMetricBind(e))
		e.scriptService.PushFunctions("CallAction", CallActionBind(e))
		e.scriptService.PushFunctions("CallScene", CallSceneBind(e))
	}
}

// SetMetric ...
func (e *supervisor) SetMetric(id common.EntityId, name string, value map[string]float32) {

	pla, err := e.GetActorById(id)
	if err != nil {
		return
	}

	pla.AddMetric(name, value)
}

// SetState ...
func (e *supervisor) SetState(id common.EntityId, params EntityStateParams) (err error) {

	pla, err := e.GetActorById(id)
	if err != nil {
		return
	}

	err = pla.SetState(params)

	return
}

// EntityIsLoaded ...
func (e *supervisor) EntityIsLoaded(id common.EntityId) (loaded bool) {

	if !e.PluginIsLoaded(id.PluginName()) {
		return
	}

	value, err := e.GetPlugin(id.PluginName())
	if err != nil {
		return
	}

	plugin := value.(Pluggable)
	loaded = plugin.EntityIsLoaded(id)

	return
}

// GetEntityById ...
func (e *supervisor) GetEntityById(id common.EntityId) (entity m.EntityShort, err error) {

	var pla PluginActor
	if pla, err = e.GetActorById(id); err != nil {
		return
	}
	entity = NewEntity(pla)
	return
}

// GetActorById ...
func (e *supervisor) GetActorById(id common.EntityId) (pla PluginActor, err error) {

	if !e.PluginIsLoaded(id.PluginName()) {
		err = errors.Wrap(ErrPluginNotLoaded, id.PluginName())
		return
	}

	var value interface{}
	if value, err = e.GetPlugin(id.PluginName()); err != nil {
		return
	}
	plugin := value.(Pluggable)

	pla, err = plugin.GetActor(id)

	return
}

// eventHandler ...
func (e *supervisor) eventHandler(_ string, message interface{}) {

	switch msg := message.(type) {
	case events.EventLoadedPlugin:
		go func() { _ = e.eventLoadedPlugin(msg) }()
	case events.EventCreatedEntity:
		go e.eventCreatedEntity(msg)
	case events.EventUpdatedEntity:
		go e.eventUpdatedEntity(msg)
	case events.CommandUnloadEntity:
		go e.commandUnloadEntity(msg)
	case events.CommandLoadEntity:
		go e.commandLoadEntity(msg)
	case events.EventEntitySetState:
		go e.eventEntitySetState(msg)
	case events.EventGetLastState:
		go e.eventLastState(msg)
	}
}

func (e *supervisor) eventLastState(msg events.EventGetLastState) {

	pla, err := e.GetActorById(msg.EntityId)
	if err != nil {
		return
	}

	if pla.GetCurrentState() == nil {
		currentState := pla.GetEventState()
		pla.SetCurrentState(currentState)
	}

	info := pla.Info()

	currentState := pla.GetCurrentState()
	if currentState.LastChanged == nil && currentState.LastUpdated == nil {
		entity, _ := e.adaptors.Entity.GetById(context.Background(), msg.EntityId)
		currentState.Attributes = entity.Attributes
	}

	e.eventBus.Publish("system/entities/"+msg.EntityId.String(), events.EventLastStateChanged{
		PluginName: info.PluginName,
		EntityId:   info.Id,
		OldState:   *currentState,
		NewState:   *currentState,
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
		go func(entity *m.Entity) {
			if err = e.AddEntity(entity); err != nil {
				log.Warnf("%s, %s", entity.Id, err.Error())
			}
		}(entity)
	}

	if len(entities) != 0 {
		page++
		goto LOOP
	}

	return
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
	if err != nil || !entity.AutoLoad {
		return
	}

	if err = e.UpdateEntity(entity); err != nil {
		log.Error(err.Error())
	}
}

func (e *supervisor) commandUnloadEntity(msg events.CommandUnloadEntity) {
	e.UnloadEntity(msg.EntityId)
}

func (e *supervisor) commandLoadEntity(msg events.CommandLoadEntity) {
	entity, err := e.adaptors.Entity.GetById(context.Background(), msg.EntityId)
	if err != nil {
		return
	}

	if !entity.AutoLoad {
		return
	}

	if err = e.AddEntity(entity); err != nil {
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

// AddEntity ...
func (e *supervisor) AddEntity(entity *m.Entity) (err error) {

	if !e.PluginIsLoaded(entity.PluginName) {
		err = errors.Wrap(ErrPluginNotLoaded, entity.PluginName)
		return
	}

	var value interface{}
	if value, err = e.GetPlugin(entity.PluginName); err != nil {
		return
	}
	plugin := value.(Pluggable)
	err = plugin.AddOrUpdateActor(entity)
	return
}

// UpdateEntity ...
func (e *supervisor) UpdateEntity(entity *m.Entity) (err error) {

	if !e.PluginIsLoaded(entity.PluginName) {
		err = errors.Wrap(ErrPluginNotLoaded, entity.PluginName)
		return
	}

	var value interface{}
	if value, err = e.GetPlugin(entity.PluginName); err != nil {
		return
	}

	plugin := value.(Pluggable)

	err = plugin.AddOrUpdateActor(entity)

	return
}

// UnloadEntity ...
func (e *supervisor) UnloadEntity(id common.EntityId) {

	if !e.PluginIsLoaded(id.PluginName()) {
		return
	}

	value, err := e.GetPlugin(id.PluginName())
	if err != nil {
		return
	}

	plugin := value.(Pluggable)
	plugin.RemoveActor(id)
}

func (e *supervisor) GetService() Service {
	return e.service
}
