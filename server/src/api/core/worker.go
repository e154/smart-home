package core

import (
	"../models"
	cr "github.com/e154/cron"
)

type Worker struct {
	Model     *models.Worker
	CronTasks map[int64]*cr.Task
	Devices   map[int64]*models.Device
}