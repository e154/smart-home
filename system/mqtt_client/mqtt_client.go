package mqtt_client

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
	"time"
)

var (
	log = logging.MustGetLogger("mqtt_client")
)

type Client struct {
	cfg        *Config
	client     MQTT.Client
	subscribes map[string]Subscribe
}

func NewClient(cfg *Config) (client *Client, err error) {

	log.Infof("new queue client(%s) uri(%s)", cfg.ClientID, cfg.Broker)

	client = &Client{
		cfg:        cfg,
		subscribes: make(map[string]Subscribe),
	}

	opts := MQTT.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID).
		SetKeepAlive(time.Duration(cfg.KeepAlive) * time.Second).
		SetPingTimeout(time.Duration(cfg.PingTimeout) * time.Second).
		SetConnectTimeout(time.Duration(cfg.ConnectTimeout) * time.Second).
		SetCleanSession(cfg.CleanSession).
		SetOnConnectHandler(client.onConnect).
		SetConnectionLostHandler(client.onConnectionLostHandler)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}

	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}

	client.client = MQTT.NewClient(opts)

	return
}

func (c *Client) Connect() (err error) {

	log.Infof("Connect to server %s", c.cfg.Broker)

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
	}

	return
}

func (c *Client) Disconnect() {
	if c.client == nil {
		return
	}

	c.UnsubscribeAll()
	c.subscribes = make(map[string]Subscribe)
	c.client.Disconnect(250)
	c.client = nil
}

func (c *Client) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) (err error) {

	if _, ok := c.subscribes[topic]; !ok {
		c.subscribes[topic] = Subscribe{
			Qos:      qos,
			Callback: callback,
		}
	}

	if token := c.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
		return token.Error()
	}
	return
}

func (c *Client) Unsubscribe(topic string) (err error) {

	if token := c.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
		return token.Error()
	}
	return
}

func (c *Client) UnsubscribeAll() {

	for topic, _ := range c.subscribes {
		if token := c.client.Unsubscribe(topic); token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
	c.subscribes = make(map[string]Subscribe)
}

func (c *Client) Publish(topic string, payload interface{}) (err error) {
	if c.client != nil && (c.client.IsConnected()) {
		c.client.Publish(topic, c.cfg.Qos, false, payload)
	}
	return
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnectionOpen()
}

func (c *Client) onConnectionLostHandler(client MQTT.Client, e error) {

	log.Debug("connection lost...")

	for topic, _ := range c.subscribes {
		if token := c.client.Unsubscribe(topic); token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
}

func (c *Client) onConnect(client MQTT.Client) {

	log.Debug("connected...")

	for topic, subscribe := range c.subscribes {
		if token := c.client.Subscribe(topic, subscribe.Qos, subscribe.Callback); token.Wait() && token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
}
