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

package management

import (
	"container/list"
	"errors"
	"time"
)

// ClientInfo represents the client information
type ClientInfo struct {
	ClientID       string    `json:"client_id"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	KeepAlive      uint16    `json:"keep_alive"`
	CleanSession   bool      `json:"clean_session"`
	WillFlag       bool      `json:"will_flag"`
	WillRetain     bool      `json:"will_retain"`
	WillQos        uint8     `json:"will_qos"`
	WillTopic      string    `json:"will_topic"`
	WillPayload    string    `json:"will_payload"`
	RemoteAddr     string    `json:"remote_addr"`
	LocalAddr      string    `json:"local_addr"`
	ConnectedAt    time.Time `json:"connected_at"`
	DisconnectedAt time.Time `json:"disconnected_at"`
}

// SessionInfo represents the session information
type SessionInfo struct {
	ClientID          string    `json:"client_id"`
	Status            string    `json:"status"`
	CleanSession      bool      `json:"clean_session"`
	Subscriptions     int64     `json:"subscriptions"`
	MaxInflight       int       `json:"max_inflight"`
	InflightLen       int64     `json:"inflight_len"`
	MaxMsgQueue       int       `json:"max_msg_queue"`
	MsgQueueLen       int64     `json:"msg_queue_len"`
	MaxAwaitRel       int       `json:"max_await_rel"`
	AwaitRelLen       int64     `json:"await_rel_len"`
	MsgDroppedTotal   int64     `json:"msg_dropped_total"`
	MsgDeliveredTotal int64     `json:"msg_delivered_total"`
	ConnectedAt       time.Time `json:"connected_at"`
	DisconnectedAt    time.Time `json:"disconnected_at"`
}

// SubscriptionInfo represents the subscription information
type SubscriptionInfo struct {
	ClientID string    `json:"client_id"`
	Qos      uint8     `json:"qos"`
	Name     string    `json:"name"`
	At       time.Time `json:"at"`
}

var ErrNotFound = errors.New("not found")

type quickList struct {
	index map[string]*list.Element
	rows  *list.List
}
