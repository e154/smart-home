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

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common/apperr"
)

type bus struct {
	sync.RWMutex
	topics map[string]*Topic
}

// NewBus ...
func NewBus() Bus {
	return &bus{
		topics: make(map[string]*Topic),
	}
}

// Publish ...
func (b *bus) Publish(topic string, args ...interface{}) {
	b.RLock()
	defer b.RUnlock()

	for t, sub := range b.topics {
		if !TopicMatch([]byte(topic), []byte(t)) {
			continue
		}
		go sub.Publish(args...)
	}
}

// Subscribe ...
func (b *bus) Subscribe(topic string, fn interface{}, options ...interface{}) (err error) {
	b.Lock()
	defer b.Unlock()

	if sub, ok := b.topics[topic]; ok {
		return sub.Subscribe(fn, options...)
	} else {
		subs := NewTopic(topic)
		b.topics[topic] = subs
		return subs.Subscribe(fn, options...)
	}
}

// Unsubscribe ...
func (b *bus) Unsubscribe(topic string, fn interface{}) (err error) {
	b.Lock()
	defer b.Unlock()

	var empty bool
	if sub, ok := b.topics[topic]; ok {
		empty, err = sub.Unsubscribe(fn)
		if err != nil {
			return err
		}
		if empty {
			delete(b.topics, topic)
		}
		return nil
	}

	return errors.Wrap(apperr.ErrInternal, fmt.Sprintf("topic %s doesn't exist", topic))
}

// CloseTopic ...
func (b *bus) CloseTopic(topic string) {
	b.Lock()
	defer b.Unlock()

	if sub, ok := b.topics[topic]; ok {
		sub.Close()
		delete(b.topics, topic)
	}
}

// Purge ...
func (b *bus) Purge() {
	b.Lock()
	defer b.Unlock()

	for topic, sub := range b.topics {
		sub.Close()
		delete(b.topics, topic)
	}
}

func (b *bus) Stat(ctx context.Context, limit, offset int64, _, _ string) (result Stats, total int64, err error) {

	b.RLock()
	var stats = make(Stats, 0, len(b.topics))
	for _, subs := range b.topics {
		stats = append(stats, subs.Stat())
	}
	b.RUnlock()

	sort.Sort(stats)

	total = int64(len(b.topics))

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
