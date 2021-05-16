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

package mqtt

import (
	"context"
	"errors"
	"fmt"
	"github.com/DrmagicE/gmqtt"
	cfg "github.com/DrmagicE/gmqtt/config"
	_ "github.com/DrmagicE/gmqtt/persistence"
	"github.com/DrmagicE/gmqtt/pkg/codes"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/DrmagicE/gmqtt/server"
	_ "github.com/DrmagicE/gmqtt/topicalias/fifo"
	"github.com/DrmagicE/gmqtt/plugin/admin"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/mqtt_authenticator"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"os"
	"sync"
)

var (
	log = common.MustGetLogger("mqtt")
)

// Mqtt ...
type Mqtt struct {
	cfg           *Config
	server        GMqttServer
	authenticator mqtt_authenticator.MqttAuthenticator
	metric        *metrics.MetricManager
	clientsLock   *sync.Mutex
	clients       map[string]MqttCli
	//admin         *admin.Admin
}

// NewMqtt ...
func NewMqtt(lc fx.Lifecycle,
	cfg *Config,
	authenticator mqtt_authenticator.MqttAuthenticator,
	scriptService scripts.ScriptService,
	metric *metrics.MetricManager,
) (mqtt MqttServ) {

	mqtt = &Mqtt{
		cfg:           cfg,
		authenticator: authenticator,
		//metric:        metric,
		clientsLock: &sync.Mutex{},
		clients:     make(map[string]MqttCli),
	}

	// javascript binding
	scriptService.PushStruct("Mqtt", NewMqttBind(mqtt))

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
	log.Info("Server exiting")

	m.clientsLock.Lock()
	for name, cli := range m.clients {
		cli.UnsubscribeAll()
		delete(m.clients, name)
	}
	m.clientsLock.Unlock()

	if m.server != nil {
		err = m.server.Stop(context.Background())
	}
	return
}

func (m *Mqtt) Start() {

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", m.cfg.Port))
	if err != nil {
		log.Error(err.Error())
	}

	//m.admin = &admin.Admin{}

	ad, _ := admin.New(cfg.DefaultConfig())

	options := []server.Options{
		server.WithTCPListener(ln),
		server.WithPlugin(ad),
		server.WithHook(server.Hooks{
			OnBasicAuth:  m.onBasicAuth,
			OnMsgArrived: m.onMsgArrived,
		}),
	}

	if m.cfg.Logging {
		options = append(options, server.WithLogger(m.logging()))
	}

	// Create a new server
	m.server = server.New(options...)

	log.Infof("Serving server at tcp://[::]:%d", m.cfg.Port)

	go func() {
		if err = m.server.Run(); err != nil {
			log.Error(err.Error())
		}
	}()
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
	log.Debugf("connect... %v", client.ClientOptions().ClientID)

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

// Management ...
func (m *Mqtt) Admin() Admin {
	//return m.admin
	return nil
}

// Publish ...
func (m *Mqtt) Publish(topic string, payload []byte, qos uint8, retain bool) (err error) {
	if qos < 0 || qos > 2 {
		err = errors.New("invalid Qos")
		return
	}
	if !packets.ValidTopicFilter(true, []byte(topic)) {
		err = errors.New("invalid topic filter")
		return
	}
	if !packets.ValidUTF8(payload) {
		err = errors.New("invalid utf-8 string")
		return
	}

	m.server.Publisher().Publish(&gmqtt.Message{
		QoS:      qos,
		Retained: retain,
		Topic:    topic,
		Payload:  payload,
	})

	// send to local subscribers
	m.onMsgArrived(nil, nil, &server.MsgArrivedRequest{
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
func (m *Mqtt) NewClient(name string) (client MqttCli) {
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
	return
}

func (m *Mqtt) logging() *zap.Logger {

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	lowLevel := zapcore.ErrorLevel
	if m.cfg.DebugMode == config.ReleaseMode {
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
	if m.cfg.DebugMode == config.ReleaseMode {
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
