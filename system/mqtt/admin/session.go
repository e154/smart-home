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

import "time"

// SessionInfo represents the session information
type SessionInfo struct {
	ClientID              string     `json:"client_id"`
	Status                string     `json:"status"`
	CleanSession          bool       `json:"clean_session"`
	Subscriptions         uint64     `json:"subscriptions"`
	MaxInflight           uint16     `json:"max_inflight"`
	InflightLen           uint64     `json:"inflight_len"`
	MaxMsgQueue           int        `json:"max_msg_queue"`
	MsgQueueLen           uint64     `json:"msg_queue_len"`
	AwaitRelLen           uint64     `json:"await_rel_len"`
	Qos0MsgDroppedTotal   uint64     `json:"qos0_msg_dropped_total"`
	Qos1MsgDroppedTotal   uint64     `json:"qos1_msg_dropped_total"`
	Qos2MsgDroppedTotal   uint64     `json:"qos2_msg_dropped_total"`
	Qos0MsgDeliveredTotal uint64     `json:"qos0_msg_delivered_total"`
	Qos1MsgDeliveredTotal uint64     `json:"qos1_msg_delivered_total"`
	Qos2MsgDeliveredTotal uint64     `json:"qos2_msg_delivered_total"`
	ConnectedAt           *time.Time `json:"connected_at"`
	DisconnectedAt        *time.Time `json:"disconnected_at"`
}
