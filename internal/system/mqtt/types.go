// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

import "github.com/e154/smart-home/internal/system/mqtt/admin"

// Admin ...
type Admin interface {
	GetClients(limit, offset uint) (list []*admin.ClientInfo, total uint32, err error)
	GetClient(clientId string) (client *admin.ClientInfo, err error)
	GetSessions(limit, offset uint) (list []*admin.SessionInfo, total int, err error)
	GetSession(clientId string) (session *admin.SessionInfo, err error)
	GetSubscriptions(clientId string, limit, offset uint) (list []*admin.SubscriptionInfo, total int, err error)
	Subscribe(clientId, topic string, qos int) (err error)
	Unsubscribe(clientId, topic string) (err error)
	Publish(topic string, qos int, payload []byte, retain bool) (err error)
	CloseClient(clientId string) (err error)
	SearchTopic(query string) (result []*admin.SubscriptionInfo, err error)
}

type MqttServAdmin interface {
	Admin() Admin
}
