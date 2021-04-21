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

type SystemTrigger struct {
	baseTrigger
}

func NewSystemTrigger(eventBus *event_bus.EventBus) ITrigger {
	return &SystemTrigger{
		baseTrigger{
			eventBus:     eventBus,
			msgQueue:     message_queue.New(SystemQueueSize),
			functionName: SystemFunctionName,
			name:         SystemName,
		},
	}
}

func (s *SystemTrigger) AsyncAttach(wg *sync.WaitGroup) {

	s.eventBus.Subscribe(event_bus.TopicSystemStart, func(msg interface{}) {
		s.msgQueue.Publish(TopicSystem, map[string]interface{}{"event": EventStart,})
	})

	s.eventBus.Subscribe(event_bus.TopicSystemStop, func(msg interface{}) {
		s.msgQueue.Publish(TopicSystem, map[string]interface{}{"event": EventStop,})
	})

	wg.Done()
}

func (b *SystemTrigger) Subscribe(_ string, fn interface{}, _ interface{}) error {
	log.Infof("subscribe topic %s", TopicSystem)
	return b.msgQueue.Subscribe(TopicSystem, fn)
}

func (b *SystemTrigger) Unsubscribe(_ string, fn interface{}, _ interface{}) error {
	log.Infof("unsubscribe topic %s", TopicSystem)
	return b.msgQueue.Unsubscribe(TopicSystem, fn)
}
