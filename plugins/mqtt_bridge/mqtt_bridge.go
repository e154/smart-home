// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package mqtt_bridge

import (
	"context"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/supervisor"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/system/mqtt"
)

type MqttBridge struct {
	cfg          *Config
	client       MQTT.Client
	server       mqtt.MqttServ
	serverClient mqtt.MqttCli
	isStarted    *atomic.Bool
	actor        *Actor
}

func NewMqttBridge(cfg *Config, server mqtt.MqttServ, actor *Actor) (bridge *MqttBridge, err error) {
	bridge = &MqttBridge{
		cfg:          cfg,
		server:       server,
		isStarted:    atomic.NewBool(false),
		serverClient: server.NewClient(uuid.NewString()),
		actor:        actor,
	}
	return
}

func (m *MqttBridge) Start(ctx context.Context) (err error) {
	if m.isStarted.Load() {
		return
	}
	m.isStarted.Store(true)
	defer func() {
		if err != nil {
			m.isStarted.Store(false)
		}
	}()

	opts := MQTT.NewClientOptions().
		AddBroker(m.cfg.Broker).
		SetClientID(m.cfg.ClientID).
		SetKeepAlive(time.Duration(m.cfg.KeepAlive) * time.Second).
		SetPingTimeout(time.Duration(m.cfg.PingTimeout) * time.Second).
		SetConnectTimeout(time.Duration(m.cfg.ConnectTimeout) * time.Second).
		SetCleanSession(m.cfg.CleanSession).
		SetOnConnectHandler(m.onConnect).
		SetConnectionLostHandler(m.onConnectionLostHandler).
		SetUsername(m.cfg.Username).
		SetPassword(m.cfg.Password)
	m.client = MQTT.NewClient(opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
	}
	return
}

func (m *MqttBridge) Shutdown(ctx context.Context) (err error) {
	if !m.isStarted.Load() {
		return
	}
	m.isStarted.Store(false)
	if m.client != nil {
		m.client.Disconnect(250)
	}
	return
}

func (m *MqttBridge) onConnect(client MQTT.Client) {

	switch m.cfg.Direction {
	case DirectionBoth:
		m.directionBoth(client)
	case DirectionIn:
		m.directionIn(client)
	case DirectionOut:
		m.directionOut(client)
	default:
		log.Debugf("unknown direction '%s'", m.cfg.Direction)
	}

	m.actor.SetState(supervisor.EntityStateParams{
		NewState:    common.String(AttrConnected),
		StorageSave: true,
	})

	log.Debug("connected ...")
}

func (m *MqttBridge) onConnectionLostHandler(client MQTT.Client, e error) {

	m.actor.SetState(supervisor.EntityStateParams{
		NewState:    common.String(AttrOffline),
		StorageSave: true,
	})

	log.Debug("connection lost...")

	for _, topic := range m.cfg.Topics {
		if token := m.client.Unsubscribe(topic); token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
	m.serverClient.UnsubscribeAll()
}

func (m *MqttBridge) directionBoth(client MQTT.Client) {

	for _, topic := range m.cfg.Topics {
		log.Debugf("subscribe: %s", topic)
		if token := client.Subscribe(topic, m.cfg.Qos, m.clientMessageHandler); token.Wait() && token.Error() != nil {
			log.Error(token.Error().Error())
		}

		if err := m.serverClient.Subscribe(topic, m.serverMessageHandler); err != nil {
			log.Error(err.Error())
		}
	}
}

func (m *MqttBridge) directionIn(client MQTT.Client) {

	for _, topic := range m.cfg.Topics {
		log.Debugf("subscribe: %s", topic)
		if token := client.Subscribe(topic, m.cfg.Qos, m.clientMessageHandler); token.Wait() && token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
}

func (m *MqttBridge) directionOut(client MQTT.Client) {

	var err error
	for _, topic := range m.cfg.Topics {
		if err = m.serverClient.Subscribe(topic, m.serverMessageHandler); err != nil {
			log.Error(err.Error())
		}
	}
}

func (m *MqttBridge) serverMessageHandler(client mqtt.MqttCli, message mqtt.Message) {
	m.client.Publish(message.Topic, message.Qos, message.Retained, message.Payload)
}

func (m *MqttBridge) clientMessageHandler(_ MQTT.Client, message MQTT.Message) {
	if err := m.server.Publish(message.Topic(), message.Payload(), message.Qos(), message.Retained()); err != nil {
		log.Error(err.Error())
	}
}
