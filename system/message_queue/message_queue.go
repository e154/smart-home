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
	"fmt"
	"reflect"
)

// New creates new MessageQueue
// queueSize sets buffered channel length per subscriber
func New(handlerQueueSize int) MessageQueue {
	if handlerQueueSize == 0 {
		panic("queueSize has to be greater then 0")
	}

	return &messageQueue{
		queueSize: handlerQueueSize,
		sub:       make(map[string]*subscribers),
	}
}

func (b *messageQueue) Publish(topic string, args ...interface{}) {
	qwe := []interface{}{topic}
	rArgs := buildHandlerArgs(append(qwe, args...))

	b.mtx.Lock()
	defer b.mtx.Unlock()

	for t, sub := range b.sub {
		if !TopicMatch([]byte(topic), []byte(t)) {
			continue
		}
		sub.lastMsg = rArgs
		for _, h := range sub.handlers {
			h.queue <- rArgs
		}
	}
}

func (b *messageQueue) Subscribe(topic string, fn interface{}, options ...interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	h := &handler{
		callback: reflect.ValueOf(fn),
		queue:    make(chan []reflect.Value, b.queueSize),
	}

	go func() {
		for args := range h.queue {
			h.callback.Call(args)
		}
	}()

	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, ok := b.sub[topic]; ok {
		b.sub[topic].handlers = append(b.sub[topic].handlers, h)
	} else {
		b.sub[topic] = &subscribers{
			handlers: []*handler{h},
		}
	}

	if len(options) > 0 {
		if retain, ok := options[0].(bool); ok && !retain {
			return nil
		}
	}

	// sand last message value
	if b.sub[topic].lastMsg != nil {
		h.callback.Call(b.sub[topic].lastMsg)
	}

	return nil
}

func (b *messageQueue) Unsubscribe(topic string, fn interface{}) error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	rv := reflect.ValueOf(fn)

	if _, ok := b.sub[topic]; ok {
		for i, h := range b.sub[topic].handlers {
			if h.callback == rv {
				close(h.queue)

				b.sub[topic].handlers = append(b.sub[topic].handlers[:i], b.sub[topic].handlers[i+1:]...)
			}
		}

		return nil
	}

	return fmt.Errorf("topic %s doesn't exist", topic)
}

func (b *messageQueue) Close(topic string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, ok := b.sub[topic]; ok {
		for _, h := range b.sub[topic].handlers {
			close(h.queue)
		}

		delete(b.sub, topic)

		return
	}
}

func buildHandlerArgs(args []interface{}) []reflect.Value {
	reflectedArgs := make([]reflect.Value, 0)

	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	return reflectedArgs
}
