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

package admin

import (
	"time"

	"github.com/DrmagicE/gmqtt/server"
)

// ClientInfo represents the client information
type ClientInfo struct {
	ClientID             string     `json:"client_id"`
	Username             string     `json:"username"`
	KeepAlive            uint16     `json:"keep_alive"`
	Version              int32      `json:"version"`
	WillRetain           bool       `json:"will_retain"`
	WillQos              uint8      `json:"will_qos"`
	WillTopic            string     `json:"will_topic"`
	WillPayload          string     `json:"will_payload"`
	RemoteAddr           string     `json:"remote_addr"`
	LocalAddr            string     `json:"local_addr"`
	SubscriptionsCurrent uint32     `json:"subscriptions_current"`
	SubscriptionsTotal   uint32     `json:"subscriptions_total"`
	PacketsReceivedBytes uint64     `json:"packets_received_bytes"`
	PacketsReceivedNums  uint64     `json:"packets_received_nums"`
	PacketsSendBytes     uint64     `json:"packets_send_bytes"`
	PacketsSendNums      uint64     `json:"packets_send_nums"`
	MessageDropped       uint64     `json:"message_dropped"`
	InflightLen          uint32     `json:"inflight_len"`
	QueueLen             uint32     `json:"queue_len"`
	ConnectedAt          time.Time  `json:"connected_at"`
	DisconnectedAt       *time.Time `json:"disconnected_at"`
}

func newClientInfo(client server.Client) *ClientInfo {
	sessionInfo := client.SessionInfo()
	optsReader := client.ClientOptions()
	conn := client.Connection()
	rs := &ClientInfo{
		ClientID:    optsReader.ClientID,
		Username:    optsReader.Username,
		KeepAlive:   optsReader.KeepAlive,
		RemoteAddr:  conn.RemoteAddr().String(),
		LocalAddr:   conn.LocalAddr().String(),
		ConnectedAt: client.ConnectedAt(),
		Version:     int32(client.Version()),
	}
	if sessionInfo.Will != nil {
		rs.WillRetain = sessionInfo.Will.Retained
		rs.WillQos = sessionInfo.Will.QoS
		rs.WillTopic = sessionInfo.Will.Topic
		rs.WillPayload = string(sessionInfo.Will.Payload)
	}
	return rs
}
