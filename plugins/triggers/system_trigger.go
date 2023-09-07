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
	"sync"

	"github.com/e154/smart-home/system/bus"
)

const (
	// TopicSystem ...
	TopicSystem = "system"
	// EventStart ...
	EventStart = "START"
	// EventStop ...
	EventStop = "STOP"
	// SystemName ...
	SystemName = "system"
	// SystemFunctionName ...
	SystemFunctionName = "automationTriggerSystem"
	// SystemQueueSize ...
	SystemQueueSize = 10
)

var _ ITrigger = (*SystemTrigger)(nil)

// SystemTrigger ...
type SystemTrigger struct {
	baseTrigger
}

// NewSystemTrigger ...
func NewSystemTrigger(eventBus bus.Bus) ITrigger {
	return &SystemTrigger{
		baseTrigger{
			eventBus:     eventBus,
			msgQueue:     bus.NewBus(),
			functionName: SystemFunctionName,
			name:         SystemName,
		},
	}
}

// AsyncAttach ...
func (t *SystemTrigger) AsyncAttach(wg *sync.WaitGroup) {

	_ = t.eventBus.Subscribe(TopicSystemStart, func(_ string, msg interface{}) {
		t.msgQueue.Publish(TopicSystem, map[string]interface{}{"event": EventStart})
	})

	_ = t.eventBus.Subscribe(TopicSystemStop, func(_ string, msg interface{}) {
		t.msgQueue.Publish(TopicSystem, map[string]interface{}{"event": EventStop})
	})

	wg.Done()
}

// Subscribe ...
func (t *SystemTrigger) Subscribe(options Subscriber) error {
	log.Infof("subscribe topic %s", TopicSystem)
	return t.msgQueue.Subscribe(TopicSystem, options.Handler)
}

// Unsubscribe ...
func (t *SystemTrigger) Unsubscribe(options Subscriber) error {
	log.Infof("unsubscribe topic %s", TopicSystem)
	return t.msgQueue.Unsubscribe(TopicSystem, options.Handler)
}
