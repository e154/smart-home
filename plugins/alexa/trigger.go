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

package alexa

import (
	"fmt"
	"sync"

	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/message_queue"
)

var _ triggers.ITrigger = (*Trigger)(nil)

const (
	// TriggerName ...
	TriggerName = "alexa"
	// TriggerFunctionName ...
	TriggerFunctionName = "automationTriggerAlexa"
	queueSize           = 10 //todo update
)

// Trigger ...
type Trigger struct {
	eventBus     event_bus.EventBus
	msgQueue     message_queue.MessageQueue
	functionName string
	name         string
}

// NewTrigger ...
func NewTrigger(eventBus event_bus.EventBus) (tr triggers.ITrigger) {
	return &Trigger{
		eventBus:     eventBus,
		msgQueue:     message_queue.New(queueSize),
		functionName: TriggerFunctionName,
		name:         TriggerName,
	}
}

// Name ...
func (t Trigger) Name() string {
	return t.name
}

// AsyncAttach ...
func (t Trigger) AsyncAttach(wg *sync.WaitGroup) {

	if err := t.eventBus.Subscribe(TopicPluginAlexa, t.eventHandler); err != nil {
		log.Error(err.Error())
	}

	wg.Done()
}

func (t *Trigger) eventHandler(topic string, msg interface{}) {
	switch v := msg.(type) {
	case EventAlexaAction:
		t.msgQueue.Publish(fmt.Sprintf("skill_%d", v.SkillId), v)
	}
}

// Subscribe ...
func (t Trigger) Subscribe(options triggers.Subscriber) error {
	if options.Payload == nil {
		return fmt.Errorf("trigger '%s' subscribe to empty topic", t.name)
	}
	topic := t.topic(options.Payload[TriggerOptionSkillId].Int64())
	log.Infof("trigger '%s' subscribe topic '%s'", t.name, topic)
	return t.msgQueue.Subscribe(topic, options.Handler)
}

// Unsubscribe ...
func (t Trigger) Unsubscribe(options triggers.Subscriber) error {
	if options.Payload == nil {
		return fmt.Errorf("trigger '%s' unsubscribe from empty topic", t.name)
	}
	topic := t.topic(options.Payload[TriggerOptionSkillId].Int64())
	log.Infof("trigger '%s' unsubscribe topic '%s'", t.name, topic)
	return t.msgQueue.Unsubscribe(topic, options.Handler)
}

// FunctionName ...
func (t Trigger) FunctionName() string {
	return t.functionName
}

func (t Trigger) topic(n int64) string {
	return fmt.Sprintf("skill_%d", n)
}

// CallManual ...
func (t Trigger) CallManual() {
	log.Warn("method not implemented")
}
