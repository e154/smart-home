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

package mqtt_client

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
	"strings"
	"sync"
	"time"
)

var (
	log = logging.MustGetLogger("mqtt_client")
)

type Client struct {
	cfg *Config
	sync.Mutex
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

	c.Lock()
	defer c.Unlock()

	log.Infof("Connect to server %s", c.cfg.Broker)

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
	}

	return
}

func (c *Client) Disconnect() {

	c.Lock()
	if c.client == nil {
		c.Unlock()
		return
	}
	c.Unlock()

	c.UnsubscribeAll()

	c.Lock()
	c.client.Disconnect(250)
	//c.client = nil
	c.Unlock()
}

func (c *Client) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) (err error) {

	c.Lock()
	defer c.Unlock()

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

	c.Lock()
	defer c.Unlock()

	if token := c.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
		return token.Error()
	}
	return
}

func (c *Client) UnsubscribeAll() {
	c.Lock()
	defer c.Unlock()

	for topic, _ := range c.subscribes {
		if token := c.client.Unsubscribe(topic); token.Error() != nil {
			log.Error(token.Error().Error())
		}
		delete(c.subscribes, topic)
	}
}

func (c *Client) Publish(topic string, payload interface{}) (err error) {
	c.Lock()
	defer c.Unlock()

	if c.client != nil && (c.client.IsConnected()) {
		c.client.Publish(topic, c.cfg.Qos, false, payload)
	}
	return
}

func (c *Client) IsConnected() bool {
	c.Lock()
	defer c.Unlock()

	return c.client.IsConnectionOpen()
}

func (c *Client) onConnectionLostHandler(client MQTT.Client, e error) {

	c.Lock()
	defer c.Unlock()

	log.Debug("connection lost...")

	for topic, _ := range c.subscribes {
		if token := c.client.Unsubscribe(topic); token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
}

func (c *Client) onConnect(client MQTT.Client) {

	c.Lock()
	defer c.Unlock()

	log.Debug("connected...")

	for topic, subscribe := range c.subscribes {
		if token := c.client.Subscribe(topic, subscribe.Qos, subscribe.Callback); token.Wait() && token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
}

func ClientIdGen(args ...interface{}) string {
	var b strings.Builder
	b.WriteString("smarthome")
	for _, n := range args {
		fmt.Fprintf(&b, "_%v", n)
	}
	return b.String()
}
