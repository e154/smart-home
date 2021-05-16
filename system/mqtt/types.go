// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"github.com/DrmagicE/gmqtt/config"
	"github.com/DrmagicE/gmqtt/retained"
	"github.com/DrmagicE/gmqtt/server"
)

// Admin ...
type Admin interface {
	GetClients(_page, _pageSize uint32) (clients, total interface{}, err error)
	//GetClient(clientId string) (client *management.ClientInfo, err error)
	//GetSessions(limit, offset int) (list []*management.SessionInfo, total int, err error)
	//GetSession(clientId string) (session *management.SessionInfo, err error)
	//GetSubscriptions(clientId string, limit, offset int) (list []*management.SubscriptionInfo, total int, err error)
	//Subscribe(clientId, topic string, qos int) (err error)
	//Unsubscribe(clientId, topic string) (err error)
	//Publish(topic string, qos int, payload []byte, retain bool) (err error)
	//CloseClient(clientId string) (err error)
	//SearchTopic(query string) (result []*management.SubscriptionInfo, err error)
}

type MqttCli interface {
	Publish(topic string, payload []byte) error
	Subscribe(topic string, handler MessageHandler) error
	Unsubscribe(topic string)
	UnsubscribeAll()
	OnMsgArrived(ctx context.Context, client server.Client, req *server.MsgArrivedRequest)
}

type MqttServ interface {
	Shutdown() error
	Start()
	Publish(topic string, payload []byte, qos uint8, retain bool) error
	NewClient(name string) MqttCli
	RemoveClient(name string)
	Admin() Admin
}

// GMqttServer ...
type GMqttServer interface {
	// ...
	Run() error
	// ...
	Stop(ctx context.Context) error
	// ...
	Init(opts ...server.Options) error
	// SubscriptionStore returns the subscription.Store.
	SubscriptionService() server.SubscriptionService
	// RetainedStore returns the retained.Store.
	RetainedStore() retained.Store
	// Publisher returns the Publisher
	Publisher() server.Publisher
	// client return the ClientService
	ClientService() server.ClientService
	// GetConfig returns the config of the server
	GetConfig() config.Config
	// StatsManager returns StatsReader
	StatsManager() server.StatsReader
}

// Message ...
type Message struct {
	Dup      bool
	Qos      uint8
	Retained bool
	Topic    string
	PacketID uint16
	Payload  []byte
}

// MessageHandler ...
type MessageHandler func(MqttCli, Message)
