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
		//TODO refactor message system
		if _, err := action.Do(); err != nil {
			log.Errorf("node: %s, device: %s error: %s", action.Node.Name, action.Device.Name, err.Error())
			continue
		}
		//TODO refactor message system
		if action.Message.Error != "" {
			continue
		}
		//TODO refactor message system
		message := NewMessage()
		*message = *action.Message
		message.Flow = w.flow.Model
		message.Device_state = func(state string) {
			action.SetDeviceState(state)
		}

		if err := w.flow.NewMessage(message); err != nil {
			log.Error(err.Error())
		}
	}
}