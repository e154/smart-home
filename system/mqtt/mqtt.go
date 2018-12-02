package mqtt

import (
	"github.com/op/go-logging"
	"github.com/surgemq/surgemq/service"
	"github.com/e154/smart-home/system/graceful_service"
	"fmt"
)

var (
	log = logging.MustGetLogger("mqtt")
)

type Mqtt struct {
	cfg     *MqttConfig
	server  *service.Server
	clients []*Client
}

func NewMqtt(cfg *MqttConfig,
	graceful *graceful_service.GracefulService) (mqtt *Mqtt) {
	mqtt = &Mqtt{
		cfg: cfg,
	}

	go mqtt.runServer()

	graceful.Subscribe(mqtt)

	return
}

func (m *Mqtt) Shutdown() {
	//if m.server != nil {
	//	m.server.Close()
	//}

	//for _, client := range m.clients {
	//	if client == nil {
	//		continue
	//	}
	//	client.Disconnect()
	//}

	log.Info("Server exiting")
}

func (m *Mqtt) runServer() {

	// Create a new server
	m.server = &service.Server{
		KeepAlive:        m.cfg.SrvKeepAlive,
		ConnectTimeout:   m.cfg.SrvConnectTimeout,
		SessionsProvider: m.cfg.SrvSessionsProvider,
		Authenticator:    m.cfg.SrvAuthenticator,
		TopicsProvider:   m.cfg.SrvTopicsProvider,
	}

	log.Infof("Serving server at tcp://[::]:%d", m.cfg.SrvPort)

	if err := m.server.ListenAndServe(fmt.Sprintf("tcp://0.0.0.0:%d", m.cfg.SrvPort)); err != nil {
		log.Error(err.Error())
	}
}

func (m *Mqtt) NewClient(topic string) (c *Client, err error) {

	uri := fmt.Sprintf("tcp://127.0.0.1:%d", m.cfg.SrvPort)

	log.Infof("new queue client %s topic(%s)", uri, topic)

	if c, err = NewClient(uri, topic); err != nil {
		return
	}

	m.clients = append(m.clients, c)

	return
}
