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

package node

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/common"
	"strings"
	"sync"
	"time"

	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	mqttClient mqtt.MqttCli
	stateMu    *sync.Mutex
	quit       chan struct{}
	lastPing   time.Time
	lastState  m.AttributeValue
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service,
	mqttClient mqtt.MqttCli) (actor *Actor) {

	actor = &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		mqttClient: mqttClient,
		stateMu:    &sync.Mutex{},
		lastPing:   time.Time{},
		lastState:  m.AttributeValue{},
	}

	if actor.Attrs == nil || len(actor.Attrs) == 0 {
		actor.Attrs = NewAttr()
	}

	if actor.States == nil || len(actor.States) == 0 {
		actor.States = NewStates()
	}

	if actor.Setts == nil || len(actor.Setts) == 0 {
		actor.Setts = NewSettings()
	}

	return actor
}

func (e *Actor) Destroy() {

	e.mqttClient.Unsubscribe(e.mqttTopic("resp/+"))
	e.mqttClient.Unsubscribe(e.mqttTopic("ping"))
	_ = e.Service.EventBus().Unsubscribe(e.localTopic("req/+"), e.onMessage)

	close(e.quit)
}

// Spawn ...
func (e *Actor) Spawn() {

	e.SetActorState(common.String("wait"))

	go func() {
		e.quit = make(chan struct{})

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				e.updateStatus()

			case <-e.quit:
				return
			}
		}
	}()

	// local sub
	_ = e.Service.EventBus().Subscribe(e.localTopic("req/+"), e.onMessage)

	// mqtt sub
	_ = e.mqttClient.Subscribe(e.mqttTopic("resp/+"), e.mqttOnMessage)
	_ = e.mqttClient.Subscribe(e.mqttTopic("ping"), e.ping)

}

// event from plugin.node/nodeName/req
func (e *Actor) onMessage(_ string, msg MessageRequest) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Error(err.Error())
	}
	_ = e.mqttClient.Publish(e.mqttTopic(fmt.Sprintf("req/%s", msg.EntityId)), b)
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
	e.Service.EventBus().Publish(e.localTopic(fmt.Sprintf("resp/%s", entityId)), resp)
}

// event from home/node/nodeName/ping
func (e *Actor) ping(_ mqtt.MqttCli, msg mqtt.Message) {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()

	e.lastPing = time.Now()

	_ = json.Unmarshal(msg.Payload, &e.lastState)
}

func (e *Actor) updateStatus() {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()

	var state = "wait"
	if time.Since(e.lastPing).Seconds() < 2 {
		state = "connected"
	}

	e.DeserializeAttr(e.lastState)
	e.SetActorState(&state)
	e.SaveState(false, true)
}

func (e *Actor) localTopic(r string) string {
	return fmt.Sprintf("%s/%s/%s", TopicPluginNode, e.Name, r)
}

func (e *Actor) mqttTopic(r string) string {
	return fmt.Sprintf("system/plugins/node/%s/%s", e.Name, r)
}
