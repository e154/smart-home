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

package mqtt

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/e154/smart-home/internal/system/logging"
	"github.com/e154/smart-home/internal/system/mqtt/admin"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/logger"
	mqttType "github.com/e154/smart-home/pkg/mqtt"
	"github.com/e154/smart-home/pkg/scripts"

	"github.com/DrmagicE/gmqtt"
	_ "github.com/DrmagicE/gmqtt/persistence"
	"github.com/DrmagicE/gmqtt/pkg/codes"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/DrmagicE/gmqtt/server"
	_ "github.com/DrmagicE/gmqtt/topicalias/fifo"
	"github.com/grandcat/zeroconf"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/e154/bus"
)

var (
	log = logger.MustGetLogger("mqtt")
)

var _ mqttType.MqttServ = (*Mqtt)(nil)
var _ MqttServAdmin = (*Mqtt)(nil)

// Mqtt ...
type Mqtt struct {
	cfg           *Config
	server        mqttType.GMqttServer
	authenticator mqttType.MqttAuthenticator
	isStarted     bool
	clientsLock   *sync.Mutex
	clients       map[string]mqttType.MqttCli
	admin         *admin.Admin
	scriptService scripts.ScriptService
	eventBus      bus.Bus
	zeroconf      *zeroconf.Server
}

// NewMqtt ...
func NewMqtt(lc fx.Lifecycle,
	cfg *Config,
	authenticator mqttType.MqttAuthenticator,
	scriptService scripts.ScriptService,
	eventBus bus.Bus) (mqtt mqttType.MqttServ) {

	mqtt = &Mqtt{
		cfg:           cfg,
		authenticator: authenticator,
		clientsLock:   &sync.Mutex{},
		clients:       make(map[string]mqttType.MqttCli),
		admin:         admin.New(),
		scriptService: scriptService,
		eventBus:      eventBus,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			mqtt.Start()
			return nil
		},
		OnStop: func(ctx context.Context) (err error) {
			return mqtt.Shutdown()
		},
	})

	return
}

// Shutdown ...
func (m *Mqtt) Shutdown() (err error) {
	if !m.isStarted {
		return
	}

	log.Info("Server exiting")

	m.scriptService.PopStruct("Mqtt")

	m.clientsLock.Lock()
	for name, cli := range m.clients {
		cli.UnsubscribeAll()
		delete(m.clients, name)
	}
	m.clientsLock.Unlock()

	if m.server != nil {
		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(100*time.Millisecond))
		err = m.server.Stop(ctx)
	}

	if m.zeroconf != nil {
		m.zeroconf.Shutdown()
	}

	m.eventBus.Publish("system/services/mqtt", events.EventServiceStopped{Service: "Mqtt"})
	return
}

