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

package mqtt

import (
	"context"
	"errors"
	"fmt"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/mqtt/management"
	"github.com/e154/smart-home/system/mqtt/metric"
	"github.com/e154/smart-home/system/mqtt_authenticator"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/zap"
	"net"
	"sync"
)

var (
	log = common.MustGetLogger("mqtt")
)

type Mqtt struct {
	cfg            *MqttConfig
	server         IMQTT
	publishService gmqtt.PublishService
	authenticator  *mqtt_authenticator.Authenticator
	management     *management.Management
	metric         *metrics.MetricManager
	clientsLock    *sync.Mutex
	clients        map[string]*Client
}

func NewMqtt(cfg *MqttConfig,
	graceful *graceful_service.GracefulService,
	authenticator *mqtt_authenticator.Authenticator,
	scriptService *scripts.ScriptService,
	metric *metrics.MetricManager) (mqtt *Mqtt) {

	mqtt = &Mqtt{
		cfg:           cfg,
		authenticator: authenticator,
		metric:        metric,
		clientsLock:   &sync.Mutex{},
		clients:       make(map[string]*Client),
	}

	// javascript binding
	scriptService.PushStruct("Mqtt", NewMqttBind(mqtt))

	mqtt.runServer()

	graceful.Subscribe(mqtt)

	return
}

func (m *Mqtt) Shutdown() {
	log.Info("Server exiting")
	if m.server != nil {
		_ = m.server.Stop(context.Background())
	}
}

func (m *Mqtt) runServer() {

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", m.cfg.Port))
	if err != nil {
		log.Error(err.Error())
	}

	m.management = management.New()

	// Create a new server
	m.server = gmqtt.NewServer(
		gmqtt.WithTCPListener(ln),
		gmqtt.WithHook(gmqtt.Hooks{
			OnConnect:        m.OnConnect,
			OnConnected:      m.OnConnected,
			OnClose:          m.OnClose,
			OnSessionCreated: m.OnSessionCreated,
			OnSessionResumed: m.OnSessionResumed,
			OnSubscribe:      m.OnSubscribe,
			OnUnsubscribed:   m.OnUnsubscribed,
			OnMsgArrived:     m.OnMsgArrived,
		}),
		gmqtt.WithPlugin(m.management),
		//gmqtt.WithPlugin(prometheus.New(&http.Server{Addr: ":8082",}, "/metrics")),
		gmqtt.WithPlugin(metric.New(m.metric, 5)),
		gmqtt.WithLogger(zap.L().Named("gmqtt")),
	)
	m.publishService = m.server.PublishService()

	log.Infof("Serving server at tcp://[::]:%d", m.cfg.Port)

	m.server.Run()
}

func (m *Mqtt) OnConnected(ctx context.Context, client gmqtt.Client) {
	log.Debugf("connected... %v", client.OptionsReader().ClientID())
}

func (m *Mqtt) OnClose(ctx context.Context, client gmqtt.Client, err error) {
	log.Debugf("disconnected... %v", client.OptionsReader().ClientID())
}

func (m *Mqtt) OnSessionCreated(ctx context.Context, client gmqtt.Client) {
	log.Debugf("session created... %v", client.OptionsReader().ClientID())
}

func (m *Mqtt) OnSessionResumed(ctx context.Context, client gmqtt.Client) {
	log.Debugf("session resumed... %v", client.OptionsReader().ClientID())
}

func (m *Mqtt) OnSubscribe(ctx context.Context, client gmqtt.Client, topic packets.Topic) (qos uint8) {
	log.Debugf("subscribe %v to topic %v", client.OptionsReader().ClientID(), topic.Name)
	return topic.Qos
}

func (m *Mqtt) OnUnsubscribed(ctx context.Context, client gmqtt.Client, topicName string) {
	log.Debugf("unsubscribe %v from topic %v", client.OptionsReader().ClientID(), topicName)
}

func (m *Mqtt) OnMsgArrived(ctx context.Context, client gmqtt.Client, msg packets.Message) (valid bool) {
	m.clientsLock.Lock()
	defer m.clientsLock.Unlock()

	for _, cli := range m.clients {
		cli.OnMsgArrived(ctx, client, msg)
	}

	return true
}

func (m *Mqtt) OnConnect(ctx context.Context, client gmqtt.Client) (code uint8) {
	log.Debugf("connect... %v", client.OptionsReader().ClientID())

	username := client.OptionsReader().Username()
	password := client.OptionsReader().Password()

	//authentication
	if err := m.authenticator.Authenticate(username, password); err != nil {
		return packets.CodeBadUsernameorPsw
	}

	return packets.CodeAccepted
}

func (m *Mqtt) Management() IManagement {
	return m.management
}

func (m *Mqtt) Publish(topic string, payload []byte, qos uint8, retain bool) (err error) {
	if qos < 0 || qos > 2 {
		err = errors.New("invalid Qos")
		return
	}
	if !packets.ValidTopicFilter([]byte(topic)) {
		err = errors.New("invalid topic filter")
		return
	}
	if !packets.ValidUTF8([]byte(payload)) {
		err = errors.New("invalid utf-8 string")
		return
	}
	m.publishService.Publish(gmqtt.NewMessage(topic, payload, qos, gmqtt.Retained(retain)))
	return
}

func (m *Mqtt) NewClient(name string) (client *Client) {
	m.clientsLock.Lock()
	defer m.clientsLock.Unlock()

	var ok bool
	if client, ok = m.clients[name]; ok {
		return
	}
	client = NewClient(m, name)
	m.clients[name] = client
	return
}
