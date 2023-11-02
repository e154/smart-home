// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package bus

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common/apperr"
)

type bus struct {
	sync.RWMutex
	sub map[string]*subscribers
}

// NewBus ...
func NewBus() Bus {
	return &bus{
		sub: make(map[string]*subscribers),
	}
}

// Publish ...
func (b *bus) Publish(topic string, args ...interface{}) {
	go b.publish(topic, args...)
}

// publish ...
func (b *bus) publish(topic string, args ...interface{}) {
	rArgs := buildHandlerArgs(append([]interface{}{topic}, args...))

	b.RLock()
	defer b.RUnlock()

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

// Subscribe ...
func (b *bus) Subscribe(topic string, fn interface{}, options ...interface{}) error {
	go b.subscribe(topic, fn, options...)
	return nil
}

// subscribe ...
func (b *bus) subscribe(topic string, fn interface{}, options ...interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return errors.Wrap(apperr.ErrInternal, fmt.Sprintf("%s is not a reflect.Func", reflect.TypeOf(fn)))
	}

	const queueSize = 1024

	h := &handler{
		callback: reflect.ValueOf(fn),
		queue:    make(chan []reflect.Value, queueSize),
	}

	b.Lock()
	defer b.Unlock()

	go func() {
		for args := range h.queue {
			go func() {
				startTime := time.Now()
				h.callback.Call(args)
				t := time.Now().Sub(startTime).Microseconds()
				if t > 5000 {
					fmt.Printf("long call! topic %s, fn: %s, Microseconds: %d\n\r", topic, reflect.ValueOf(fn).String(), t)
				}
			}()
		}
	}()

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
		go h.callback.Call(b.sub[topic].lastMsg)
	}

	return nil
}

// Unsubscribe ...
func (b *bus) Unsubscribe(topic string, fn interface{}) error {
	go b.unsubscribe(topic, fn)
	return nil
}

// unsubscribe ...
func (b *bus) unsubscribe(topic string, fn interface{}) error {
	b.Lock()
	defer b.Unlock()

	rv := reflect.ValueOf(fn)

	if _, ok := b.sub[topic]; ok {
		for i, h := range b.sub[topic].handlers {
			if h.callback == rv || h.callback.Pointer() == rv.Pointer() {
				close(h.queue)
				b.sub[topic].handlers = append(b.sub[topic].handlers[:i], b.sub[topic].handlers[i+1:]...)
				if len(b.sub[topic].handlers) == 0 {
					delete(b.sub, topic)
				}
				return nil
			}
		}

		return nil
	}

	return errors.Wrap(apperr.ErrInternal, fmt.Sprintf("topic %s doesn't exist", topic))
}

// CloseTopic ...
func (b *bus) CloseTopic(topic string) {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.sub[topic]; ok {
		for _, h := range b.sub[topic].handlers {
			close(h.queue)
		}
		delete(b.sub, topic)
		return
	}
}

// Purge ...
func (b *bus) Purge() {
	b.Lock()
	defer b.Unlock()

	for topic, s := range b.sub {
		for _, h := range s.handlers {
			close(h.queue)
		}
		delete(b.sub, topic)
	}
}

func (b *bus) Stat(ctx context.Context, limit, offset int64, _, _ string) (result Stats, total int64, err error) {
	b.RLock()
	defer b.RUnlock()

	var stats Stats
	for topic, subs := range b.sub {
		stats = append(stats, Stat{
			Topic:       topic,
			Subscribers: len(subs.handlers),
		})
	}

	sort.Sort(stats)

	total = int64(len(stats))

	if offset > total {
		offset = total
	}

	end := offset + limit
	if end > total {
		end = total
	}

	result = stats[offset:end]

	return
}

func buildHandlerArgs(args []interface{}) []reflect.Value {
	reflectedArgs := make([]reflect.Value, 0)

	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	return reflectedArgs
}
