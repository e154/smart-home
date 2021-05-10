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

func (a *StateChangeTrigger) AsyncAttach(wg *sync.WaitGroup) {

	if err := a.eventBus.Subscribe(event_bus.TopicEntities, a.eventHandler); err != nil {
		log.Error(err.Error())
	}

	wg.Done()
}

func (a *StateChangeTrigger) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case event_bus.EventStateChanged:
		a.msgQueue.Publish(string(v.EntityId), v)
	}
}

func (b *StateChangeTrigger) Subscribe(topic string, fn interface{}, _ interface{}) error {
	log.Infof("subscribe topic %s", topic)
	return b.msgQueue.Subscribe(topic, fn)
}

func (b *StateChangeTrigger) Unsubscribe(topic string, fn interface{}, _ interface{}) error {
	log.Infof("unsubscribe topic %s", topic)
	return b.msgQueue.Unsubscribe(topic, fn)
}
