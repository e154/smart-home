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

// EventStateChanged

package triggers

import (
	"sync"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/system/bus"
)

const (
	// StateChangeName ...
	StateChangeName = "state_change"
	// StateChangeFunctionName ...
	StateChangeFunctionName = "automationTriggerStateChanged"
)

var _ ITrigger = (*StateChangeTrigger)(nil)

// StateChangeTrigger ...
type StateChangeTrigger struct {
	baseTrigger
}

// NewStateChangedTrigger ...
func NewStateChangedTrigger(eventBus bus.Bus) ITrigger {
	return &StateChangeTrigger{
		baseTrigger{
			eventBus:     eventBus,
			msgQueue:     bus.NewBus(),
			functionName: StateChangeFunctionName,
			name:         StateChangeName,
		},
	}
}

// AsyncAttach ...
func (t *StateChangeTrigger) AsyncAttach(wg *sync.WaitGroup) {

	if err := t.eventBus.Subscribe("system/entities/+", t.eventHandler); err != nil {
		log.Error(err.Error())
	}

	wg.Done()
}

func (t *StateChangeTrigger) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventStateChanged:
		t.msgQueue.Publish(string(v.EntityId), msg)
	}
}

// Subscribe ...
func (t *StateChangeTrigger) Subscribe(options Subscriber) error {
	//log.Infof("subscribe topic %s", options.EntityId)
	return t.msgQueue.Subscribe(options.EntityId.String(), options.Handler)
}

// Unsubscribe ...
func (t *StateChangeTrigger) Unsubscribe(options Subscriber) error {
	//log.Infof("unsubscribe topic %s", options.EntityId)
	return t.msgQueue.Unsubscribe(options.EntityId.String(), options.Handler)
}
