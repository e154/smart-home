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
	"github.com/e154/smart-home/common/events"
	"go.uber.org/atomic"
	"sync"

	"github.com/e154/smart-home/system/bus"
)

const (
	// TopicSystem ...
	TopicSystem = "system/#"
	// SystemName ...
	SystemName = "system"
	// SystemFunctionName ...
	SystemFunctionName = "automationTriggerSystem"
)

var _ ITrigger = (*SystemTrigger)(nil)

// SystemTrigger ...
type SystemTrigger struct {
	baseTrigger
	counter *atomic.Int32
}

// NewSystemTrigger ...
func NewSystemTrigger(eventBus bus.Bus) ITrigger {
	return &SystemTrigger{
		baseTrigger: baseTrigger{
			eventBus:     eventBus,
			msgQueue:     bus.NewBus(),
			functionName: SystemFunctionName,
			name:         SystemName,
		},
		counter: atomic.NewInt32(0),
	}
}

// AsyncAttach ...
func (t *SystemTrigger) AsyncAttach(wg *sync.WaitGroup) {

	if err := t.eventBus.Subscribe(TopicSystem, t.eventHandler); err != nil {
		log.Error(err.Error())
	}

	wg.Done()
}

func (t *SystemTrigger) eventHandler(topic string, event interface{}) {
	if t.counter.Load() <= 0 {
		return
	}
	switch event.(type) {
	case events.EventStateChanged:
		return
	}

	t.msgQueue.Publish(topic, &SystemTriggerMessage{
		Topic:     topic,
		EventName: events.EventName(event),
		Event:     event,
	})
}

// Subscribe ...
func (t *SystemTrigger) Subscribe(options Subscriber) error {
	//log.Infof("subscribe topic %s", TopicSystem)
	t.counter.Inc()
	return t.msgQueue.Subscribe(TopicSystem, options.Handler)
}

// Unsubscribe ...
func (t *SystemTrigger) Unsubscribe(options Subscriber) error {
	//log.Infof("unsubscribe topic %s", TopicSystem)
	t.counter.Dec()
	return t.msgQueue.Unsubscribe(TopicSystem, options.Handler)
}
