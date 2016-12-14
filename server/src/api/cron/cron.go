package cron

import (
	"github.com/e154/cron"
	"../log"
)

// Singleton
var instantiated *cron.Cron = nil

func Cron() *cron.Cron {
	return instantiated
}

func Initialize() {
	log.Info("Crontab initialize...")

	if instantiated == nil {
		instantiated = cron.NewCron()
		instantiated.Run()
	}

	return
}