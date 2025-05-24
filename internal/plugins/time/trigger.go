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

package time

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/e154/smart-home/pkg/plugins/triggers"
	"github.com/e154/smart-home/pkg/scheduler"

	"github.com/e154/bus"
)

var _ triggers.ITrigger = (*Trigger)(nil)

type subscribe struct {
	callback reflect.Value
	entryID  scheduler.EntryID
}

// Trigger ...
type Trigger struct {
	eventBus     bus.Bus
	functionName string
	name         string
	scheduler    scheduler.Scheduler
	sync.Mutex
	subscribers map[string][]*subscribe
}

// NewTrigger ...
func NewTrigger(eventBus bus.Bus,
	scheduler scheduler.Scheduler) triggers.ITrigger {

	return &Trigger{
		eventBus:     eventBus,
		scheduler:    scheduler,
		subscribers:  make(map[string][]*subscribe),
		functionName: FunctionName,
		name:         Name,
	}
}

// Name ...
func (t *Trigger) Name() string {
	return t.name
}

// AsyncAttach ...
func (t *Trigger) AsyncAttach(wg *sync.WaitGroup) {

	wg.Done()
}

// Subscribe ...
func (t *Trigger) Subscribe(options triggers.Subscriber) error {
	if options.Payload == nil {
		return fmt.Errorf("payload is nil")
	}
	if _, ok := options.Payload[AttrCron]; !ok {
		return fmt.Errorf("cron attribute is nil")
	}
	schedule := options.Payload[AttrCron].String()
	if schedule == "" {
		return fmt.Errorf("error static cast to string %v", options.Payload)
	}
	callback := reflect.ValueOf(options.Handler)
	entryID, err := t.scheduler.AddFunc(schedule, func() {
		callback.Call([]reflect.Value{reflect.ValueOf(""), reflect.ValueOf(time.Now())})
	})

	if err != nil {
		return err
	}

	sub := &subscribe{
		callback: callback,
		entryID:  entryID,
	}
	t.Lock()
	t.subscribers[schedule] = append(t.subscribers[schedule], sub)
	t.Unlock()

	return nil
}

// Unsubscribe ...
func (t *Trigger) Unsubscribe(options triggers.Subscriber) error {
	if options.Payload == nil {
		return fmt.Errorf("payload is nil")
	}
	schedule := options.Payload[AttrCron].String()
	if schedule == "" {
		return fmt.Errorf("error static cast to string %v", options.Payload)
	}
	rv := reflect.ValueOf(options.Handler)

	t.Lock()
	defer t.Unlock()

	var indexesToDelete []int

	for i, sub := range t.subscribers[schedule] {
		if sub.callback == rv || sub.callback.Pointer() == rv.Pointer() {
			indexesToDelete = append(indexesToDelete, i)
		}
	}

	for i := len(indexesToDelete) - 1; i >= 0; i-- {
		index := indexesToDelete[i]
		t.scheduler.Remove(t.subscribers[schedule][index].entryID)
		t.subscribers[schedule] = append(t.subscribers[schedule][:index], t.subscribers[schedule][index+1:]...)
	}

	return nil
}

// FunctionName ...
func (t *Trigger) FunctionName() string {
	return t.functionName
}
