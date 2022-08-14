package scheduler

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/fx"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/logger"
)

var (
	log = logger.MustGetLogger("scheduler")
)

type Scheduler struct {
	adaptors  *adaptors.Adaptors
	scheduler *gocron.Scheduler
}

func NewScheduler(lc fx.Lifecycle,
	adaptors *adaptors.Adaptors) (scheduler *Scheduler, err error) {
	scheduler = &Scheduler{adaptors: adaptors}

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
	c.scheduler = gocron.NewScheduler(time.UTC)

	// every day at 00:00 am
	_, _ = c.scheduler.Cron("0 0 * * *").Do(func() {
		if err := c.adaptors.MetricBucket.DeleteOldest(60); err != nil {
			log.Error(err.Error())
		}
	})

	c.scheduler.StartAsync()
	log.Info("started ...")

	return nil
}

func (c *Scheduler) Shutdown(_ context.Context) error {
	c.scheduler.Stop()
	log.Info("shutdown ...")
	return nil
}
