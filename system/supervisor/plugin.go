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

	"github.com/pkg/errors"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
)

// Plugin ...
type Plugin struct {
	Service          Service
	IsStarted        *atomic.Bool
	Actors           *sync.Map
	actorConstructor ActorConstructor
}

// NewPlugin ...
func NewPlugin() *Plugin {
	return &Plugin{
		IsStarted: atomic.NewBool(false),
		Actors:    &sync.Map{},
	}
}

// Load ...
func (p *Plugin) Load(ctx context.Context, service Service, actorConstructor ActorConstructor) error {
	p.Service = service
	p.actorConstructor = actorConstructor

	if p.IsStarted.Load() {
		return ErrPluginIsLoaded
	}
	p.IsStarted.Store(true)

	return nil
}

// Unload ...
func (p *Plugin) Unload(ctx context.Context) error {

	if !p.IsStarted.Load() {
		return ErrPluginIsUnloaded
	}

	p.Actors.Range(func(key, value any) bool {
		if pla, ok := value.(PluginActor); ok {
			p.removePluginActor(pla)
		}
		return true
	})

	p.IsStarted.Store(false)

	return nil
}

// Name ...
func (p *Plugin) Name() string {
	panic("implement me")
}

// Type ...
func (p *Plugin) Type() PluginType {
	panic("implement me")
}

// Depends ...
func (p *Plugin) Depends() []string {
	panic("implement me")
}

// Version ...
func (p *Plugin) Version() string {
	panic("implement me")
}

// Options ...
func (p *Plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		ActorCustomAttrs: false,
		ActorAttrs:       nil,
		ActorSetts:       nil,
	}
}

// LoadSettings ...
func (p *Plugin) LoadSettings(pl Pluggable) (settings m.Attributes, err error) {
	var plugin *m.Plugin
	if plugin, err = p.Service.Adaptors().Plugin.GetByName(context.Background(), pl.Name()); err != nil {
		return
	}
	settings = pl.Options().Setts
	if settings == nil {
		settings = make(m.Attributes)
		return
	}
	_, err = settings.Deserialize(plugin.Settings)
	return
}

func (p *Plugin) AddOrUpdateActor(entity *m.Entity) (err error) {

	if p.actorConstructor == nil {
		return
	}

	var pla PluginActor
	if pla, err = p.actorConstructor(entity); err != nil {
		return
	}

	item, ok := p.Actors.Load(entity.Id)
	if ok && item != nil {
		_ = p.RemoveActor(entity.Id)
	}

	err = p.AddActor(pla, entity)

	return
}

func (p *Plugin) AddActor(pla PluginActor, entity *m.Entity) (err error) {

	if entity == nil {
		return
	}

	pla.Spawn()
	p.Actors.Store(entity.Id, pla)
	log.Infof("entity '%v' loaded", entity.Id)

	currentState := pla.GetEventState()
	pla.SetCurrentState(currentState)

	p.Service.EventBus().Publish("system/entities/"+entity.Id.String(), events.EventEntityLoaded{
		EntityId:   entity.Id,
		PluginName: entity.PluginName,
	})

	if _, err = p.Service.Adaptors().Entity.GetById(context.Background(), entity.Id); err == nil {
		return
	}

	err = p.Service.Adaptors().Entity.Add(context.Background(), &m.Entity{
		Id:          entity.Id,
		Description: entity.Description,
		PluginName:  entity.PluginName,
		Icon:        entity.Icon,
		Area:        entity.Area,
		Hidden:      entity.Hidden,
		AutoLoad:    entity.AutoLoad,
		ParentId:    entity.ParentId,
		Attributes:  entity.Attributes.Signature(),
		Settings:    entity.Settings,
	})

	return
}

func (p *Plugin) RemoveActor(entityId common.EntityId) (err error) {

	item, ok := p.Actors.Load(entityId)
	if !ok || item == nil {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("failed remove \"%s\"", entityId))
		return
	}

	pla := item.(PluginActor)
	p.removePluginActor(pla)
	return
}

func (p *Plugin) removePluginActor(pla PluginActor) {

	info := pla.Info()
	entityId := info.Id

	pla.Destroy()
	pla.StopWatchers()
	p.Actors.Delete(entityId)

	p.Service.EventBus().Publish("system/entities/"+entityId.String(), events.EventEntityUnloaded{
		PluginName: entityId.PluginName(),
		EntityId:   entityId,
	})

	log.Infof("entity '%v' unloaded", entityId)
}

func (p *Plugin) EntityIsLoaded(id common.EntityId) bool {
	value, ok := p.Actors.Load(id)
	return ok && value != nil
}

func (p *Plugin) GetActor(id common.EntityId) (pla PluginActor, err error) {
	value, ok := p.Actors.Load(id)
	if !ok || value == nil {
		err = errors.Wrap(apperr.ErrEntityNotFound, id.String())
		return
	}
	pla = value.(PluginActor)
	return
}
