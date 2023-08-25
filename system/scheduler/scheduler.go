package scheduler

import (
	"context"
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

func (c *Scheduler) Start(_ context.Context) error {

	c.cron = cron.New(
		cron.WithSeconds(),
		cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
		)))

	// every day at 00:00 am
	c.cron.AddFunc("0 0 0 * * *", func() {
		go func() {
			log.Info("deleting obsolete metric entries ...")
			if err := c.adaptors.MetricBucket.DeleteOldest(60); err != nil {
				log.Error(err.Error())
			}
		}()
		go func() {
			log.Info("deleting obsolete log entries ...")
			if err := c.adaptors.Log.DeleteOldest(60); err != nil {
				log.Error(err.Error())
			}
		}()
		go func() {
			log.Info("deleting obsolete entity storage entries ...")
			if err := c.adaptors.EntityStorage.DeleteOldest(60); err != nil {
				log.Error(err.Error())
			}
		}()
	})

	c.cron.Start()
	c.cron.Run()
	c.eventBus.Publish("system/services/scheduler", events.EventServiceStarted{})
	log.Info("started ...")

	return nil
}

func (c *Scheduler) Shutdown(_ context.Context) error {
	c.cron.Stop()
	c.eventBus.Publish("system/services/scheduler", events.EventServiceStopped{})
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
