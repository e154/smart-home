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
	"github.com/gammazero/deque"
	"sync"
	"time"
)

// History ...
type History struct {
	Items []HistoryItem `json:"items"`
}

// HistoryManager ...
type HistoryManager struct {
	publisher  IPublisher
	adaptors   *adaptors.Adaptors
	updateLock *sync.Mutex
	queue      *deque.Deque
}

// NewHistoryManager ...
func NewHistoryManager(publisher IPublisher, adaptors *adaptors.Adaptors) *HistoryManager {
	manager := &HistoryManager{
		publisher:  publisher,
		adaptors:   adaptors,
		queue:      &deque.Deque{},
		updateLock: &sync.Mutex{},
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
		maxItems = 8
	)

	for d.queue.Len() > maxItems {
		d.queue.PopBack()
	}

	d.queue.PushFront(HistoryItem{
		DeviceName:        item.DeviceName,
		DeviceDescription: item.DeviceDescription,
		Type:              item.Type,
		LogLevel:          item.LogLevel,
		Description:       item.Description,
		CreatedAt:         item.CreatedAt,
	})
}

// Snapshot ...
func (d HistoryManager) Snapshot() History {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	items := make([]HistoryItem, d.queue.Len())
	for i := 0; i < d.queue.Len(); i++ {
		if item, ok := d.queue.At(i).(HistoryItem); ok {
			items[i] = item
		}
	}

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

	//if list, err := d.adaptors.EntityHistory.List(8, 0); err == nil {
	//	for _, item := range list {
	//		d.queue.PushBack(HistoryItem{
	//			DeviceName:        item.Entity.Name,
	//			DeviceDescription: item.Entity.Description,
	//			Type:              string(item.Type),
	//			LogLevel:          string(item.LogLevel),
	//			Description:       item.Description,
	//			CreatedAt:         item.CreatedAt,
	//		})
	//	}
	//}

}

// HistoryItem ...
type HistoryItem struct {
	DeviceName        string    `json:"device_name"`
	DeviceDescription string    `json:"device_description"`
	Type              string    `json:"type"`
	LogLevel          string    `json:"log_level"`
	Description       string    `json:"description"`
	CreatedAt         time.Time `json:"created_at"`
}
