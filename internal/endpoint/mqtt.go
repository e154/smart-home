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

package endpoint

import (
	"context"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/internal/system/mqtt"
	admin2 "github.com/e154/smart-home/internal/system/mqtt/admin"
	"github.com/e154/smart-home/pkg/apperr"
)

// MqttEndpoint ...
type MqttEndpoint struct {
	*CommonEndpoint
}

// NewMqttEndpoint ...
func NewMqttEndpoint(common *CommonEndpoint) *MqttEndpoint {
	return &MqttEndpoint{
		CommonEndpoint: common,
	}
}

// GetClientList ...
func (m *MqttEndpoint) GetClientList(ctx context.Context, pagination common.PageParams) (list []*admin2.ClientInfo, total uint32, err error) {
	admin, ok := m.mqtt.(mqtt.MqttServAdmin)
	if !ok {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	list, total, err = admin.Admin().GetClients(uint(pagination.Limit), uint(pagination.Offset))
	return
}

// GetClientById ...
func (m *MqttEndpoint) GetClientById(ctx context.Context, clientId string) (client *admin2.ClientInfo, err error) {
	admin, ok := m.mqtt.(mqtt.MqttServAdmin)
	if !ok {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	client, err = admin.Admin().GetClient(clientId)
	return
}

// GetSessions ...
func (m *MqttEndpoint) GetSessions(limit, offset uint) (list []*admin2.SessionInfo, total int, err error) {
	admin, ok := m.mqtt.(mqtt.MqttServAdmin)
	if !ok {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	list, total, err = admin.Admin().GetSessions(limit, offset)
	return
}

// GetSession ...
func (m *MqttEndpoint) GetSession(clientId string) (session *admin2.SessionInfo, err error) {
	admin, ok := m.mqtt.(mqtt.MqttServAdmin)
	if !ok {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	session, err = admin.Admin().GetSession(clientId)
	return
}

// GetSubscriptionList ...
func (m *MqttEndpoint) GetSubscriptionList(ctx context.Context, clientId *string, pagination common.PageParams) (list []*admin2.SubscriptionInfo, total int, err error) {
	admin, ok := m.mqtt.(mqtt.MqttServAdmin)
	if !ok {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	if clientId == nil {
		err = apperr.ErrBadRequestParams
		return
	}
	list, total, err = admin.Admin().GetSubscriptions(*clientId, uint(pagination.Limit), uint(pagination.Offset))
	return
}

// Subscribe ...
func (m *MqttEndpoint) Subscribe(clientId, topic string, qos int) (err error) {
	admin, ok := m.mqtt.(mqtt.MqttServAdmin)
	if !ok {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	err = admin.Admin().Subscribe(clientId, topic, qos)
	return
}

// Unsubscribe ...
func (m *MqttEndpoint) Unsubscribe(clientId, topic string) (err error) {
	admin, ok := m.mqtt.(mqtt.MqttServAdmin)
	if !ok {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	err = admin.Admin().Unsubscribe(clientId, topic)
	return
}

// Publish ...
func (m *MqttEndpoint) Publish(topic string, qos int, payload []byte, retain bool) (err error) {
	admin, ok := m.mqtt.(mqtt.MqttServAdmin)
	if !ok {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	err = admin.Admin().Publish(topic, qos, payload, retain)
	return
}

// CloseClient ...
func (m *MqttEndpoint) CloseClient(clientId string) (err error) {
	admin, ok := m.mqtt.(mqtt.MqttServAdmin)
	if !ok {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	err = admin.Admin().CloseClient(clientId)
	return
}

// SearchTopic ...
func (m *MqttEndpoint) SearchTopic(query string, limit, offset int) (result []*admin2.SubscriptionInfo, total int64, err error) {
	admin, ok := m.mqtt.(mqtt.MqttServAdmin)
	if !ok {
		err = apperr.ErrMqttServerNoWorked
		return
	}

	if result, err = admin.Admin().SearchTopic(query); err != nil {
		return
	}

	// add custom text as topic
	var exist bool
	for _, sub := range result {
		if sub.Name == query {
			exist = true
		}
	}

	if !exist {
		sub := &admin2.SubscriptionInfo{
			Name: query,
		}
		result = append(result, sub)
	}

	return
}
