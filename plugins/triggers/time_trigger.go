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

package triggers

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/e154/smart-home/system/bus"
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
	schedule := options.Payload[CronOptionTrigger].String()
	if schedule == "" {
		return fmt.Errorf("error static cast to string %v", options.Payload)
	}
	callback := reflect.ValueOf(options.Handler)

	t.Lock()
	defer t.Unlock()

	if len(t.subscribers[schedule]) == 1 {
		t.subscribers[schedule] = []*subscribe{}
		return nil
	}

	for i, sub := range t.subscribers[schedule] {
		if sub.callback == callback {
			t.scheduler.Remove(sub.entryID)
			t.subscribers[schedule] = append(t.subscribers[schedule][:i], t.subscribers[schedule][i+1:]...)
		}
	}

	return nil
}

// CallManual ...
func (t *TimeTrigger) CallManual() {
	log.Warn("method not implemented")
}
