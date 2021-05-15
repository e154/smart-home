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

package alexa

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.server")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	entityManager entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	isStarted     *atomic.Bool
	eventBus      event_bus.EventBus
	server        IServer
	actorsLock    *sync.Mutex
	registrar     triggers.IRegistrar
}

func New() plugins.Plugable {
	return &plugin{
		isStarted:  atomic.NewBool(false),
		actorsLock: &sync.Mutex{},
	}
}

func (p *plugin) Load(service plugins.Service) error {
	p.adaptors = service.Adaptors()
	p.eventBus = service.EventBus()
	p.entityManager = service.EntityManager()
	p.scriptService = service.ScriptService()

	if p.isStarted.Load() {
		return nil
	}
	p.isStarted.Store(true)

	// register trigger
	if triggersPlugin, ok := service.Plugins()[triggers.Name]; ok {
		if p.registrar, ok = triggersPlugin.(triggers.IRegistrar); ok {
			if err := p.registrar.RegisterTrigger(NewTrigger(p.eventBus)); err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}

	// run server
	p.server = NewServer(p.adaptors,
		NewConfig(service.AppConfig()),
		p.scriptService,
		service.GateClient(),
		p.eventBus)

	p.server.Start()

	p.eventBus.Subscribe(TopicPluginAlexa, p.eventHandler)

	return nil
}

func (p *plugin) Unload() error {
	if !p.isStarted.Load() {
		return nil
	}
	p.isStarted.Store(false)

	p.eventBus.Unsubscribe(TopicPluginAlexa, p.eventHandler)

	p.server.Stop()
	p.server = nil

	if err := p.registrar.UnregisterTrigger(TriggerName); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (p *plugin) Name() string {
	return Name
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {

	return
}

func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {

	return
}

func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
}

func (p *plugin) Depends() []string {
	return []string{"triggers"}
}

func (p *plugin) Version() string {
	return "0.0.1"
}

func (p *plugin) Server() IServer {
	return p.server
}

func (p *plugin) eventHandler(_ string, event interface{}) {

	switch v := event.(type) {
	case EventAlexaAddSkill:
		p.server.AddSkill(v.Skill)
	case EventAlexaUpdateSkill:
		p.server.UpdateSkill(v.Skill)
	case EventAlexaDeleteSkill:
		p.server.DeleteSkill(v.Skill)
	}
}
