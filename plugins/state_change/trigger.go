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

// EventStateChanged

package state_change

import (
	"sync"

	"github.com/e154/bus"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/plugins/triggers"
)

var _ triggers.ITrigger = (*Trigger)(nil)

type Trigger struct {
	eventBus     bus.Bus
	msgQueue     bus.Bus
	counter      *atomic.Int32
	functionName string
	name         string
}

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

func (t *Trigger) AsyncAttach(wg *sync.WaitGroup) {

	if err := t.eventBus.Subscribe("system/entities/+", t.eventHandler); err != nil {
		log.Error(err.Error())
	}

	wg.Done()
}

func (t *Trigger) eventHandler(_ string, event interface{}) {
	if t.counter.Load() <= 0 {
		return
	}
	switch v := event.(type) {
	case events.EventStateChanged:
		message := TriggerStateChangedMessage{
			StorageSave:     v.StorageSave,
			DoNotSaveMetric: v.DoNotSaveMetric,
			PluginName:      v.PluginName,
			EntityId:        v.EntityId,
			OldState: EventEntityState{
				EntityId:    v.OldState.EntityId,
				Value:       v.OldState.Value,
				State:       v.OldState.State,
				Attributes:  v.OldState.Attributes.Serialize(),
				Settings:    v.OldState.Settings.Serialize(),
				LastChanged: v.OldState.LastChanged,
				LastUpdated: v.OldState.LastUpdated,
			},
			NewState: EventEntityState{
				EntityId:    v.NewState.EntityId,
				Value:       v.NewState.Value,
				State:       v.NewState.State,
				Attributes:  v.NewState.Attributes.Serialize(),
				Settings:    v.NewState.Settings.Serialize(),
				LastChanged: v.NewState.LastChanged,
				LastUpdated: v.NewState.LastUpdated,
			},
		}
		t.msgQueue.Publish(v.EntityId.String(), message)
	}
}

// Subscribe ...
func (t *Trigger) Subscribe(options triggers.Subscriber) error {
	//log.Infof("subscribe topic %s", options.EntityId)
	t.counter.Inc()
	return t.msgQueue.Subscribe(options.EntityId.String(), options.Handler)
}

// Unsubscribe ...
func (t *Trigger) Unsubscribe(options triggers.Subscriber) error {
	//log.Infof("unsubscribe topic %s", options.EntityId)
	t.counter.Dec()
	return t.msgQueue.Unsubscribe(options.EntityId.String(), options.Handler)
}

// FunctionName ...
func (t *Trigger) FunctionName() string {
	return t.functionName
}
