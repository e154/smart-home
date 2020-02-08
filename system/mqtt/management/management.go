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

package management

import (
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/pkg/errors"
)

type Management struct {
	monitor *monitor
	service gmqtt.ServerService
}

func NewManagement() *Management {
	return &Management{
		monitor: newMonitor(),
	}
}

func (m *Management) Load(service gmqtt.ServerService) error {
	m.service = service
	m.monitor.config = service.GetConfig()
	return nil
}

func (m *Management) Unload() error {
	return nil
}

func (m *Management) HookWrapper() gmqtt.HookWrapper {
	return gmqtt.HookWrapper{
		OnSessionCreatedWrapper:    m.OnSessionCreatedWrapper,
		OnSessionResumedWrapper:    m.OnSessionResumedWrapper,
		OnSessionTerminatedWrapper: m.OnSessionTerminatedWrapper,
		OnSubscribedWrapper:        m.OnSubscribedWrapper,
		OnUnsubscribedWrapper:      m.OnUnsubscribedWrapper,
	}
}

func (m *Management) Name() string {
	return "management"
}

// OnSessionCreatedWrapper store the client when session created
func (m *Management) OnSessionCreatedWrapper(created gmqtt.OnSessionCreated) gmqtt.OnSessionCreated {
	return func(cs gmqtt.ChainStore, client gmqtt.Client) {
		m.monitor.addClient(client)
		created(cs, client)
	}
}

// OnSessionResumedWrapper refresh the client when session resumed
func (m *Management) OnSessionResumedWrapper(resumed gmqtt.OnSessionResumed) gmqtt.OnSessionResumed {
	return func(cs gmqtt.ChainStore, client gmqtt.Client) {
		m.monitor.addClient(client)
		resumed(cs, client)
	}
}

// OnSessionTerminated remove the client when session terminated
func (m *Management) OnSessionTerminatedWrapper(terminated gmqtt.OnSessionTerminated) gmqtt.OnSessionTerminated {
	return func(cs gmqtt.ChainStore, client gmqtt.Client, reason gmqtt.SessionTerminatedReason) {
		m.monitor.deleteClient(client.OptionsReader().ClientID())
		m.monitor.deleteClientSubscriptions(client.OptionsReader().ClientID())
		terminated(cs, client, reason)
	}
}

// OnSubscribedWrapper store the subscription
func (m *Management) OnSubscribedWrapper(subscribed gmqtt.OnSubscribed) gmqtt.OnSubscribed {
	return func(cs gmqtt.ChainStore, client gmqtt.Client, topic packets.Topic) {
		m.monitor.addSubscription(client.OptionsReader().ClientID(), topic)
		subscribed(cs, client, topic)
	}
}

// OnUnsubscribedWrapper remove the subscription
func (m *Management) OnUnsubscribedWrapper(unsubscribe gmqtt.OnUnsubscribed) gmqtt.OnUnsubscribed {
	return func(cs gmqtt.ChainStore, client gmqtt.Client, topicName string) {
		m.monitor.deleteSubscription(client.OptionsReader().ClientID(), topicName)
		unsubscribe(cs, client, topicName)
	}
}

// GetClients
func (m *Management) GetClients(limit, offset int) (list []*ClientInfo, total int, err error) {
	list, total, err = m.monitor.GetClients(offset, limit)
	return
}

//GetClient
func (m *Management) GetClient(clientId string) (client *ClientInfo, err error) {
	client, err = m.monitor.GetClientByID(clientId)
	return
}

//GetSessions
func (m *Management) GetSessions(limit, offset int) (list []*SessionInfo, total int, err error) {
	list, total, err = m.monitor.GetSessions(offset, limit)
	return
}

//GetSession
func (m *Management) GetSession(clientId string) (session *SessionInfo, err error) {
	session, err = m.monitor.GetSessionByID(clientId)
	return
}

//GetSubscriptions
func (m *Management) GetSubscriptions(clientId string, limit, offset int) (list []*SubscriptionInfo, total int, err error) {
	list, total, err = m.monitor.GetClientSubscriptions(clientId, offset, limit)
	return
}

//Subscribe
func (m *Management) Subscribe(clientId, topic string, qos int) (err error) {
	if qos < 0 || qos > 2 {
		err = errors.New("invalid Qos")
		return
	}
	if !packets.ValidTopicFilter([]byte(topic)) {
		err = errors.New("invalid topic filter")
		return
	}
	if clientId == "" {
		err = errors.New("invalid clientID")
		return
	}
	m.service.Subscribe(clientId, []packets.Topic{
		{
			Qos:  uint8(qos),
			Name: topic,
		},
	})
	return
}

//Unsubscribe
func (m *Management) Unsubscribe(clientId, topic string) (err error) {
	if !packets.ValidTopicFilter([]byte(topic)) {
		err = errors.New("invalid topic filter")
		return
	}
	if clientId == "" {
		err = errors.New("invalid clientID")
		return
	}
	m.service.UnSubscribe(clientId, []string{topic})
	return
}

//Publish
func (m *Management) Publish(topic string, qos int, payload []byte, retain bool) (err error) {
	if qos < 0 || qos > 2 {
		err = errors.New("invalid Qos")
		return
	}
	if !packets.ValidTopicFilter([]byte(topic)) {
		err = errors.New("invalid topic filter")
		return
	}
	if !packets.ValidUTF8([]byte(payload)) {
		err = errors.New("invalid utf-8 string")
		return
	}
	m.service.Publish(topic, []byte(payload), uint8(qos), retain)
	return
}

//CloseClient
func (m *Management) CloseClient(clientId string) (err error) {
	if clientId == "" {
		err = errors.New("invalid clientID")
		return
	}
	client := m.service.Client(clientId)
	if client != nil {
		client.Close()
	}

	return
}

//SearchTopic
func (m *Management) SearchTopic(query string) (result []*SubscriptionInfo, err error) {
	result, err = m.monitor.SearchTopic(query)
	return
}
