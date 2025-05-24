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

package admin

// SubscriptionInfo represents the subscription information
type SubscriptionInfo struct {
	Id                uint32 `json:"id"`
	ClientID          string `json:"client_id"`
	TopicName         string `json:"topic_name"`
	Name              string `json:"name"`
	Qos               uint32 `json:"qos"`
	NoLocal           bool   `json:"no_local"`
	RetainAsPublished bool   `json:"retain_as_published"`
	RetainHandling    uint32 `json:"retain_handling"`
}
