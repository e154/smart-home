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

package event_bus

import (
	"github.com/e154/smart-home/system/message_queue"
)

const (
	queueSize = 100
)

type EventBus struct {
	bus message_queue.MessageQueue
}

func NewEventBus() *EventBus {
	return &EventBus{
		bus: message_queue.New(queueSize),
	}
}

func (e *EventBus) Publish(topic string, args ...interface{}) {
	e.bus.Publish(topic, args...)
}

func (e *EventBus) Close(topic string) {
	e.bus.Close(topic)
}

func (e *EventBus) Subscribe(topic string, fn interface{}) error {
	return e.bus.Subscribe(topic, fn)
}

func (e *EventBus) Unsubscribe(topic string, fn interface{}) error {
	return e.bus.Unsubscribe(topic, fn)
}
