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

package metrics

import (
	"go.uber.org/atomic"
	"sync"
)

type Mqtt struct {
	ClientState MqttClientStats `json:"client_state"`
}

type MqttManager struct {
	publisher   IPublisher
	isStarted   atomic.Bool
	quit        chan struct{}
	updateLock  sync.Mutex
	clientState MqttClientStats
}

func NewMqttManager(publisher IPublisher) *MqttManager {
	return &MqttManager{publisher: publisher}
}

func (d *MqttManager) update(t interface{}) {
	switch v := t.(type) {
	case MqttClientStats:
		d.updateClientState(v)
	default:
		return
	}

	d.broadcast()
}

func (d *MqttManager) updateClientState(clientState MqttClientStats) {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	d.clientState = clientState
}

func (d *MqttManager) Snapshot() Mqtt {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	return Mqtt{
		ClientState: d.clientState,
	}
}

func (d *MqttManager) broadcast() {
	go d.publisher.Broadcast("mqtt")
}

type MqttClientStats struct {
	ConnectedTotal    uint64 `json:"connected_total"`
	DisconnectedTotal uint64 `json:"disconnected_total"`
	ActiveCurrent     uint64 `json:"active_current"`
	InactiveCurrent   uint64 `json:"inactive_current"`
	ExpiredTotal      uint64 `json:"expired_total"`
}
