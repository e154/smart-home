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

package triggers

import (
	"sync"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
)

const (
	// Name ...
	Name = "triggers"
	// TopicSystemStart ...
	TopicSystemStart = "system/event/start" //todo move
	// TopicSystemStop ...
	TopicSystemStop = "system/event/stop" //todo move

	Version = "0.0.1"
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

type SystemTriggerMessage struct {
	Topic     string      `json:"topic"`
	EventName string      `json:"event_name"`
	Event     interface{} `json:"event"`
}

type EventEntityState struct {
	EntityId    common.EntityId     `json:"entity_id"`
	Value       interface{}         `json:"value"`
	State       *events.EntityState `json:"state"`
	Attributes  m.AttributeValue    `json:"attributes"`
	Settings    m.AttributeValue    `json:"settings"`
	LastChanged *time.Time          `json:"last_changed"`
	LastUpdated *time.Time          `json:"last_updated"`
}

type TriggerStateChangedMessage struct {
	StorageSave     bool             `json:"storage_save"`
	DoNotSaveMetric bool             `json:"do_not_save_metric"`
	PluginName      string           `json:"plugin_name"`
	EntityId        common.EntityId  `json:"entity_id"`
	OldState        EventEntityState `json:"old_state"`
	NewState        EventEntityState `json:"new_state"`
}
