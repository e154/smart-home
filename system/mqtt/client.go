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
	"errors"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"sync"
)

type Client struct {
	mqtt            *Mqtt
	name            string
	subscribersLock *sync.Mutex
	subscribers     map[string]MessageHandler
}

func NewClient(mqtt *Mqtt, name string) *Client {
	return &Client{
		mqtt:            mqtt,
		name:            name,
		subscribersLock: &sync.Mutex{},
		subscribers:     make(map[string]MessageHandler),
	}
}

func (m *Client) Publish(topic string, payload []byte) (err error) {
	m.mqtt.Publish(topic, payload, 0, false)
	return
}

func (m *Client) Subscribe(topic string, handler MessageHandler) (err error) {
	m.subscribersLock.Lock()
	defer m.subscribersLock.Unlock()
	if m.subscribers[topic] != nil {
		delete(m.subscribers, topic)
	}
	if !packets.ValidTopicFilter([]byte(topic)) {
		err = errors.New("invalid topic filter")
		return
	}
	m.subscribers[topic] = handler
	return
}

func (m *Client) Unsubscribe(topic string) {
	m.subscribersLock.Lock()
	defer m.subscribersLock.Unlock()
	if m.subscribers[topic] != nil {
		delete(m.subscribers, topic)
	}
}

func (m *Client) UnsubscribeAll() {
	m.subscribersLock.Lock()
	defer m.subscribersLock.Unlock()
	for topic, _ := range m.subscribers {
		delete(m.subscribers, topic)
	}
}

func (m *Client) OnMsgArrived(ctx context.Context, client gmqtt.Client, msg packets.Message) {
	m.subscribersLock.Lock()
	defer m.subscribersLock.Unlock()
	for topic, handler := range m.subscribers {
		if !packets.TopicMatch([]byte(msg.Topic()), []byte(topic)) {
			continue
		}
		handler(m, Message{
			Dup:      msg.Dup(),
			Qos:      msg.Qos(),
			Retained: msg.Retained(),
			Topic:    msg.Topic(),
			PacketID: msg.PacketID(),
			Payload:  msg.Payload(),
		})
	}
}
