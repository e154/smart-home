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
)

type baseTrigger struct {
	eventBus     event_bus.EventBus
	msgQueue     message_queue.MessageQueue
	functionName string
	name         string
}

func (b *baseTrigger) Name() string {
	return b.name
}

func (b *baseTrigger) FunctionName() string {
	return b.functionName
}

func (b *baseTrigger) Subscribe(topic string, fn interface{}, _ interface{}) error {
	log.Infof("subscribe topic %s", topic)
	return b.msgQueue.Subscribe(topic, fn)
}

func (b *baseTrigger) Unsubscribe(topic string, fn interface{}, _ interface{}) error {
	log.Infof("unsubscribe topic %s", topic)
	return b.msgQueue.Unsubscribe(topic, fn)
}
