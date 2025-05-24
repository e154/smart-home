// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package dto

import (
	"github.com/e154/smart-home/internal/api/stub"
	admin2 "github.com/e154/smart-home/internal/system/mqtt/admin"
)

// Mqtt ...
type Mqtt struct{}

// NewMqttDto ...
func NewMqttDto() Mqtt {
	return Mqtt{}
}

// GetClientById ...
func (r Mqtt) GetClientById(from *admin2.ClientInfo) (client *stub.ApiClient) {
	if from == nil {
		return
	}
	client = &stub.ApiClient{
		ClientId:             from.ClientID,
		Username:             from.Username,
		KeepAlive:            from.KeepAlive,
		Version:              from.Version,
		WillRetain:           from.WillRetain,
		WillQos:              from.WillQos,
		WillTopic:            from.WillTopic,
		WillPayload:          from.WillPayload,
		RemoteAddr:           from.RemoteAddr,
		LocalAddr:            from.LocalAddr,
		SubscriptionsCurrent: from.SubscriptionsCurrent,
		SubscriptionsTotal:   from.SubscriptionsTotal,
		PacketsReceivedBytes: from.PacketsReceivedBytes,
		PacketsReceivedNums:  from.PacketsReceivedNums,
		PacketsSendBytes:     from.PacketsSendBytes,
		PacketsSendNums:      from.PacketsSendNums,
		MessageDropped:       from.MessageDropped,
		InflightLen:          from.InflightLen,
		QueueLen:             from.QueueLen,
		ConnectedAt:          from.ConnectedAt,
		DisconnectedAt:       from.DisconnectedAt,
	}

	return
}

// ToListResult ...
func (r Mqtt) ToListResult(list []*admin2.ClientInfo) []*stub.ApiClient {

	items := make([]*stub.ApiClient, 0, len(list))

	for _, i := range list {
		items = append(items, r.GetClientById(i))
	}

	return items
}

// GetSubscriptiontById ...
func (r Mqtt) GetSubscriptiontById(from *admin2.SubscriptionInfo) (client *stub.ApiSubscription) {
	if from == nil {
		return
	}
	client = &stub.ApiSubscription{
		Id:                from.Id,
		ClientId:          from.ClientID,
		TopicName:         from.TopicName,
		Name:              from.Name,
		Qos:               from.Qos,
		NoLocal:           from.NoLocal,
		RetainAsPublished: from.RetainAsPublished,
		RetainHandling:    from.RetainHandling,
	}
	return
}

// GetSubscriptionList ...
func (r Mqtt) GetSubscriptionList(list []*admin2.SubscriptionInfo) []*stub.ApiSubscription {

	items := make([]*stub.ApiSubscription, 0, len(list))

	for _, i := range list {
		items = append(items, r.GetSubscriptiontById(i))
	}

	return items
}
