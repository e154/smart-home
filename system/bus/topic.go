// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common/apperr"
)

const queueSize = 1024

type Topic struct {
	name string
	*Statistic
	sync.RWMutex
	handlers []*handler
	lastMsg  []reflect.Value
}

func NewTopic(name string) *Topic {
	return &Topic{
		name:      name,
		handlers:  make([]*handler, 0),
		Statistic: NewStatistic(),
	}
}

func (t *Topic) Publish(args ...interface{}) {
	t.RLock()
	defer t.RUnlock()

	if len(t.handlers) == 0 {
		return
	}

	rArgs := buildHandlerArgs(append([]interface{}{t.name}, args...))

	t.lastMsg = rArgs
	for _, h := range t.handlers {
		t.rps.Inc()
		h.queue <- rArgs
	}
}

func (t *Topic) Subscribe(fn interface{}, options ...interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return errors.Wrap(apperr.ErrInternal, fmt.Sprintf("%s is not a reflect.Func", reflect.TypeOf(fn)))
	}

	h := &handler{
		callback: reflect.ValueOf(fn),
		queue:    make(chan []reflect.Value, queueSize),
	}

	t.Lock()
	t.handlers = append(t.handlers, h)
	t.Unlock()

	go func(h *handler) {
		var startTime time.Time
		for args := range h.queue {
			go func(args []reflect.Value) {
				startTime = time.Now()
				h.callback.Call(args)
				t.setTime(time.Since(startTime))
			}(args)
		}
	}(h)

	if len(options) > 0 {
		if retain, ok := options[0].(bool); ok && !retain {
			return nil
		}
	}

	t.RLock()
	// sand last message value
	if t.lastMsg != nil {
		go h.callback.Call(t.lastMsg)
	}
	t.RUnlock()

	return nil
}

func (t *Topic) Unsubscribe(fn interface{}) (empty bool, err error) {
	t.Lock()
	defer t.Unlock()

	rv := reflect.ValueOf(fn)

	for i, h := range t.handlers {
		if h.callback == rv || h.callback.Pointer() == rv.Pointer() {
			close(h.queue)
			t.handlers = append(t.handlers[:i], t.handlers[i+1:]...)
		}
	}

	empty = len(t.handlers) == 0

	return
}

func (t *Topic) Close() {
	t.Lock()
	defer t.Unlock()

	for _, h := range t.handlers {
		close(h.queue)
	}
	t.handlers = make([]*handler, 0)
	t.rps.Stop()
}

func (t *Topic) Stat() *StatItem {
	t.RLock()
	defer t.RUnlock()
	return &StatItem{
		Topic:       t.name,
		Subscribers: len(t.handlers),
		Min:         t.min.Load(),
		Max:         t.max.Load(),
		Avg:         t.avg.Load(),
		Rps:         t.rps.Value(),
	}
}
