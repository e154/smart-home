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

package zigbee2mqtt

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/plugins"
)

var (
	log = common.MustGetLogger("plugins.zigbee2mqtt")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	actorsLock *sync.Mutex
	actors     map[string]*Actor
	mqttServ   mqtt.MqttServ
	mqttClient mqtt.MqttCli
	mqttSubs   sync.Map
}

// New ...
func New() plugins.Plugable {
	return &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]*Actor),
		mqttSubs:   sync.Map{},
	}
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.mqttServ = service.MqttServ()

	p.mqttClient = p.mqttServ.NewClient("plugins.zigbee2mqtt")
	if err := p.EventBus.Subscribe(event_bus.TopicEntities, p.eventHandler); err != nil {
		log.Error(err.Error())
	}
	return nil
}

// Unload ...
func (p plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.mqttServ.RemoveClient("plugins.zigbee2mqtt")
	p.EventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)
	return
}

// Name ...
func (p plugin) Name() string {
	return Name
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) error {
	return p.addOrUpdateEntity(entity, entity.Attributes.Serialize())
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {
	return p.removeEntity(entityId.Name())
}

func (p *plugin) addOrUpdateEntity(entity *m.Entity, attributes m.AttributeValue) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	name := entity.Id.Name()
	if _, ok := p.actors[name]; ok {
		return
	}

	if actor, ok := p.actors[name]; ok {
		// update
		actor.SetState(entity_manager.EntityStateParams{
			AttributeValues: attributes,
		})
		return
	}

	var actor *Actor
	if actor, err = NewActor(entity, attributes,
		p.Adaptors, p.ScriptService, p.EntityManager, p.EventBus); err != nil {
		return
	}
	p.actors[name] = actor
	p.EntityManager.Spawn(p.actors[name].Spawn)

	var br *m.Zigbee2mqtt
	if br, err = p.Adaptors.Zigbee2mqtt.GetById(actor.zigbee2mqttDevice.Zigbee2mqttId); err != nil {
		return
	}

	if _, ok := p.mqttSubs.Load(br.Id); !ok {
		p.mqttClient.Subscribe(p.topic(br.BaseTopic), p.mqttOnPublish)
		p.mqttSubs.Store(br.Id, nil)
	}

	return
}

func (p *plugin) removeEntity(name string) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[name]; !ok {
		err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("failed remove '%s", name))
		return
	}

	delete(p.actors, name)

	return
}

func (p *plugin) topic(bridgeId string) string {
	return fmt.Sprintf("%s/#", bridgeId)
}

func (p *plugin) mqttOnPublish(client mqtt.MqttCli, msg mqtt.Message) {

	var topic = strings.Split(msg.Topic, "/")

	if len(topic) == 0 {
		return
	}

	var actor *Actor
	var err error
	if actor, err = p.getActorByZigbeeDeviceId(topic[1]); err != nil {
		//log.Warn(err.Error())
		return
	}

	actor.mqttOnPublish(client, msg)
}

func (p *plugin) getActorByZigbeeDeviceId(deviceId string) (actor *Actor, err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	for _, actor = range p.actors {
		if actor.zigbee2mqttDevice.Id == deviceId {
			return
		}
	}

	err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("device \"%s\" not found", deviceId))

	return
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case event_bus.EventCallAction:
		actor, ok := p.actors[v.EntityId.Name()]
		if !ok {
			return
		}
		actor.addAction(v)

	default:
		//fmt.Printf("new event: %v\n", reflect.TypeOf(v).String())
	}
}

// Type ...
func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginInstallable
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Version ...
func (p *plugin) Version() string {
	return "0.0.1"
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Triggers:           false,
		Actors:             true,
		ActorCustomAttrs:   true,
		//ActorAttrs:         NewAttr(),
		ActorCustomActions: true,
		//ActorActions:       entity_manager.ToEntityActionShort(NewActions()),
		ActorCustomStates:  true,
		//ActorStates:        entity_manager.ToEntityStateShort(NewStates()),
		ActorCustomSetts:   true,
		//ActorSetts:         NewSettings(),
		Setts:              nil,
	}
}
