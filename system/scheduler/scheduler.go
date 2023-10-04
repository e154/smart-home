// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package scheduler

import (
	"context"
	"strconv"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/system/bus"

	"github.com/robfig/cron/v3"
	"go.uber.org/fx"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/logger"
)

type EntryID int

var (
	log = logger.MustGetLogger("scheduler")
)

type Scheduler struct {
	adaptors *adaptors.Adaptors
	cron     *cron.Cron
	eventBus bus.Bus
}

func NewScheduler(lc fx.Lifecycle,
	adaptors *adaptors.Adaptors,
	eventBus bus.Bus) (scheduler *Scheduler, err error) {
	scheduler = &Scheduler{
		adaptors: adaptors,
		eventBus: eventBus,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return scheduler.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return scheduler.Shutdown(ctx)
		},
	})

	return
}

func (c *Scheduler) Start(ctx context.Context) error {

	c.cron = cron.New(
		cron.WithSeconds(),
		cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
		)))

	// every hour
	c.cron.AddFunc("0 0 * * * *", func() {
		go func() {
			//log.Info("deleting obsolete metric entries ...")
			if err := c.adaptors.MetricBucket.DeleteOldest(context.Background(), c.getDays("clearMetricsDays", 60)); err != nil {
				log.Error(err.Error())
			}
		}()
		go func() {
			//log.Info("deleting obsolete log entries ...")
			if err := c.adaptors.Log.DeleteOldest(context.Background(), c.getDays("clearLogsDays", 60)); err != nil {
				log.Error(err.Error())
			}
		}()
		go func() {

			//log.Info("deleting obsolete entity storage entries ...")
			if err := c.adaptors.EntityStorage.DeleteOldest(context.Background(), c.getDays("clearEntityStorageDays", 60)); err != nil {
				log.Error(err.Error())
			}
		}()
		go func() {
			//log.Info("deleting obsolete run history entries ...")
			if err := c.adaptors.RunHistory.DeleteOldest(context.Background(), c.getDays("clearRunHistoryDays", 60)); err != nil {
				log.Error(err.Error())
			}
		}()
	})

	c.cron.Start()
	c.cron.Run()
	c.eventBus.Publish("system/services/scheduler", events.EventServiceStarted{Service: "Scheduler"})
	log.Info("started ...")

	return nil
}

func (c *Scheduler) Shutdown(_ context.Context) error {
	c.cron.Stop()
	c.eventBus.Publish("system/services/scheduler", events.EventServiceStopped{Service: "Scheduler"})
	log.Info("shutdown ...")
	return nil
}

func (c *Scheduler) AddFunc(spec string, cmd func()) (id EntryID, err error) {
	var entryId cron.EntryID
	if entryId, err = c.cron.AddFunc(spec, cmd); err != nil {
		return
	}
	id = EntryID(entryId)
	return
}

func (c *Scheduler) Remove(id EntryID) {
	c.cron.Remove(cron.EntryID(id))
}

func (c *Scheduler) getDays(varName string, def int) int {
	if clearMetricsDays, err := c.adaptors.Variable.GetByName(context.Background(), varName); err == nil {
		var days int
		if days, err = strconv.Atoi(clearMetricsDays.Value); err == nil {
			return days
		}
	}
	return def
}
