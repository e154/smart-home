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
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/message_queue"
	"sync"
)

const (
	TopicSystem        = "system"
	EventStart         = "START"
	EventStop          = "STOP"
	SystemName         = "system"
	SystemFunctionName = "automationTriggerSystem"
	SystemQueueSize    = 10
)

var _ ITrigger = (*SystemTrigger)(nil)

type SystemTrigger struct {
	baseTrigger
}

func NewSystemTrigger(eventBus event_bus.EventBus) ITrigger {
	return &SystemTrigger{
		baseTrigger{
			eventBus:     eventBus,
			msgQueue:     message_queue.New(SystemQueueSize),
			functionName: SystemFunctionName,
			name:         SystemName,
		},
	}
}

func (t *SystemTrigger) AsyncAttach(wg *sync.WaitGroup) {

	t.eventBus.Subscribe(TopicSystemStart, func(_ string, msg interface{}) {
		t.msgQueue.Publish(TopicSystem, map[string]interface{}{"event": EventStart,})
	})

	t.eventBus.Subscribe(TopicSystemStop, func(_ string, msg interface{}) {
		t.msgQueue.Publish(TopicSystem, map[string]interface{}{"event": EventStop,})
	})

	wg.Done()
}

func (t *SystemTrigger) Subscribe(options Subscriber) error {
	log.Infof("subscribe topic %s", TopicSystem)
	return t.msgQueue.Subscribe(TopicSystem, options.Handler)
}

func (t *SystemTrigger) Unsubscribe(options Subscriber) error {
	log.Infof("unsubscribe topic %s", TopicSystem)
	return t.msgQueue.Unsubscribe(TopicSystem, options.Handler)
}
