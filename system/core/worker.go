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
	"go.uber.org/atomic"
	"sync"
	"time"
)

type Worker struct {
	Model       *m.Worker
	flow        *Flow
	CronTask    *cr.Task
	cron        *cr.Cron
	isRuning    atomic.Bool
	actionsLock sync.Mutex
	actions     map[int64]*Action
	cancelFunc  map[int64]context.CancelFunc
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
	w.CronTask = w.cron.NewTask(w.Model.Time, w.Do)
}

func (w *Worker) Stop() () {

	w.actionsLock.Lock()
	defer w.actionsLock.Unlock()

	if w.CronTask == nil {
		return
	}

	w.CronTask.Disable()
	w.cron.RemoveTask(w.CronTask)
	w.CronTask = nil

	w.unsafeRemoveActions()

	for {
		time.Sleep(time.Millisecond * 500)
		if !w.isRuning.Load() {
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
	w.actionsLock.Lock()
	defer w.actionsLock.Unlock()

	if _, ok := w.actions[action.Device.Id]; ok {
		return
	}

	w.actions[action.Device.Id] = action
}

func (w *Worker) unsafeRemoveActions() {

	for i, action := range w.actions {
		if cancel, ok := w.cancelFunc[action.Device.Id]; ok {
			cancel()
		}
		delete(w.actions, i)
	}
}

// Run worker script, and send result to flow as message struct
func (w *Worker) Do() {

	if w.isRuning.Load() || !w.flow.Node.IsConnected() {
		return
	}

	w.actionsLock.Lock()
	defer func() {
		w.isRuning.Store(false)
		w.actionsLock.Unlock()
	}()

	w.isRuning.Store(true)

	for _, action := range w.actions {
		w.doAction(action)
	}
}

func (w *Worker) doAction(action *Action) {

	//fmt.Println("<---- start", action.deviceAction.Name)
	//defer fmt.Println("end ---->", action.deviceAction.Name)

	if _, err := action.Do(); err != nil {
		log.Errorf("node: %s, device: %s error: %s", action.Node.Model().Name, action.Device.Name, err.Error())
		return
	}

	if w.flow.message.Error != "" {
		return
	}

	// create context
	ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	ctx = context.WithValue(ctx, "msg", w.flow.message.Copy())

	w.cancelFunc[action.Device.Id] = cancelFunc

	done := make(chan struct{})
	defer close(done)
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
	case <-ctx.Done():

	}

	if _, ok := w.cancelFunc[action.Device.Id]; ok {
		delete(w.cancelFunc, action.Device.Id)
	}

}
