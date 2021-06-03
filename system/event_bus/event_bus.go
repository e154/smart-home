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

package event_bus

import (
	"github.com/e154/smart-home/system/message_queue"
)

const (
	queueSize = 100
)

// EventBus ...
type EventBus interface {
	Publish(topic string, args ...interface{})
	Close(topic string)
	Subscribe(topic string, fn interface{}, options ...interface{}) error
	Unsubscribe(topic string, fn interface{}) error
	Stat() (stats message_queue.Stats, err error)
	Purge()
}

type eventBus struct {
	queue message_queue.MessageQueue
}

// NewEventBus ...
func NewEventBus() EventBus {
	return &eventBus{
		queue: message_queue.New(queueSize),
	}
}

// Publish ...
func (e *eventBus) Publish(topic string, args ...interface{}) {
	e.queue.Publish(topic, args...)
}

// Close ...
func (e *eventBus) Close(topic string) {
	e.queue.Close(topic)
}

// Subscribe ...
func (e *eventBus) Subscribe(topic string, fn interface{}, options ...interface{}) error {
	return e.queue.Subscribe(topic, fn, options...)
}

// Unsubscribe ...
func (e *eventBus) Unsubscribe(topic string, fn interface{}) error {
	return e.queue.Unsubscribe(topic, fn)
}

// Stat ...
func (e *eventBus) Stat() (message_queue.Stats, error) {
	return e.queue.Stat()
}

// Purge ...
func (e *eventBus) Purge() {
	e.queue.Purge()
}
