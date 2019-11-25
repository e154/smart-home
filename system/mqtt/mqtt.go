package mqtt

import (
	"context"
	"fmt"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/DrmagicE/gmqtt/plugin/management"
	"github.com/e154/smart-home/system/graceful_service"
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
}

func NewMqtt(cfg *MqttConfig,
	graceful *graceful_service.GracefulService,
	authenticator *Authenticator,
	scriptService *scripts.ScriptService ) (mqtt *Mqtt) {

	mqtt = &Mqtt{
		cfg:           cfg,
		authenticator: authenticator,
	}

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
		panic(err.Error())
	}

	m.server.AddTCPListenner(ln)

	m.server.AddPlugins(management.New(":8081", nil))

	m.hooks()

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

func (m *Mqtt) hooks() {

	//authentication
	m.server.RegisterOnConnect(func(cs gmqtt.ChainStore, client gmqtt.Client) (code uint8) {

		username := client.OptionsReader().Username()
		password := client.OptionsReader().Password()

		if err := m.authenticator.Authenticate(username, password); err != nil {
			return packets.CodeBadUsernameorPsw
		}

		return packets.CodeAccepted
	})

	//m.server.RegisterOnConnected(func(cs gmqtt.ChainStore, client gmqtt.Client) {
	//	log.Debug("connected...")
	//})
	//
	//m.server.RegisterOnSessionCreated(func(cs gmqtt.ChainStore, client gmqtt.Client) {
	//	log.Debug("session created...")
	//})
	//
	//m.server.RegisterOnSessionResumed(func(cs gmqtt.ChainStore, client gmqtt.Client) {
	//	log.Debug("session resumed...")
	//})

	m.server.RegisterOnSubscribe(func(cs gmqtt.ChainStore, client gmqtt.Client, topic packets.Topic) (qos uint8) {
		log.Debugf("subscribe: %v", topic.Name)
		if topic.Name == "test/nosubscribe" {
			return packets.SUBSCRIBE_FAILURE
		}
		return topic.Qos
	})
}
