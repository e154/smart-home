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

package zigbee2mqtt

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/plugin_manager"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
	"strings"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.zigbee2mqtt")
)

type pluginZigbee2mqtt struct {
	entityManager *entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	scriptService *scripts.ScriptService
	isStarted     *atomic.Bool
	eventBus      *event_bus.EventBus
	actorsLock    *sync.Mutex
	actors        map[string]*EntityActor
	mqtt          *mqtt.Mqtt
	mqttClient    *mqtt.Client
	mqttSubs      sync.Map
}

func Register(manager *plugin_manager.PluginManager,
	entityManager *entity_manager.EntityManager,
	eventBus *event_bus.EventBus,
	adaptors *adaptors.Adaptors,
	mqtt *mqtt.Mqtt,
	scriptService *scripts.ScriptService) {
	manager.Register(&pluginZigbee2mqtt{
		entityManager: entityManager,
		adaptors:      adaptors,
		scriptService: scriptService,
		isStarted:     atomic.NewBool(false),
		eventBus:      eventBus,
		actorsLock:    &sync.Mutex{},
		actors:        make(map[string]*EntityActor),
		mqtt:          mqtt,
		mqttClient:    mqtt.NewClient("plugins.zigbee2mqtt"),
		mqttSubs:      sync.Map{},
	})
}

func (p pluginZigbee2mqtt) Load(service plugin_manager.IPluginManager, plugins map[string]interface{}) (err error) {
	if err := p.eventBus.Subscribe(event_bus.TopicEntities, p.eventHandler); err != nil {
		log.Error(err.Error())
	}
	return
}

func (p pluginZigbee2mqtt) Unload() (err error) {
	return
}

func (p pluginZigbee2mqtt) Name() string {
	return Name
}

func (p *pluginZigbee2mqtt) AddOrUpdateEntity(entity *m.Entity) error {
	return p.addOrUpdateEntity(entity, entity.Attributes.Serialize())
}

func (p *pluginZigbee2mqtt) RemoveEntity(entity *m.Entity) (err error) {
	return p.removeEntity(entity.Id.Name())
}

func (p *pluginZigbee2mqtt) addOrUpdateEntity(entity *m.Entity, attributes m.EntityAttributeValue) (err error) {
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

	log.Infof("Add zigbee2mqtt '%s'", name)

	var actor *EntityActor
	if actor, err = NewEntityActor(entity, attributes, p.adaptors, p.scriptService); err != nil {
		return
	}
	p.actors[name] = actor
	p.entityManager.Spawn(p.actors[name].Spawn)

	var br *m.Zigbee2mqtt
	if br, err = p.adaptors.Zigbee2mqtt.GetById(actor.zigbee2mqttDevice.Zigbee2mqttId); err != nil {
		return
	}

	if _, ok := p.mqttSubs.Load(br.Id); !ok {
		p.mqttClient.Subscribe(p.topic(br.BaseTopic), p.mqttOnPublish)
		p.mqttSubs.Store(br.Id, nil)
	}

	return
}

func (p *pluginZigbee2mqtt) removeEntity(name string) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[name]; !ok {
		return
	}

	delete(p.actors, name)

	return
}

func (p *pluginZigbee2mqtt) topic(bridgeId string) string {
	return fmt.Sprintf("%s/#", bridgeId)
}

func (p *pluginZigbee2mqtt) mqttOnPublish(client *mqtt.Client, msg mqtt.Message) {

	var topic = strings.Split(msg.Topic, "/")

	if len(topic) == 0 {
		return
	}

	var actor *EntityActor
	var err error
	if actor, err = p.getActorByZigbeeDeviceId(topic[1]); err != nil {
		log.Warn(err.Error())
		return
	}

	actor.mqttOnPublish(client, msg)
}

func (p *pluginZigbee2mqtt) getActorByZigbeeDeviceId(deviceId string) (actor *EntityActor, err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	for _, actor = range p.actors {
		if actor.zigbee2mqttDevice.Id == deviceId {
			return
		}
	}

	err = fmt.Errorf("device \"%s\" not found", deviceId)

	return
}

func (p *pluginZigbee2mqtt) eventHandler(msg interface{}) {

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

func (p *pluginZigbee2mqtt) Type() plugin_manager.PlugableType {
	return plugin_manager.PlugableBuiltIn
}

func (p *pluginZigbee2mqtt) Depends() []string {
	return nil
}
