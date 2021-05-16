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
	"github.com/DrmagicE/gmqtt/server"
	"time"
)

// ClientInfo represents the client information
type ClientInfo struct {
	ClientID       string     `json:"client_id"`
	Username       string     `json:"username"`
	KeepAlive      uint16     `json:"keep_alive"`
	CleanSession   bool       `json:"clean_session"`
	WillFlag       bool       `json:"will_flag"`
	WillRetain     bool       `json:"will_retain"`
	WillQos        uint8      `json:"will_qos"`
	WillTopic      string     `json:"will_topic"`
	WillPayload    string     `json:"will_payload"`
	RemoteAddr     string     `json:"remote_addr"`
	LocalAddr      string     `json:"local_addr"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at"`
}

func newClientInfo(client server.Client) *ClientInfo {
	sessionInfo := client.SessionInfo()
	optsReader := client.ClientOptions()
	conn := client.Connection()
	rs := &ClientInfo{
		ClientID:  optsReader.ClientID,
		Username:  optsReader.Username,
		KeepAlive: optsReader.KeepAlive,
		//CleanSession:   optsReader.CleanSession(),
		//WillFlag:       optsReader.WillFlag(),
		RemoteAddr:  conn.RemoteAddr().String(),
		LocalAddr:   conn.LocalAddr().String(),
		ConnectedAt: client.ConnectedAt(),
		//DisconnectedAt: client.DisconnectedAt(),
	}
	if sessionInfo.Will != nil {
		rs.WillRetain = sessionInfo.Will.Retained
		rs.WillQos = sessionInfo.Will.QoS
		rs.WillTopic = sessionInfo.Will.Topic
		rs.WillPayload = string(sessionInfo.Will.Payload)
	}
	return rs
}
