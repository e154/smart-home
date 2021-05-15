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

// EventStateChanged

package triggers

import (
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/message_queue"
	"sync"
)

const (
	StateChangeName         = "state_change"
	StateChangeFunctionName = "automationTriggerStateChanged"
	StateChangeQueueSize    = 10
)

var _ ITrigger = (*StateChangeTrigger)(nil)

type StateChangeTrigger struct {
	baseTrigger
}

func NewStateChangedTrigger(eventBus event_bus.EventBus) ITrigger {
	return &StateChangeTrigger{
		baseTrigger{
			eventBus:     eventBus,
			msgQueue:     message_queue.New(StateChangeQueueSize),
			functionName: StateChangeFunctionName,
			name:         StateChangeName,
		},
	}
}

func (t *StateChangeTrigger) AsyncAttach(wg *sync.WaitGroup) {

	if err := t.eventBus.Subscribe(event_bus.TopicEntities, t.eventHandler); err != nil {
		log.Error(err.Error())
	}

	wg.Done()
}

func (t *StateChangeTrigger) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case event_bus.EventStateChanged:
		t.msgQueue.Publish(string(v.EntityId), v)
	}
}

func (t *StateChangeTrigger) Subscribe(options Subscriber) error {
	log.Infof("subscribe topic %s", options.EntityId)
	return t.msgQueue.Subscribe(options.EntityId.String(), options.Handler)
}

func (t *StateChangeTrigger) Unsubscribe(options Subscriber) error {
	log.Infof("unsubscribe topic %s", options.EntityId)
	return t.msgQueue.Unsubscribe(options.EntityId.String(), options.Handler)
}
