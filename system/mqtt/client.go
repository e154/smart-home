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
	"sync"

	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/DrmagicE/gmqtt/server"
)

// client ...
type client struct {
	mqtt            *Mqtt
	name            string
	subscribersLock *sync.Mutex
	subscribers     map[string]MessageHandler
}

// NewClient ...
func NewClient(mqtt *Mqtt, name string) MqttCli {
	return &client{
		mqtt:            mqtt,
		name:            name,
		subscribersLock: &sync.Mutex{},
		subscribers:     make(map[string]MessageHandler),
	}
}

// Publish ...
func (m *client) Publish(topic string, payload []byte) (err error) {
	_ = m.mqtt.Publish(topic, payload, 1, true)
	return
}

// Subscribe ...
func (m *client) Subscribe(topic string, handler MessageHandler) (err error) {
	log.Infof("client %s, subscribed to topic '%s'", m.name, topic)
	m.subscribersLock.Lock()
	defer m.subscribersLock.Unlock()
	if m.subscribers[topic] != nil {
		log.Infof("DELETE topic %s", topic)
		delete(m.subscribers, topic)
	}
	if !packets.ValidTopicFilter(true, []byte(topic)) {
		err = ErrInvalidTopicFilter
		return
	}
	m.subscribers[topic] = handler
	return
}

// Unsubscribe ...
func (m *client) Unsubscribe(topic string) {
	m.subscribersLock.Lock()
	defer m.subscribersLock.Unlock()
	if m.subscribers[topic] != nil {
		delete(m.subscribers, topic)
	}
}

// UnsubscribeAll ...
func (m *client) UnsubscribeAll() {
	m.subscribersLock.Lock()
	defer m.subscribersLock.Unlock()
	for topic := range m.subscribers {
		delete(m.subscribers, topic)
	}
}

// OnMsgArrived ...
func (m *client) OnMsgArrived(ctx context.Context, client server.Client, req *server.MsgArrivedRequest) {
	m.subscribersLock.Lock()
	defer m.subscribersLock.Unlock()
	for topic, handler := range m.subscribers {
		msg := req.Message
		if !packets.TopicMatch([]byte(msg.Topic), []byte(topic)) {
			continue
		}
		go handler(m, Message{
			Dup:      msg.Dup,
			Qos:      msg.QoS,
			Retained: msg.Retained,
			Topic:    msg.Topic,
			PacketID: msg.PacketID,
			Payload:  msg.Payload,
		})
	}
}
