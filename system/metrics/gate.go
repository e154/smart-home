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
	"sync"
)

type Gate struct {
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
}

type GateManager struct {
	updateLock  sync.Mutex
	status      string
	accessToken string
	publisher   IPublisher
}

func NewGateManager(publisher IPublisher) *GateManager {
	return &GateManager{publisher: publisher}
}

func (d *GateManager) update(t interface{}) {
	switch v := t.(type) {
	case GateUpdate:
		if !d.selfUpdate(v) {
			return
		}
	default:
		return
	}

	d.broadcast()
}

func (d *GateManager) selfUpdate(v GateUpdate) (broadcast bool) {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	broadcast = d.status != v.Status || d.accessToken != v.AccessToken

	if v.Status != "" {
		d.status = v.Status
	}
	d.accessToken = v.AccessToken

	return
}

func (d *GateManager) Snapshot() Gate {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	return Gate{
		Status:      d.status,
		AccessToken: d.accessToken,
	}
}

func (d *GateManager) broadcast() {
	go d.publisher.Broadcast("gate")
}

type GateUpdate struct {
	Status, AccessToken string
}
