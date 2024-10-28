// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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

package system

import (
	"sync"

	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/plugins/triggers"

	"github.com/e154/bus"
	"go.uber.org/atomic"
)

var _ triggers.ITrigger = (*Trigger)(nil)

// Trigger ...
type Trigger struct {
	eventBus     bus.Bus
	msgQueue     bus.Bus
	functionName string
	name         string
	counter      *atomic.Int32
}

// NewTrigger ...
func NewTrigger(eventBus bus.Bus) triggers.ITrigger {
	return &Trigger{
		eventBus:     eventBus,
		msgQueue:     bus.NewBus(),
		functionName: FunctionName,
		name:         Name,
		counter:      atomic.NewInt32(0),
	}
}

func (t *Trigger) Name() string {
	return t.name
}

// AsyncAttach ...
func (t *Trigger) AsyncAttach(wg *sync.WaitGroup) {

	if err := t.eventBus.Subscribe(TopicSystem, t.eventHandler); err != nil {
		log.Error(err.Error())
	}

	wg.Done()
}

func (t *Trigger) eventHandler(topic string, event interface{}) {
	if t.counter.Load() <= 0 {
		return
	}
	switch event.(type) {
	case events.EventStateChanged:
		return
	}

	t.msgQueue.Publish(topic, &TriggerMessage{
		Topic:     topic,
		EventName: events.EventName(event),
		Event:     event,
	})
}

// Subscribe ...
func (t *Trigger) Subscribe(options triggers.Subscriber) error {
	//log.Infof("subscribe topic %s", TopicSystem)
	t.counter.Inc()
	return t.msgQueue.Subscribe(TopicSystem, options.Handler)
}

// Unsubscribe ...
func (t *Trigger) Unsubscribe(options triggers.Subscriber) error {
	//log.Infof("unsubscribe topic %s", TopicSystem)
	t.counter.Dec()
	return t.msgQueue.Unsubscribe(TopicSystem, options.Handler)
}

// FunctionName ...
func (t *Trigger) FunctionName() string {
	return t.functionName
}
