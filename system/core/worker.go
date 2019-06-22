package core

import (
	m "github.com/e154/smart-home/models"
	cr "github.com/e154/smart-home/system/cron"
	"sync"
)

type Worker struct {
	Model    *m.Worker
	flow     *Flow
	CronTask *cr.Task
	cron     *cr.Cron
	sync.Mutex
	isRuning bool
	actions  map[int64]*Action
}

func NewWorker(model *m.Worker, flow *Flow, cron *cr.Cron) (worker *Worker) {

	worker = &Worker{
		Model:   model,
		flow:    flow,
		cron:    cron,
		actions: make(map[int64]*Action),
	}

	return
}

func (w *Worker) RemoveTask() (err error) {

	if w.CronTask == nil {
		return
	}

	w.CronTask.Disable()

	// remove task from cron
	w.cron.RemoveTask(w.CronTask)

	return
}

func (w *Worker) Actions() map[int64]*Action {
	w.Lock()
	defer w.Unlock()
	return w.actions
}

func (w *Worker) AddAction(action *Action) {
	w.Lock()
	defer w.Unlock()

	if _, ok := w.actions[action.Device.Id]; ok {
		return
	}

	w.actions[action.Device.Id] = action
}

func (w *Worker) RemoveActions() {
	w.Lock()
	w.actions = nil
	w.Unlock()
}

func (w *Worker) RegTask() {
	w.CronTask = w.cron.NewTask(w.Model.Time, func() {
		w.Do()
	})
}

// Run worker script, and send result to flow as message struct
func (w *Worker) Do() {

	if w.isRuning || !w.flow.Node.IsConnected {
		return
	}

	w.Lock()
	defer func() {
		w.isRuning = false
		w.Unlock()
	}()

	w.isRuning = true

	for _, action := range w.actions {
		if _, err := action.Do(); err != nil {
			log.Errorf("node: %s, device: %s error: %s", action.Node.Name, action.Device.Name, err.Error())
			continue
		}

		if action.Message.Error != "" {
			continue
		}

		message := action.Message.Copy()

		if err := w.flow.NewMessage(message); err != nil {
			log.Error(err.Error())
		}
	}
}
