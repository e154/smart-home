package core

import (
	"../models"
	cr "github.com/e154/cron"
)

type Workflow struct {
	model   *models.Workflow
	Flows   map[int64]*models.Flow
	Workers map[int64]*models.Worker
	CronTasks map[int64]*cr.Task
	Nodes   map[int64]*models.Node
}

func (wf *Workflow) Run() (err error) {

	return
}

func (wf *Workflow) Stop() (err error) {

	return
}