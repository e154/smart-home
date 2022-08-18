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

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	// Name ...
	Name = "triggers"
	// TopicSystemStart ...
	TopicSystemStart = "system_start"
	// TopicSystemStop ...
	TopicSystemStop = "system_stop"
)

// IGetTrigger ...
type IGetTrigger interface {
	GetTrigger(string) (ITrigger, error)
}

// IRegistrar ...
type IRegistrar interface {
	RegisterTrigger(ITrigger) error
	UnregisterTrigger(string) error
	TriggerList() []string
}

// todo deAttach
type ITrigger interface {
	Name() string
	AsyncAttach(wg *sync.WaitGroup)
	Subscribe(Subscriber) error
	Unsubscribe(Subscriber) error
	FunctionName() string
	CallManual()
}

// Subscriber ...
type Subscriber struct {
	EntityId *common.EntityId
	Handler  interface{}
	Payload  m.Attributes
}

const (
	// CronOptionTrigger ...
	CronOptionTrigger = "cron"
)
