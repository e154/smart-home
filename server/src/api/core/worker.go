package core

import (
	"../models"
	cr "github.com/e154/cron"
	"sync"
	"../log"
)

func NewWorker(model *models.Worker, flow *Flow) (worker *Worker) {

	worker = &Worker{
		Model: model,
		flow:	flow,
		actions:	make(map[int64]*Action),
	}

	return
}

type Worker struct {
	Model    *models.Worker
	flow     *Flow
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

func (w *Worker) Actions() map[int64]*Action {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.actions
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

// Run worker script, and send result to flow as message struct
func (w *Worker) Do() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, action := range w.actions {
		res, err := action.Do()
		if err != nil {
			log.Error(err.Error())
		}

		message := &Message{
			Result: res,
			Flow: w.flow.Model,
			Device: action.Device,
			Node: action.Node,
			Device_state: func(state string) {
				action.SetDeviceState(state)
			},
		}

		if err = w.flow.NewMessage(message); err != nil {
			log.Error(err.Error())
		}
	}
}