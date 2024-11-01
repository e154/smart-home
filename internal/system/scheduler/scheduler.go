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

	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/logger"
	"github.com/e154/smart-home/pkg/scheduler"

	"github.com/robfig/cron/v3"
	"go.uber.org/fx"

	"github.com/e154/bus"
)

var (
	log = logger.MustGetLogger("scheduler")
)

var _ scheduler.Scheduler = (*Scheduler)(nil)

type Scheduler struct {
	adaptors    *adaptors.Adaptors
	cron        *cron.Cron
	eventBus    bus.Bus
	backupEntry cron.EntryID
}

func NewScheduler(lc fx.Lifecycle,
	adaptors *adaptors.Adaptors,
	eventBus bus.Bus) (scheduler.Scheduler, error) {
	scheduler := &Scheduler{
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

	return scheduler, nil
}

func (c *Scheduler) Start(ctx context.Context) error {

	c.cron = cron.New(
		cron.WithSeconds(),
		cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
		)))

	// every hour
	_, _ = c.cron.AddFunc("0 0 * * * *", func() {
		go func() {
			//log.Info("deleting obsolete metric entries ...")
			if err := c.adaptors.MetricBucket.DeleteOldest(context.Background(), c.getNumber("clearMetricsDays", 60)); err != nil {
				log.Error(err.Error())
			}
		}()
		go func() {
			//log.Info("deleting obsolete log entries ...")
			if err := c.adaptors.Log.DeleteOldest(context.Background(), c.getNumber("clearLogsDays", 60)); err != nil {
				log.Error(err.Error())
			}
		}()
		go func() {
			//log.Info("deleting obsolete entity storage entries ...")
			if err := c.adaptors.EntityStorage.DeleteOldest(context.Background(), c.getNumber("clearEntityStorageDays", 60)); err != nil {
				log.Error(err.Error())
			}
		}()
		go func() {
			//log.Info("deleting obsolete run history entries ...")
			if err := c.adaptors.RunHistory.DeleteOldest(context.Background(), c.getNumber("clearRunHistoryDays", 60)); err != nil {
				log.Error(err.Error())
			}
		}()
	})

	c.updateBackupScheduler()

	c.cron.Start()
	c.cron.Run()
	c.eventBus.Publish("system/services/scheduler", events.EventServiceStarted{Service: "Scheduler"})
	log.Info("started ...")

	_ = c.eventBus.Subscribe("system/models/variables/+", c.eventHandler)
	_ = c.eventBus.Subscribe("system/services/backup", c.eventHandler)

	return nil
}

func (c *Scheduler) Shutdown(_ context.Context) error {
	_ = c.eventBus.Unsubscribe("system/models/variables/+", c.eventHandler)
	_ = c.eventBus.Unsubscribe("system/services/backup", c.eventHandler)

	c.cron.Stop()
	c.eventBus.Publish("system/services/scheduler", events.EventServiceStopped{Service: "Scheduler"})
	log.Info("shutdown ...")
	return nil
}

func (c *Scheduler) AddFunc(spec string, cmd func()) (id scheduler.EntryID, err error) {
	var entryId cron.EntryID
	if entryId, err = c.cron.AddFunc(spec, cmd); err != nil {
		return
	}
	id = scheduler.EntryID(entryId)
	return
}

func (c *Scheduler) Remove(id scheduler.EntryID) {
	c.cron.Remove(cron.EntryID(id))
}

func (c *Scheduler) getNumber(varName string, def int) int {
	if variable, err := c.adaptors.Variable.GetByName(context.Background(), varName); err == nil {
		var num int
		if num, err = strconv.Atoi(variable.Value); err == nil {
			return num
		}
	}
	return def
}

func (c *Scheduler) getString(varName, def string) string {
	if variable, err := c.adaptors.Variable.GetByName(context.Background(), varName); err == nil {
		return variable.Value
	}
	return def
}

func (c *Scheduler) updateBackupScheduler() {

	if c.backupEntry != 0 {
		c.cron.Remove(c.backupEntry)
		c.backupEntry = 0
	}

	var err error
	c.backupEntry, err = c.cron.AddFunc(c.getString("createBackupAt", "0 0 0 * * *"), func() {
		c.eventBus.Publish("system/services/backup", events.CommandCreateBackup{
			Scheduler: true,
		})
	})
	if err != nil {
		log.Error(err.Error())
	}
}

func (c *Scheduler) eventHandler(_ string, message interface{}) {
	switch v := message.(type) {
	case events.EventUpdatedVariableModel:
		if v.Name == "createBackupAt" && v.Value != "" {
			c.updateBackupScheduler()
		}
	case events.EventCreatedBackup:
		c.eventBus.Publish("system/services/backup", events.CommandClearStorage{
			Num: int64(c.getNumber("maximumNumberOfBackups", 60)),
		})

		if !v.Scheduler {
			return
		}
		tgVariable, err := c.adaptors.Variable.GetByName(context.Background(), "sendbackuptoTelegramBot")
		if err != nil {
			return
		}
		if tgVariable.Value == "" {
			return
		}
		chunkSize := c.getNumber("sendTheBackupInPartsMb", 0)
		c.eventBus.Publish("system/services/backup", events.CommandSendFileToTelegram{
			Filename:  v.Name,
			EntityId:  common.EntityId(tgVariable.Value),
			Chunks:    chunkSize > 0,
			ChunkSize: chunkSize,
		})
	}
}
