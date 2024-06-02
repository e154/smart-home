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

package triggers

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/e154/bus"
	"github.com/e154/smart-home/system/scheduler"
)

const (
	// TimeName ...
	TimeName = "time"
	// TimeFunctionName ...
	TimeFunctionName = "automationTriggerTime"
	// TimeQueueSize ...
	TimeQueueSize = 10
)

var _ ITrigger = (*TimeTrigger)(nil)

type subscribe struct {
	callback reflect.Value
	entryID  scheduler.EntryID
}

// TimeTrigger ...
type TimeTrigger struct {
	baseTrigger
	scheduler *scheduler.Scheduler
	sync.Mutex
	subscribers map[string][]*subscribe
}

// NewTimeTrigger ...
func NewTimeTrigger(eventBus bus.Bus,
	scheduler *scheduler.Scheduler) ITrigger {

	return &TimeTrigger{
		scheduler:   scheduler,
		subscribers: make(map[string][]*subscribe),
		baseTrigger: baseTrigger{
			eventBus:     eventBus,
			msgQueue:     bus.NewBus(),
			functionName: TimeFunctionName,
			name:         TimeName,
		},
	}
}

// AsyncAttach ...
func (t *TimeTrigger) AsyncAttach(wg *sync.WaitGroup) {

	wg.Done()
}

// Subscribe ...
func (t *TimeTrigger) Subscribe(options Subscriber) error {
	if options.Payload == nil {
		return fmt.Errorf("payload is nil")
	}
	if _, ok := options.Payload[CronOptionTrigger]; !ok {
		return fmt.Errorf("cron attribute is nil")
	}
	schedule := options.Payload[CronOptionTrigger].String()
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
func (t *TimeTrigger) Unsubscribe(options Subscriber) error {
	if options.Payload == nil {
		return fmt.Errorf("payload is nil")
	}
	schedule := options.Payload[CronOptionTrigger].String()
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
