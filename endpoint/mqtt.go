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

package endpoint

import (
	"context"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/system/mqtt/admin"
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
func (m *MqttEndpoint) GetClientList(ctx context.Context, pagination common.PageParams) (list []*admin.ClientInfo, total uint32, err error) {
	if m.mqtt.Admin() == nil {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	list, total, err = m.mqtt.Admin().GetClients(uint(pagination.Limit), uint(pagination.Offset))
	return
}

// GetClientById ...
func (m *MqttEndpoint) GetClientById(ctx context.Context, clientId string) (client *admin.ClientInfo, err error) {
	if m.mqtt.Admin() == nil {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	client, err = m.mqtt.Admin().GetClient(clientId)
	return
}

// GetSessions ...
func (m *MqttEndpoint) GetSessions(limit, offset uint) (list []*admin.SessionInfo, total int, err error) {
	if m.mqtt.Admin() == nil {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	list, total, err = m.mqtt.Admin().GetSessions(limit, offset)
	return
}

// GetSession ...
func (m *MqttEndpoint) GetSession(clientId string) (session *admin.SessionInfo, err error) {
	if m.mqtt.Admin() == nil {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	session, err = m.mqtt.Admin().GetSession(clientId)
	return
}

// GetSubscriptionList ...
func (m *MqttEndpoint) GetSubscriptionList(ctx context.Context, clientId *string, pagination common.PageParams) (list []*admin.SubscriptionInfo, total int, err error) {
	if m.mqtt.Admin() == nil {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	if clientId == nil {
		err = apperr.ErrBadRequestParams
		return
	}
	list, total, err = m.mqtt.Admin().GetSubscriptions(*clientId, uint(pagination.Limit), uint(pagination.Offset))
	return
}

// Subscribe ...
func (m *MqttEndpoint) Subscribe(clientId, topic string, qos int) (err error) {
	if m.mqtt.Admin() == nil {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	err = m.mqtt.Admin().Subscribe(clientId, topic, qos)
	return
}

// Unsubscribe ...
func (m *MqttEndpoint) Unsubscribe(clientId, topic string) (err error) {
	if m.mqtt.Admin() == nil {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	err = m.mqtt.Admin().Unsubscribe(clientId, topic)
	return
}

// Publish ...
func (m *MqttEndpoint) Publish(topic string, qos int, payload []byte, retain bool) (err error) {
	if m.mqtt.Admin() == nil {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	err = m.mqtt.Admin().Publish(topic, qos, payload, retain)
	return
}

// CloseClient ...
func (m *MqttEndpoint) CloseClient(clientId string) (err error) {
	if m.mqtt.Admin() == nil {
		err = apperr.ErrMqttServerNoWorked
		return
	}
	err = m.mqtt.Admin().CloseClient(clientId)
	return
}

// SearchTopic ...
func (m *MqttEndpoint) SearchTopic(query string, limit, offset int) (result []*admin.SubscriptionInfo, total int64, err error) {
	if m.mqtt.Admin() == nil {
		err = apperr.ErrMqttServerNoWorked
		return
	}

	if result, err = m.mqtt.Admin().SearchTopic(query); err != nil {
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
		sub := &admin.SubscriptionInfo{
			Name: query,
		}
		result = append(result, sub)
	}

	return
}
