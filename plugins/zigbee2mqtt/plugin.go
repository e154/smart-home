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

package zigbee2mqtt

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/e154/smart-home/common"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.zigbee2mqtt")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	mqttServ   mqtt.MqttServ
	mqttClient mqtt.MqttCli
	mqttSubs   sync.Map
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin:   supervisor.NewPlugin(),
		mqttSubs: sync.Map{},
	}
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	p.mqttServ = service.MqttServ()

	p.mqttClient = p.mqttServ.NewClient("plugins.zigbee2mqtt")
	if err := p.Service.EventBus().Subscribe("system/entities/+", p.eventHandler); err != nil {
		log.Error(err.Error())
	}
	return nil
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	p.mqttServ.RemoveClient("plugins.zigbee2mqtt")
	_ = p.Service.EventBus().Unsubscribe("system/entities/+", p.eventHandler)
	return
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (supervisor.PluginActor, error) {

	actor, err := NewActor(entity, p.Service)
	if err != nil {
		return nil, err
	}
	var br *m.Zigbee2mqtt
	if br, err = p.Service.Adaptors().Zigbee2mqtt.GetById(context.Background(), actor.zigbee2mqttDevice.Zigbee2mqttId); err != nil {
		return nil, err
	}

	if _, ok := p.mqttSubs.Load(br.Id); !ok {
		_ = p.mqttClient.Subscribe(p.topic(br.BaseTopic), p.mqttOnPublish)
		p.mqttSubs.Store(br.Id, nil)
	}
	return actor, nil
}

// Name ...
func (p plugin) Name() string {
	return Name
}

func (p *plugin) topic(bridgeId string) string {
	return fmt.Sprintf("%s/#", bridgeId)
}

func (p *plugin) mqttOnPublish(client mqtt.MqttCli, msg mqtt.Message) {

	var topic = strings.Split(msg.Topic, "/")

	if len(topic) == 0 {
		return
	}

	value, ok := p.Actors.Load(common.EntityId(fmt.Sprintf("zigbee2mqtt.%s", topic[1])))
	if ok {
		actor := value.(*Actor)
		actor.mqttOnPublish(client, msg)
	}
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventCallEntityAction:
		value, ok := p.Actors.Load(v.EntityId)
		if !ok {
			return
		}
		actor := value.(*Actor)
		actor.addAction(v)
	}
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginInstallable
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
		ActorCustomActions: true,
		ActorCustomStates:  true,
		ActorCustomSetts:   true,
		Setts:              nil,
	}
}
