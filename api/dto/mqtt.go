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
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/mqtt/admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mqtt ...
type Mqtt struct{}

// NewMqttDto ...
func NewMqttDto() Mqtt {
	return Mqtt{}
}

// GetClientById ...
func (r Mqtt) GetClientById(from *admin.ClientInfo) (client *api.Client) {
	if from == nil {
		return
	}
	client = &api.Client{
		ClientId:             from.ClientID,
		Username:             from.Username,
		KeepAlive:            uint32(from.KeepAlive),
		Version:              from.Version,
		WillRetain:           from.WillRetain,
		WillQos:              uint32(from.WillQos),
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
		ConnectedAt:          timestamppb.New(from.ConnectedAt),
	}
	if from.DisconnectedAt != nil {
		client.DisconnectedAt = timestamppb.New(*from.DisconnectedAt)
	}
	return
}

// ToListResult ...
func (r Mqtt) ToListResult(list []*admin.ClientInfo, total uint64, pagination common.PageParams) *api.GetClientListResult {

	items := make([]*api.Client, 0, len(list))

	for _, i := range list {
		items = append(items, r.GetClientById(i))
	}

	return &api.GetClientListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

// GetSubscriptiontById ...
func (r Mqtt) GetSubscriptiontById(from *admin.SubscriptionInfo) (client *api.Subscription) {
	if from == nil {
		return
	}
	client = &api.Subscription{
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
func (r Mqtt) GetSubscriptionList(list []*admin.SubscriptionInfo, total uint64, pagination common.PageParams) *api.GetSubscriptionListResult {

	items := make([]*api.Subscription, 0, len(list))

	for _, i := range list {
		items = append(items, r.GetSubscriptiontById(i))
	}

	return &api.GetSubscriptionListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}
