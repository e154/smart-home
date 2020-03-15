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
	"github.com/e154/smart-home/adaptors"
	"sync"
	"time"
)

type History struct {
	Items []HistoryItem `json:"items"`
}

type HistoryManager struct {
	publisher  IPublisher
	adaptors   *adaptors.Adaptors
	updateLock sync.Mutex
	pool       []HistoryItem
}

func NewHistoryManager(publisher IPublisher, adaptors *adaptors.Adaptors) *HistoryManager {
	manager := &HistoryManager{
		publisher: publisher,
		adaptors:  adaptors,
		pool:      make([]HistoryItem, 0),
	}
	manager.init()
	return manager
}

func (d *HistoryManager) update(t interface{}) {
	switch v := t.(type) {
	case HistoryItem:
		d.updatePool(v)
	default:
		return
	}

	d.broadcast()
}

func (d *HistoryManager) updatePool(item HistoryItem) {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	const (
		maxItems = 20
	)

	lenPool := len(d.pool)
	i := 1

	if lenPool > maxItems {
		i = lenPool - maxItems
	}

	d.pool = d.pool[i:lenPool]
	d.pool = append(d.pool, item)
}

func (d HistoryManager) Snapshot() History {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	items := make([]HistoryItem, len(d.pool))
	copy(items, d.pool)

	return History{
		Items: items,
	}
}

func (d *HistoryManager) broadcast() {
	go d.publisher.Broadcast("history")
}

func (d *HistoryManager) init() {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	if list, err := d.adaptors.MapDeviceHistory.List(20, 0); err == nil {
		for _, item := range list {
			d.pool = append(d.pool, HistoryItem{
				DeviceName:        item.MapElement.Name,
				DeviceDescription: item.MapElement.Description,
				Type:              string(item.Type),
				Description:       item.Description,
				CreatedAt:         item.CreatedAt,
			})
		}
	}

}

type HistoryItem struct {
	DeviceName        string    `json:"device_name"`
	DeviceDescription string    `json:"device_description"`
	Type              string    `json:"type"`
	Description       string    `json:"description"`
	CreatedAt         time.Time `json:"created_at"`
}
