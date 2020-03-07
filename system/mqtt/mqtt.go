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
	"fmt"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/e154/smart-home/system/mqtt/prometheus"
	"net/http"

	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/mqtt/management"
	"github.com/e154/smart-home/system/mqtt_authenticator"
	"github.com/e154/smart-home/system/mqtt_client"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/uuid"
	"github.com/op/go-logging"
	"net"
)

var (
	log = logging.MustGetLogger("mqtt")
)

type Mqtt struct {
	cfg           *MqttConfig
	server        IMQTT
	clients       []*mqtt_client.Client
	authenticator *mqtt_authenticator.Authenticator
	management    *management.Management
}

func NewMqtt(cfg *MqttConfig,
	graceful *graceful_service.GracefulService,
	authenticator *mqtt_authenticator.Authenticator,
	scriptService *scripts.ScriptService) (mqtt *Mqtt) {

	mqtt = &Mqtt{
		cfg:           cfg,
		authenticator: authenticator,
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
		gmqtt.WithPlugin(m.management),
		gmqtt.WithHook(gmqtt.Hooks{
			OnConnect:        m.OnConnect,
			OnConnected:      m.OnConnected,
			OnClose:          m.OnClose,
			OnSessionCreated: m.OnSessionCreated,
			OnSessionResumed: m.OnSessionResumed,
		}),
		gmqtt.WithPlugin(prometheus.New(&http.Server{
			Addr: ":8082",
		}, "/metrics")),
	)

	log.Infof("Serving server at tcp://[::]:%d", m.cfg.Port)

	m.server.Run()
}

func (m *Mqtt) NewClient(cfg *mqtt_client.Config) (c *mqtt_client.Client, err error) {

	if cfg == nil {
		cfg = &mqtt_client.Config{
			KeepAlive:      5,
			PingTimeout:    5,
			ConnectTimeout: 5,
			Qos:            0,
			CleanSession:   true,
		}
	}

	cfg.Username = m.authenticator.Login()
	cfg.Password = m.authenticator.Password()
	if cfg.ClientID == "" {
		cfg.ClientID = uuid.NewV4().String()
	}
	cfg.Broker = fmt.Sprintf("tcp://127.0.0.1:%d", m.cfg.Port)

	if c, err = mqtt_client.NewClient(cfg); err != nil {
		return
	}

	m.clients = append(m.clients, c)

	return
}

func (m *Mqtt) OnConnected(ctx context.Context, client gmqtt.Client) {
	log.Debugf("%v connected...", client.OptionsReader().ClientID())
}

func (m *Mqtt) OnClose(ctx context.Context, client gmqtt.Client, err error) {
	log.Debugf("%v disconnected...", client.OptionsReader().ClientID())
}

func (m *Mqtt) OnSessionCreated(ctx context.Context, client gmqtt.Client) {
	log.Debug("session created...")
}

func (m *Mqtt) OnSessionResumed(ctx context.Context, client gmqtt.Client) {
	log.Debug("session resumed...")
}

func (m *Mqtt) OnConnect(ctx context.Context, client gmqtt.Client) (code uint8) {

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

func (m *Mqtt) Publish(topic string, payload []byte, qos uint8, retain bool) {
	m.server.PublishService().Publish(gmqtt.NewMessage(topic, payload, qos, gmqtt.Retained(retain)))
}
