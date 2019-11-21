package mqtt

import (
	"context"
	"fmt"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/e154/smart-home/system/graceful_service"
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
	authenticator *Authenticator) (mqtt *Mqtt) {

	mqtt = &Mqtt{
		cfg:           cfg,
		authenticator: authenticator,
	}

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
		RetryInterval:              20 * time.Second,
		RetryCheckInterval:         20 * time.Second,
		SessionExpiryInterval:      0,
		SessionExpireCheckInterval: 0,
		QueueQos0Messages:          true,
		MaxInflight:                32,
		MaxAwaitRel:                100,
		MaxMsgQueue:                1000,
		DeliverMode:                gmqtt.OnlyOnce,
	}

	// Create a new server
	m.server = gmqtt.NewServer(config)

	log.Infof("Serving server at tcp://[::]:%d", m.cfg.SrvPort)

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", m.cfg.SrvPort))
	if err != nil {
		log.Error(err.Error())
	}

	m.server.AddTCPListenner(ln)

	m.hooks()

	m.server.Run()
}

func (m *Mqtt) NewClient(topic string) (c *Client, err error) {

	uri := fmt.Sprintf("tcp://127.0.0.1:%d", m.cfg.SrvPort)

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
}
