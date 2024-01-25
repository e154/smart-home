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

package webpush

import m "github.com/e154/smart-home/models"

// EventAddWebPushSubscription ...
type EventAddWebPushSubscription struct {
	UserID       int64           `json:"user_id"`
	SessionID    string          `json:"session_id"`
	Subscription *m.Subscription `json:"subscription"`
}

// EventGetWebPushPublicKey ...
type EventGetWebPushPublicKey struct {
	UserID    int64  `json:"user_id,omitempty"`
	SessionID string `json:"session_id"`
}

// EventNewWebPushPublicKey ...
type EventNewWebPushPublicKey struct {
	UserID    int64  `json:"user_id,omitempty"`
	SessionID string `json:"session_id"`
	PublicKey string `json:"public_key"`
}
