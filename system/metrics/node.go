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
	"github.com/rcrowley/go-metrics"
	"sync"
)

type Node struct {
	Total  int64            `json:"total"`
	Status map[int64]string `json:"status"`
}

type NodeManager struct {
	publisher  IPublisher
	total      metrics.Counter
	updateLock sync.Mutex
	status     map[int64]string
}

func NewNodeManager(publisher IPublisher) *NodeManager {
	return &NodeManager{
		publisher: publisher,
		status:    make(map[int64]string),
		total:     metrics.NewCounter(),
	}
}

func (d *NodeManager) update(t interface{}) {
	switch v := t.(type) {
	case NodeAdd:
		d.total.Inc(v.Num)
	case NodeDelete:
		d.total.Dec(v.Num)
	case NodeUpdateStatus:
		d.updateLock.Lock()
		defer d.updateLock.Unlock()

		if d.status[v.Id] == v.Status {
			return
		}
		d.status[v.Id] = v.Status

	default:
		return
	}

	d.broadcast()
}

func (d *NodeManager) GetStatus(nodeId int64) (status string, err error) {

	var ok bool
	if status, ok = d.status[nodeId]; ok {
		return
	}

	err = ErrRecordNotFound

	return
}

func (d *NodeManager) Snapshot() Node {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	status := make(map[int64]string)
	for k, v := range d.status {
		status[k] = v
	}
	return Node{
		Total:  d.total.Count(),
		Status: status,
	}
}

func (d *NodeManager) broadcast() {
	go d.publisher.Broadcast("node")
}

type NodeUpdateStatus struct {
	Id     int64
	Status string
}

type NodeAdd struct {
	Num int64
}

type NodeDelete struct {
	Num int64
}
