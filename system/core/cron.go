package core

import (
	cr "github.com/e154/smart-home/system/cron"
)

func NewCron() (cron *cr.Cron) {
	cron = cr.NewCron()
	cron.Run()
	return
}