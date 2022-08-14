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

package logs

import (
	"fmt"
	m "github.com/e154/smart-home/models"
	"sync"

	"github.com/rcrowley/go-metrics"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/cron"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/event_bus/events"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	cores         int64
	model         string
	ErrTotal      metrics.Counter
	ErrToday      metrics.Counter
	ErrYesterday  metrics.Counter
	WarnTotal     metrics.Counter
	WarnToday     metrics.Counter
	WarnYesterday metrics.Counter
	eventBus      event_bus.EventBus
	updateLock    *sync.Mutex
	cron          *cron.Cron
}

// NewActor ...
func NewActor(entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus, entity *m.Entity) *Actor {

	actor := &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:                common.EntityId(fmt.Sprintf("%s.%s", EntityLogs, Name)),
			Name:              Name,
			EntityType:        EntityLogs,
			UnitOfMeasurement: "",
			AttrMu:            &sync.RWMutex{},
			Attrs:             NewAttr(),
			Manager:           entityManager,
		},
		eventBus:      eventBus,
		ErrTotal:      metrics.NewCounter(),
		ErrToday:      metrics.NewCounter(),
		ErrYesterday:  metrics.NewCounter(),
		WarnTotal:     metrics.NewCounter(),
		WarnToday:     metrics.NewCounter(),
		WarnYesterday: metrics.NewCounter(),
		updateLock:    &sync.Mutex{},
	}

	if entity != nil {
		actor.Metric = entity.Metrics
		attrs := entity.Attributes
		actor.ErrTotal.Inc(attrs[AttrErrTotal].Int64())
		actor.ErrToday.Inc(attrs[AttrErrToday].Int64())
		actor.ErrYesterday.Inc(attrs[AttrErrYesterday].Int64())
		actor.WarnTotal.Inc(attrs[AttrWarnTotal].Int64())
		actor.WarnToday.Inc(attrs[AttrWarnToday].Int64())
		actor.WarnYesterday.Inc(attrs[AttrWarnYesterday].Int64())
	}
	return actor
}

func (e *Actor) Spawn() entity_manager.PluginActor {
	return e
}

func (u *Actor) selfUpdate() {

	u.updateLock.Lock()
	defer u.updateLock.Unlock()

	oldState := u.GetEventState(u)
	u.Now(oldState)

	u.AttrMu.Lock()
	u.Attrs[AttrErrTotal].Value = u.ErrTotal.Count()
	u.Attrs[AttrErrToday].Value = u.ErrToday.Count()
	u.Attrs[AttrErrYesterday].Value = u.ErrYesterday.Count()
	u.Attrs[AttrWarnTotal].Value = u.WarnTotal.Count()
	u.Attrs[AttrWarnToday].Value = u.WarnToday.Count()
	u.Attrs[AttrWarnYesterday].Value = u.WarnYesterday.Count()
	u.AttrMu.Unlock()

	u.eventBus.Publish(event_bus.TopicEntities, events.EventStateChanged{
		StorageSave: true,
		PluginName:  u.Id.PluginName(),
		EntityId:    u.Id,
		OldState:    oldState,
		NewState:    u.GetEventState(u),
	})
}

func (u *Actor) LogsHook(level common.LogLevel) {

	switch level {
	case common.LogLevelError:
		u.ErrTotal.Inc(1)
		u.ErrToday.Inc(1)
	case common.LogLevelWarning:
		u.WarnTotal.Inc(1)
		u.WarnToday.Inc(1)
	case common.LogLevelInfo:
	case common.LogLevelDebug:

	}
	u.selfUpdate()
}

func (u *Actor) UpdateDay() {
	fmt.Println("update day")
	u.ErrYesterday.Clear()
	u.ErrYesterday.Inc(u.ErrToday.Count())
	u.WarnYesterday.Clear()
	u.WarnYesterday.Inc(u.WarnToday.Count())
	u.ErrToday.Clear()
	u.WarnToday.Clear()

	u.selfUpdate()
}
