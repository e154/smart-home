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

package admin

import (
	"context"
	"errors"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/DrmagicE/gmqtt/server"
	"github.com/e154/smart-home/common"
)

var _ server.Plugin = (*Admin)(nil)

var (
	log = common.MustGetLogger("mqtt.admin")
)

const Name = "admin"

func New() *Admin {
	return &Admin{}
}

type Admin struct {
	monitor             *monitor
	statsReader         server.StatsReader
	publisher           server.Publisher
	clientService       server.ClientService
	subscriptionService server.SubscriptionService
}

func (a *Admin) HookWrapper() server.HookWrapper {
	return server.HookWrapper{
		OnSessionCreatedWrapper:    a.OnSessionCreatedWrapper,
		OnSessionResumedWrapper:    a.OnSessionResumedWrapper,
		OnClosedWrapper:            a.OnClosedWrapper,
		OnSessionTerminatedWrapper: a.OnSessionTerminatedWrapper,
		OnSubscribedWrapper:        a.OnSubscribedWrapper,
		OnUnsubscribedWrapper:      a.OnUnsubscribedWrapper,
	}
}

func (a *Admin) Load(service server.Server) error {

	a.monitor = newMonitor(service.SubscriptionService(), service.GetConfig(), service.StatsManager())
	a.statsReader = service.StatsManager()
	a.publisher = service.Publisher()
	a.clientService = service.ClientService()
	a.subscriptionService = service.SubscriptionService()

	log.Info("loaded ...")

	return nil
}

func (a *Admin) Unload() error {
	log.Info("unloaded ...")
	return nil
}

func (a *Admin) Name() string {
	return Name
}

// OnSessionCreatedWrapper store the client when session created
func (a *Admin) OnSessionCreatedWrapper(created server.OnSessionCreated) server.OnSessionCreated {
	return func(cs context.Context, client server.Client) {
		a.monitor.addClient(client)
		created(cs, client)
	}
}

// OnSessionResumedWrapper refresh the client when session resumed
func (a *Admin) OnSessionResumedWrapper(resumed server.OnSessionResumed) server.OnSessionResumed {
	return func(cs context.Context, client server.Client) {
		a.monitor.addClient(client)
		resumed(cs, client)
	}
}

// OnClosedWrapper refresh the client when session resumed
func (a *Admin) OnClosedWrapper(pre server.OnClosed) server.OnClosed {
	return func(cs context.Context, client server.Client, err error) {
		a.monitor.setClientDisconnected(client.ClientOptions().ClientID)
		pre(cs, client, err)
	}
}

// OnSessionTerminated remove the client when session terminated
func (a *Admin) OnSessionTerminatedWrapper(terminated server.OnSessionTerminated) server.OnSessionTerminated {
	return func(cs context.Context, client string, reason server.SessionTerminatedReason) {
		a.monitor.deleteClient(client)
		a.monitor.deleteClientSubscriptions(client)
		terminated(cs, client, reason)
	}
}

// OnSubscribedWrapper store the subscription
func (a *Admin) OnSubscribedWrapper(subscribed server.OnSubscribed) server.OnSubscribed {
	return func(cs context.Context, client server.Client, subscription *gmqtt.Subscription) {
		a.monitor.addSubscription(client.ClientOptions().ClientID, subscription)
		subscribed(cs, client, subscription)
	}
}

// OnUnsubscribedWrapper remove the subscription
func (a *Admin) OnUnsubscribedWrapper(unsubscribe server.OnUnsubscribed) server.OnUnsubscribed {
	return func(cs context.Context, client server.Client, topicName string) {
		a.monitor.deleteSubscription(client.ClientOptions().ClientID, topicName)
		unsubscribe(cs, client, topicName)
	}
}

// GetClients ...
func (a *Admin) GetClients(limit, offset int) (list []*ClientInfo, total int, err error) {
	list, total, err = a.monitor.GetClients(offset, limit)
	return
}

// GetClient ...
func (a *Admin) GetClient(clientId string) (client *ClientInfo, err error) {
	client, err = a.monitor.GetClientByID(clientId)
	return
}

// GetSessions ...
func (a *Admin) GetSessions(limit, offset int) (list []*SessionInfo, total int, err error) {
	list, total, err = a.monitor.GetSessions(offset, limit)
	return
}

// GetSession ...
func (a *Admin) GetSession(clientId string) (session *SessionInfo, err error) {
	session, err = a.monitor.GetSessionByID(clientId)
	return
}

// GetSubscriptions ...
func (a *Admin) GetSubscriptions(clientId string, limit, offset int) (list []*SubscriptionInfo, total int, err error) {
	list, total, err = a.monitor.GetClientSubscriptions(clientId, offset, limit)
	return
}

// Subscribe ...
func (a *Admin) Subscribe(clientId, topic string, qos int) (err error) {
	if qos < 0 || qos > 2 {
		err = errors.New("invalid Qos")
		return
	}
	if !packets.ValidTopicFilter(true, []byte(topic)) {
		err = errors.New("invalid topic filter")
		return
	}
	if clientId == "" {
		err = errors.New("invalid clientID")
		return
	}
	_, err = a.subscriptionService.Subscribe(clientId, &gmqtt.Subscription{
		TopicFilter: topic,
		QoS:         uint8(qos),
	})
	return
}

// Unsubscribe ...
func (a *Admin) Unsubscribe(clientId, topic string) (err error) {
	if !packets.ValidTopicFilter(true, []byte(topic)) {
		err = errors.New("invalid topic filter")
		return
	}
	if clientId == "" {
		err = errors.New("invalid clientID")
		return
	}
	a.subscriptionService.Unsubscribe(clientId, topic)
	return
}

// Publish ...
func (a *Admin) Publish(topic string, qos int, payload []byte, retain bool) (err error) {
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
	a.publisher.Publish(&gmqtt.Message{
		QoS:      uint8(qos),
		Retained: retain,
		Topic:    topic,
		Payload:  payload,
	})
	return
}

// CloseClient ...
func (a *Admin) CloseClient(clientId string) (err error) {
	if clientId == "" {
		err = errors.New("invalid clientID")
		return
	}
	client := a.clientService.GetClient(clientId)
	if client != nil {
		client.Close()
	}

	return
}

// SearchTopic ...
func (a *Admin) SearchTopic(query string) (result []*SubscriptionInfo, err error) {
	result, err = a.monitor.SearchTopic(query)
	return
}
