package core

import (
	"../models"
	cr "github.com/e154/cron"
	"sync"
)

func NewWorker(model *models.Worker) (worker *Worker) {

	worker = &Worker{
		Model: model,
		actions:	make(map[int64]*Action),
	}

	return
}

type Worker struct {
	Model    *models.Worker
	CronTask *cr.Task
	mu       sync.Mutex
	actions  map[int64]*Action
}

func (w *Worker) RemoveTask() (err error) {

	if w.CronTask == nil {
		return
	}

	w.CronTask.Disable()

	// remove task from cron
	cron.RemoveTask(w.CronTask)

	return
}

func (w *Worker) AddAction(action *Action) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.actions[action.Device.Id]; ok {
		return
	}

	w.actions[action.Device.Id] = action
}

func (w *Worker) RegTask() {
	w.CronTask = cron.NewTask(w.Model.Time, func() {
		w.Do()
	})
}

func (w *Worker) Do() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, action := range w.actions {
		action.Do()
	}
}