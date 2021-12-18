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

package node

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	eventBus      event_bus.EventBus
	mqttClient    mqtt.MqttCli
	stateMu       *sync.Mutex
	quit          chan struct{}
	lastPing      time.Time
	lastState     m.AttributeValue
}

// NewActor ...
func NewActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	eventBus event_bus.EventBus,
	mqttClient mqtt.MqttCli) (actor *Actor) {

	actor = &Actor{
		BaseActor:     entity_manager.NewBaseActor(entity, scriptService, adaptors),
		adaptors:      adaptors,
		scriptService: scriptService,
		eventBus:      eventBus,
		mqttClient:    mqttClient,
		stateMu:       &sync.Mutex{},
		lastPing:      time.Time{},
		lastState:     m.AttributeValue{},
	}

	actor.Manager = entityManager
	actor.Attrs = NewAttr()
	actor.States = NewStates()
	actor.Setts = entity.Settings

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	return actor
}

func (e *Actor) destroy() {

	e.mqttClient.Unsubscribe(e.mqttTopic("resp/+"))
	e.mqttClient.Unsubscribe(e.mqttTopic("ping"))
	e.eventBus.Unsubscribe(e.localTopic("req/+"), e.onMessage)

	e.quit <- struct{}{}
}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {

	e.quit = make(chan struct{})

	state := e.States["wait"]
	e.State = &state

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				e.updateStatus()

			case <-e.quit:
				//close(e.quit)
				return
			}
		}
	}()

	// local sub
	e.eventBus.Subscribe(e.localTopic("req/+"), e.onMessage)

	// mqtt sub
	e.mqttClient.Subscribe(e.mqttTopic("resp/+"), e.mqttOnMessage)
	e.mqttClient.Subscribe(e.mqttTopic("ping"), e.ping)

	return e
}

// event from plugin.node/nodeName/req
func (e *Actor) onMessage(_ string, msg MessageRequest) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Error(err.Error())
	}
	e.mqttClient.Publish(e.mqttTopic(fmt.Sprintf("req/%s", msg.EntityId)), b)
}

// event from home/node/nodeName/#
func (e *Actor) mqttOnMessage(_ mqtt.MqttCli, msg mqtt.Message) {
	// resend msg to plugin.node/nodeName/resp/entityId
	resp := MessageResponse{}
	if err := json.Unmarshal(msg.Payload, &resp); err != nil {
		log.Warn(err.Error())
		return
	}
	items := strings.Split(msg.Topic, "/")
	entityId := items[len(items)-1]
	e.eventBus.Publish(e.localTopic(fmt.Sprintf("resp/%s", entityId)), resp)
}

// event from home/node/nodeName/ping
func (e *Actor) ping(_ mqtt.MqttCli, msg mqtt.Message) {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()

	e.lastPing = time.Now()

	json.Unmarshal(msg.Payload, &e.lastState)
}

func (e *Actor) updateStatus() {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()

	var state = "wait"
	if time.Now().Sub(e.lastPing).Seconds() < 2 {
		state = "connected"
	}

	oldState := e.GetEventState(e)
	e.Now(oldState)

	e.AttrMu.Lock()
	changed, _ := e.Attrs.Deserialize(e.lastState)
	e.AttrMu.Unlock()

	if !changed && e.State.Name == state {
		return
	}

	if state, ok := e.States[state]; ok {
		e.State = &state
	}

	e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventStateChanged{
		Type:     e.Id.Type(),
		EntityId: e.Id,
		OldState: oldState,
		NewState: e.GetEventState(e),
	})
}

func (e *Actor) localTopic(r string) string {
	return fmt.Sprintf("%s/%s/%s", TopicPluginNode, e.Name, r)
}

func (e *Actor) mqttTopic(r string) string {
	return fmt.Sprintf("home/node/%s/%s", e.Name, r)
}
