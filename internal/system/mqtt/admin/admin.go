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

package admin

import (
	"context"

	"github.com/e154/smart-home/pkg/logger"

	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/DrmagicE/gmqtt/server"
)

var _ server.Plugin = (*Admin)(nil)

var (
	log = logger.MustGetLogger("mqtt.admin")
)

// Name ...
const Name = "admin"

// New ...
func New() *Admin {
	return &Admin{}
}

// Admin ...
type Admin struct {
	store               *store
	statsReader         server.StatsReader
	publisher           server.Publisher
	clientService       server.ClientService
	subscriptionService server.SubscriptionService
}

// HookWrapper ...
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

// Load ...
func (a *Admin) Load(service server.Server) error {

	a.store = newStore(service.StatsManager(), service.SubscriptionService(), service.ClientService())
	a.statsReader = service.StatsManager()
	a.publisher = service.Publisher()
	a.clientService = service.ClientService()
	a.subscriptionService = service.SubscriptionService()

	log.Info("loaded ...")

	return nil
}

// Unload ...
func (a *Admin) Unload() error {
	log.Info("unloaded ...")
	return nil
}

// Name ...
func (a *Admin) Name() string {
	return Name
}

// OnSessionCreatedWrapper store the client when session created
func (a *Admin) OnSessionCreatedWrapper(pre server.OnSessionCreated) server.OnSessionCreated {
	return func(cs context.Context, client server.Client) {
		pre(cs, client)
		a.store.addClient(client)
	}
}

// OnSessionResumedWrapper refresh the client when session resumed
func (a *Admin) OnSessionResumedWrapper(pre server.OnSessionResumed) server.OnSessionResumed {
	return func(cs context.Context, client server.Client) {
		pre(cs, client)
		a.store.addClient(client)
	}
}

// OnClosedWrapper refresh the client when session resumed
func (a *Admin) OnClosedWrapper(pre server.OnClosed) server.OnClosed {
	return func(cs context.Context, client server.Client, err error) {
		pre(cs, client, err)
		a.store.setClientDisconnected(client.ClientOptions().ClientID)
	}
}

// OnSessionTerminated remove the client when session terminated
func (a *Admin) OnSessionTerminatedWrapper(pre server.OnSessionTerminated) server.OnSessionTerminated {
	return func(cs context.Context, client string, reason server.SessionTerminatedReason) {
		pre(cs, client, reason)
		a.store.removeClient(client)
	}
}

// OnSubscribedWrapper store the subscription
func (a *Admin) OnSubscribedWrapper(pre server.OnSubscribed) server.OnSubscribed {
	return func(cs context.Context, client server.Client, subscription *gmqtt.Subscription) {
		pre(cs, client, subscription)
		a.store.addSubscription(client.ClientOptions().ClientID, subscription)
	}
}

// OnUnsubscribedWrapper remove the subscription
func (a *Admin) OnUnsubscribedWrapper(pre server.OnUnsubscribed) server.OnUnsubscribed {
	return func(cs context.Context, client server.Client, topicName string) {
		pre(cs, client, topicName)
		a.store.removeSubscription(client.ClientOptions().ClientID, topicName)
	}
}

// GetClients ...
func (a *Admin) GetClients(limit, offset uint) (list []*ClientInfo, total uint32, err error) {
	list, total, err = a.store.GetClients(limit, offset)
	return
}

// GetClient ...
func (a *Admin) GetClient(clientId string) (client *ClientInfo, err error) {
	client = a.store.GetClientByID(clientId)
	return
}

// GetSessions ...
func (a *Admin) GetSessions(limit, offset uint) (list []*SessionInfo, total int, err error) {
	list, total, err = a.store.GetSessions(offset, limit)
	return
}

// GetSession ...
func (a *Admin) GetSession(clientId string) (session *SessionInfo, err error) {
	session, err = a.store.GetSessionByID(clientId)
	return
}

// GetSubscriptions ...
func (a *Admin) GetSubscriptions(clientId string, limit, offset uint) (list []*SubscriptionInfo, total int, err error) {
	list, total, err = a.store.GetClientSubscriptions(clientId, offset, limit)
	return
}

// Subscribe ...
func (a *Admin) Subscribe(clientId, topic string, qos int) (err error) {
	if qos < 0 || qos > 2 {
		err = ErrInvalidQos
		return
	}
	if !packets.ValidTopicFilter(true, []byte(topic)) {
		err = ErrInvalidTopicFilter
		return
	}
	if clientId == "" {
		err = ErrInvalidClientID
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
		err = ErrInvalidTopicFilter
		return
	}
	if clientId == "" {
		err = ErrInvalidClientID
		return
	}
	_ = a.subscriptionService.Unsubscribe(clientId, topic)
	return
}

// Publish ...
func (a *Admin) Publish(topic string, qos int, payload []byte, retain bool) (err error) {
	if qos < 0 || qos > 2 {
		err = ErrInvalidQos
		return
	}
	if !packets.ValidTopicFilter(true, []byte(topic)) {
		err = ErrInvalidTopicFilter
		return
	}
	if !packets.ValidUTF8(payload) {
		err = ErrInvalidUtf8String
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
		err = ErrInvalidClientID
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
	result, err = a.store.SearchTopic(query)
	return
}
