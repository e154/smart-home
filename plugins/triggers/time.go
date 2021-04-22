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

package triggers

import (
	"fmt"
	"github.com/e154/smart-home/system/cron"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/message_queue"
	"reflect"
	"sync"
	"time"
)

const (
	TimeName         = "time"
	TimeFunctionName = "automationTriggerTime"
	TimeQueueSize    = 10
)

type subscribe struct {
	callback reflect.Value
	task     *cron.Task
}

type TimeTrigger struct {
	baseTrigger
	cron *cron.Cron
	sync.Mutex
	subscribers map[string][]*subscribe
}

func NewTimeTrigger(eventBus event_bus.EventBus) ITrigger {
	c := cron.NewCron()
	go c.Run()
	return &TimeTrigger{
		cron:        c,
		subscribers: make(map[string][]*subscribe),
		baseTrigger: baseTrigger{
			eventBus:     eventBus,
			msgQueue:     message_queue.New(TimeQueueSize),
			functionName: TimeFunctionName,
			name:         TimeName,
		},
	}
}

func (t *TimeTrigger) AsyncAttach(wg *sync.WaitGroup) {

	wg.Done()
}

func (t *TimeTrigger) Subscribe(_ string, fn interface{}, payload interface{}) error {
	schedule, ok := payload.(string)
	if !ok {
		return fmt.Errorf("error static cast to string %v", payload)
	}
	callback := reflect.ValueOf(fn)
	task, err := t.cron.NewTask(schedule, func() {
		callback.Call([]reflect.Value{reflect.ValueOf(time.Now())})
	})

	if err != nil {
		return err
	}

	sub := &subscribe{
		callback: callback,
		task:     task,
	}
	t.Lock()
	t.subscribers[schedule] = append(t.subscribers[schedule], sub)
	t.Unlock()

	return nil
}

func (t *TimeTrigger) Unsubscribe(_ string, fn interface{}, payload interface{}) error {
	schedule, ok := payload.(string)
	if !ok {
		return fmt.Errorf("error static cast to string %v", payload)
	}
	callback := reflect.ValueOf(fn)

	t.Lock()
	defer t.Unlock()

	if len(t.subscribers[schedule]) == 1 {
		t.subscribers[schedule] = []*subscribe{}
		return nil
	}

	for i, sub := range t.subscribers[schedule] {
		if sub.callback == callback {
			t.cron.RemoveTask(sub.task)
			t.subscribers[schedule] = append(t.subscribers[schedule][:i], t.subscribers[schedule][i+1:]...)
		}
	}

	return nil
}
