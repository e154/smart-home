// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package core

import (
	"context"
	m "github.com/e154/smart-home/models"
	cr "github.com/e154/smart-home/system/cron"
	"sync"
	"time"
)

type Worker struct {
	Model    *m.Worker
	flow     *Flow
	CronTask *cr.Task
	cron     *cr.Cron
	sync.Mutex
	isRuning   bool
	actions    map[int64]*Action
	cancelFunc map[int64]context.CancelFunc
}

func NewWorker(model *m.Worker, flow *Flow, cron *cr.Cron) (worker *Worker) {

	worker = &Worker{
		Model:      model,
		flow:       flow,
		cron:       cron,
		actions:    make(map[int64]*Action),
		cancelFunc: make(map[int64]context.CancelFunc),
	}

	return
}

func (w *Worker) Start() {
	w.CronTask = w.cron.NewTask(w.Model.Time, func() {
		w.Do()
	})
}

func (w *Worker) Stop() () {

	w.Lock()
	defer w.Unlock()

	if w.CronTask == nil {
		return
	}

	w.CronTask.Disable()
	w.cron.RemoveTask(w.CronTask)
	w.CronTask = nil

	w.removeActions()

	for {
		time.Sleep(time.Millisecond * 500)
		if !w.isRuning {
			log.Infof("worker %v ... ok", w.Model.Id)
			break
		}

		select {
		case <-time.After(3 * time.Second):
			return
		default:

		}
	}

	return
}

func (w *Worker) AddAction(action *Action) {
	w.Lock()
	defer w.Unlock()

	if _, ok := w.actions[action.Device.Id]; ok {
		return
	}

	w.actions[action.Device.Id] = action
}

func (w *Worker) removeActions() {

	for i, action := range w.actions {
		if cancel, ok := w.cancelFunc[action.Device.Id]; ok {
			cancel()
		}
		delete(w.actions, i)
	}
}

// Run worker script, and send result to flow as message struct
func (w *Worker) Do() {

	w.Lock()
	if w.isRuning || !w.flow.Node.IsConnected() {
		w.Unlock()
		return
	}

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

		if w.flow.message.Error != "" {
			continue
		}

		// create context
		ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
		ctx = context.WithValue(ctx, "msg", w.flow.message.Copy())

		w.cancelFunc[action.Device.Id] = cancelFunc

		done := make(chan struct{})
		go func() {
			if err := w.flow.NewMessage(ctx); err != nil {
				//log.Errorf("flow '%v' end with error: '%+v'", action.flow.Model.Name, err.Error())
			}

			if ctx.Err() != nil {
				//log.Errorf("flow '%v' end with error: '%+v'", action.flow.Model.Name, ctx.Err())
			}

			done <- struct{}{}
		}()

		select {
		case <-done:
			close(done)
		case <-ctx.Done():

		}

		if _, ok := w.cancelFunc[action.Device.Id]; ok {
			delete(w.cancelFunc, action.Device.Id)
		}
	}
}
