package mqtt

import (
	"context"
	"fmt"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/mqtt/management"
	"github.com/e154/smart-home/system/scripts"
	"github.com/op/go-logging"
	"net"
	"time"
)

var (
	log = logging.MustGetLogger("mqtt")
)

type Mqtt struct {
	cfg           *MqttConfig
	server        *gmqtt.Server
	clients       []*Client
	authenticator *Authenticator
	management    *management.Management
}

func NewMqtt(cfg *MqttConfig,
	graceful *graceful_service.GracefulService,
	authenticator *Authenticator,
	scriptService *scripts.ScriptService) (mqtt *Mqtt) {

	mqtt = &Mqtt{
		cfg:           cfg,
		authenticator: authenticator,
		management:    management.NewManagement(),
	}

	// javascript binding
	scriptService.PushStruct("Mqtt", NewMqttBind(mqtt))

	go mqtt.runServer()

	graceful.Subscribe(mqtt)

	return
}

func (m *Mqtt) Shutdown() {
	log.Info("Server exiting")
	m.server.Stop(context.Background())
}

func (m *Mqtt) runServer() {

	config := gmqtt.Config{
		RetryInterval:              m.cfg.RetryInterval * time.Second,
		RetryCheckInterval:         m.cfg.RetryCheckInterval * time.Second,
		SessionExpiryInterval:      m.cfg.SessionExpiryInterval,
		SessionExpireCheckInterval: m.cfg.SessionExpireCheckInterval,
		QueueQos0Messages:          m.cfg.QueueQos0Messages,
		MaxInflight:                m.cfg.MaxInflight,
		MaxAwaitRel:                m.cfg.MaxAwaitRel,
		MaxMsgQueue:                m.cfg.MaxMsgQueue,
		DeliverMode:                m.cfg.DeliverMode,
	}

	// Create a new server
	m.server = gmqtt.NewServer(config)

	log.Infof("Serving server at tcp://[::]:%d", m.cfg.Port)

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", m.cfg.Port))
	if err != nil {
		log.Error(err.Error())
	}

	m.server.AddTCPListenner(ln)

	// management
	m.server.AddPlugins(m.management)

	m.server.RegisterOnConnect(m.OnConnect)
	m.server.RegisterOnConnected(m.OnConnected)
	m.server.RegisterOnSessionCreated(m.OnSessionCreated)
	m.server.RegisterOnSessionResumed(m.OnSessionResumed)

	m.server.Run()
}

func (m *Mqtt) NewClient(topic string) (c *Client, err error) {

	uri := fmt.Sprintf("tcp://127.0.0.1:%d", m.cfg.Port)
	log.Infof("new queue client %s topic(%s)", uri, topic)

	if c, err = NewClient(uri, topic, m.authenticator); err != nil {
		return
	}

	m.clients = append(m.clients, c)

	return
}

func (m *Mqtt) OnConnected(cs gmqtt.ChainStore, client gmqtt.Client) {
	log.Debug("connected...")
}

func (m *Mqtt) OnSessionCreated(cs gmqtt.ChainStore, client gmqtt.Client) {
	log.Debug("session created...")
}

func (m *Mqtt) OnSessionResumed(cs gmqtt.ChainStore, client gmqtt.Client) {
	log.Debug("session resumed...")
}

func (m *Mqtt) OnConnect(cs gmqtt.ChainStore, client gmqtt.Client) (code uint8) {

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
