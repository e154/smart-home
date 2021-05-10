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

package message_queue

import (
	"reflect"
	"sync"
)

// MessageQueue implements publish/subscribe messaging paradigm
type MessageQueue interface {
	// Publish publishes arguments to the given topic subscribers
	// Publish block only when the buffer of one of the subscribers is full.
	Publish(topic string, args ...interface{})
	// Close unsubscribe all subscribers from given topic
	Close(topic string)
	// Subscribe subscribes to the given topic
	Subscribe(topic string, fn interface{}, options ...interface{}) error
	// Unsubscribe unsubscribe handler from the given topic
	Unsubscribe(topic string, fn interface{}) error
}

type handler struct {
	callback reflect.Value
	queue    chan []reflect.Value
}

type subscribers struct {
	handlers []*handler
	lastMsg  []reflect.Value
}

type messageQueue struct {
	queueSize int
	mtx       sync.RWMutex
	sub       map[string]*subscribers
}
