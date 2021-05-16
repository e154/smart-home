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
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/message_queue"
	"sync"
)

var _ triggers.ITrigger = (*Trigger)(nil)

const (
	TriggerName         = "alexa"
	TriggerFunctionName = "automationTriggerAlexa"
	queueSize           = 10 //todo update
)

type Trigger struct {
	eventBus     event_bus.EventBus
	msgQueue     message_queue.MessageQueue
	functionName string
	name         string
}

func NewTrigger(eventBus event_bus.EventBus) (tr triggers.ITrigger) {
	return &Trigger{
		eventBus:     eventBus,
		msgQueue:     message_queue.New(queueSize),
		functionName: TriggerFunctionName,
		name:         TriggerName,
	}
}

func (t Trigger) Name() string {
	return t.name
}

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

func (t Trigger) Subscribe(options triggers.Subscriber) error {
	if options.Payload == nil {
		return fmt.Errorf("trigger '%s' subscribe to empty topic", t.name)
	}
	log.Infof("trigger '%s' subscribe topic '%s'", t.name, t.topic(options.Payload))
	return t.msgQueue.Subscribe(t.topic(options.Payload), options.Handler)
}

func (t Trigger) Unsubscribe(options triggers.Subscriber) error {
	if options.Payload == nil {
		return fmt.Errorf("trigger '%s' unsubscribe from empty topic", t.name)
	}
	log.Infof("trigger '%s' unsubscribe topic '%s'", t.name, t.topic(options.Payload))
	return t.msgQueue.Unsubscribe(t.topic(options.Payload), options.Handler)
}

func (t Trigger) FunctionName() string {
	return t.functionName
}

func (t Trigger) topic(n interface{}) string {
	return fmt.Sprintf("skill_%v", n)
}
