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

package metrics

import (
	"github.com/e154/smart-home/common"
	"sync"

	"github.com/rcrowley/go-metrics"
)

// Node ...
type Node struct {
	Total  int64            `json:"total"`
	Status map[int64]string `json:"status"`
}

// NodeManager ...
type NodeManager struct {
	publisher  IPublisher
	total      metrics.Counter
	updateLock sync.Mutex
	status     map[int64]string
}

// NewNodeManager ...
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
		if d.updateStatus(v) {
			return
		}

	default:
		return
	}

	d.broadcast()
}

func (d *NodeManager) updateStatus(v NodeUpdateStatus) (exist bool) {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	if d.status[v.Id] == v.Status {
		exist = true
		return
	}
	d.status[v.Id] = v.Status

	return
}

// GetStatus ...
func (d *NodeManager) GetStatus(nodeId int64) (status string, err error) {

	var ok bool
	if status, ok = d.status[nodeId]; ok {
		return
	}

	err = common.ErrNotFound

	return
}

// Snapshot ...
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

// NodeUpdateStatus ...
type NodeUpdateStatus struct {
	Id     int64
	Status string
}

// NodeAdd ...
type NodeAdd struct {
	Num int64
}

// NodeDelete ...
type NodeDelete struct {
	Num int64
}
