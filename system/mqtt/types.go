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
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/retained"
	"github.com/DrmagicE/gmqtt/subscription"
	"github.com/e154/smart-home/system/mqtt/management"
)

type IManagement interface {
	GetClients(limit, offset int) (list []*management.ClientInfo, total int, err error)
	GetClient(clientId string) (client *management.ClientInfo, err error)
	GetSessions(limit, offset int) (list []*management.SessionInfo, total int, err error)
	GetSession(clientId string) (session *management.SessionInfo, err error)
	GetSubscriptions(clientId string, limit, offset int) (list []*management.SubscriptionInfo, total int, err error)
	Subscribe(clientId, topic string, qos int) (err error)
	Unsubscribe(clientId, topic string) (err error)
	Publish(topic string, qos int, payload []byte, retain bool) (err error)
	CloseClient(clientId string) (err error)
	SearchTopic(query string) (result []*management.SubscriptionInfo, err error)
}

type IMQTT interface {
	Run()
	Stop(ctx context.Context) error
	// SubscriptionStore returns the subscription.Store.
	SubscriptionStore() subscription.Store
	// RetainedStore returns the retained.Store.
	RetainedStore() retained.Store
	// PublishService returns the PublishService
	PublishService() gmqtt.PublishService
	// Client return the client specified by clientID.
	Client(clientID string) gmqtt.Client
	// GetConfig returns the config of the server
	GetConfig() gmqtt.Config
	// GetStatsManager returns StatsManager
	GetStatsManager() gmqtt.StatsManager
}

type Message struct {
	Dup       bool
	Qos       uint8
	Retained  bool
	Topic     string
	PacketID  uint16
	Payload   []byte
}
type MessageHandler func(*Client, Message)
