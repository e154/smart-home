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
	"os"
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
			sub.rps.Inc()
			h.queue <- rArgs
		}
	}
}

// Subscribe ...
func (b *bus) Subscribe(topic string, fn interface{}, options ...interface{}) (err error) {
	if err = b.subscribe(topic, fn, options...); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	return err
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
	if _, ok := b.sub[topic]; ok {
		b.sub[topic].handlers = append(b.sub[topic].handlers, h)
	} else {
		b.sub[topic] = newSubscibers(h)
	}
	b.Unlock()

	go func(subs *subscribers) {
		var startTime time.Time
		for args := range h.queue {
			go func(args []reflect.Value) {
				startTime = time.Now()
				h.callback.Call(args)
				subs.setTime(time.Since(startTime))
			}(args)
		}
	}(b.sub[topic])

	if len(options) > 0 {
		if retain, ok := options[0].(bool); ok && !retain {
			return nil
		}
	}

	b.RLock()
	// sand last message value
	if b.sub[topic].lastMsg != nil {
		go h.callback.Call(b.sub[topic].lastMsg)
	}
	b.RUnlock()

	return nil
}

// Unsubscribe ...
func (b *bus) Unsubscribe(topic string, fn interface{}) (err error) {
	if err = b.unsubscribe(topic, fn); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	return
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
					b.sub[topic].rps.Stop()
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

	if s, ok := b.sub[topic]; ok {
		s.stop()
		delete(b.sub, topic)
		return
	}
}

// Purge ...
func (b *bus) Purge() {
	b.Lock()
	defer b.Unlock()

	for topic, s := range b.sub {
		s.stop()
		delete(b.sub, topic)
	}
}

func (b *bus) Stat(ctx context.Context, limit, offset int64, _, _ string) (result Stats, total int64, err error) {

	var stats Stats

	b.RLock()
	for topic, subs := range b.sub {
		stats = append(stats, StatItem{
			Topic:       topic,
			Subscribers: len(subs.handlers),
			Min:         subs.min.Load(),
			Max:         subs.max.Load(),
			Avg:         subs.avg.Load(),
			Rps:         subs.rps.Value(),
		})
	}
	b.RUnlock()

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