// Start ...
func (m *Mqtt) Start() {

	if m.isStarted {
		return
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", m.cfg.Port))
	if err != nil {
		log.Error(err.Error())
	}

	defer func() {
		if err == nil {
			m.isStarted = true
		}
	}()

	options := []server.Options{
		server.WithTCPListener(ln),
		server.WithPlugin(m.admin),
		server.WithHook(server.Hooks{
			OnBasicAuth:  m.onBasicAuth,
			OnMsgArrived: m.onMsgArrived,
			OnConnected: func(ctx context.Context, client server.Client) {
				m.eventBus.Publish("system/services/mqtt", events.EventMqttNewClient{
					ClientId: client.ClientOptions().ClientID,
				})
			},
		}),
	}

	if m.cfg.Logging {
		options = append(options, server.WithLogger(m.logging()))
	}

	// Create a new server
	m.server = server.New(options...)

	log.Infof("Serving MQTT server at tcp://[::]:%d", m.cfg.Port)

	m.scriptService.PushStruct("Mqtt", NewMqttBind(m))

	go func() {
		if err = m.server.Run(); err != nil {
			log.Error(err.Error())
		}
	}()

	if m.zeroconf, err = zeroconf.Register("smart-home", "_mqtt._tcp", "local.", m.cfg.Port, nil, nil); err != nil {
		log.Error(err.Error())
	}

	m.eventBus.Publish("system/services/mqtt", events.EventServiceStarted{Service: "Mqtt"})
}

// OnMsgArrived ...
func (m *Mqtt) onMsgArrived(ctx context.Context, client server.Client, msg *server.MsgArrivedRequest) (err error) {
	m.clientsLock.Lock()
	defer m.clientsLock.Unlock()

	for _, cli := range m.clients {
		cli.OnMsgArrived(ctx, client, msg)
	}

	return
}

// OnConnect ...
func (m *Mqtt) onBasicAuth(ctx context.Context, client server.Client, req *server.ConnectRequest) (err error) {
	log.Debugf("connect client version %v ...", client.Version())

	username := string(req.Connect.Username)
	password := string(req.Connect.Password)

	//authentication
	if err = m.authenticator.Authenticate(username, password); err == nil {
		return
	}

	// check the client version, return a compatible reason code.
	switch client.Version() {
	case packets.Version5:
		return codes.NewError(codes.BadUserNameOrPassword)
	case packets.Version311:
		return codes.NewError(codes.V3BadUsernameorPassword)
	}
	// return nil if pass authentication.
	return nil
}

// Admin ...
func (m *Mqtt) Admin() Admin {
	return m.admin
}

// Publish ...
func (m *Mqtt) Publish(topic string, payload []byte, qos uint8, retain bool) (err error) {
	if qos < 0 || qos > 2 {
		err = mqttType.ErrInvalidQos
		return
	}
	if !packets.ValidTopicFilter(true, []byte(topic)) {
		err = mqttType.ErrInvalidTopicFilter
		return
	}
	if !packets.ValidUTF8(payload) {
		err = mqttType.ErrInvalidUtf8String
		return
	}

	m.server.Publisher().Publish(&gmqtt.Message{
		QoS:      qos,
		Retained: retain,
		Topic:    topic,
		Payload:  payload,
	})

	// send to local subscribers
	_ = m.onMsgArrived(context.TODO(), nil, &server.MsgArrivedRequest{
		Message: &gmqtt.Message{
			QoS:      qos,
			Retained: retain,
			Topic:    topic,
			Payload:  payload,
		},
	})
	return
}

// NewClient ...
func (m *Mqtt) NewClient(name string) (client mqttType.MqttCli) {
	m.clientsLock.Lock()
	defer m.clientsLock.Unlock()

	var ok bool
	if client, ok = m.clients[name]; ok {
		return
	}
	client = NewClient(m, name)
	m.clients[name] = client
	log.Infof("new mqtt client '%s'", name)
	return
}

// RemoveClient ...
func (m *Mqtt) RemoveClient(name string) {
	m.clientsLock.Lock()
	defer m.clientsLock.Unlock()

	var ok bool
	if _, ok = m.clients[name]; !ok {
		return
	}
	delete(m.clients, name)
}

func (m *Mqtt) logging() *zap.Logger {

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	lowLevel := zapcore.ErrorLevel
	if m.cfg.DebugMode == common.ReleaseMode {
		lowLevel = zapcore.DebugLevel
	}
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < lowLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	var encConfig zapcore.EncoderConfig
	if m.cfg.DebugMode == common.ReleaseMode {
		encConfig = zap.NewProductionEncoderConfig()
	} else {
		encConfig = zap.NewDevelopmentEncoderConfig()
	}

	encConfig.EncodeTime = nil
	encConfig.EncodeName = logging.CustomNameEncoder
	encConfig.EncodeCaller = logging.CustomCallerEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encConfig)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Named("mqtt")
}

// Authenticator ...
func (m *Mqtt) Authenticator() mqttType.MqttAuthenticator {
	return m.authenticator
}
