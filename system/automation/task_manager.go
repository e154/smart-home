// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package automation

import (
	"context"
	"sync"

	"github.com/e154/bus"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/telemetry"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/plugins/triggers/types"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
)

type taskManager struct {
	eventBus      bus.Bus
	scriptService scripts.ScriptService
	supervisor    supervisor.Supervisor
	adaptors      *adaptors.Adaptors
	taskCount     *atomic.Uint64
	isStarted     *atomic.Bool
	rawPlugin     types.IGetTrigger
	sync.Mutex
	tasks map[int64]*Task
}

// NewTaskManager ...
func NewTaskManager(
	eventBus bus.Bus,
	scriptService scripts.ScriptService,
	sup supervisor.Supervisor,
	adaptors *adaptors.Adaptors) (manager *taskManager) {

	manager = &taskManager{
		eventBus:      eventBus,
		scriptService: scriptService,
		tasks:         make(map[int64]*Task),
		supervisor:    sup,
		adaptors:      adaptors,
		isStarted:     atomic.NewBool(false),
		taskCount:     atomic.NewUint64(0),
	}

	return
}

// Start ...
func (a *taskManager) Start() {
	a.load()
	_ = a.eventBus.Subscribe("system/models/tasks/+", a.eventHandler, false)
	_ = a.eventBus.Subscribe("system/automation/tasks/+", a.eventHandler, false)

	log.Info("Started")
}

// Shutdown ...
func (a *taskManager) Shutdown() {
	a.unload()
	_ = a.eventBus.Unsubscribe("system/models/tasks/+", a.eventHandler)
	_ = a.eventBus.Unsubscribe("system/automation/tasks/+", a.eventHandler)
	log.Info("shutdown")
}

func (a *taskManager) load() {
	if a.isStarted.Load() {
		return
	}

	// load triggers plugin
	plugin, err := a.supervisor.GetPlugin(triggers.Name)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if rawPlugin, ok := plugin.(types.IGetTrigger); ok {
		a.rawPlugin = rawPlugin
	} else {
		log.Fatal("bad static cast triggers.IGetTrigger")
	}

	a.AddTasks()

	a.isStarted.Store(true)

	log.Info("Loaded ...")
}

func (a *taskManager) unload() {
	if !a.isStarted.Load() {
		return
	}

	for id := range a.tasks {
		a.unloadTask(id)
	}
	a.isStarted.Store(false)

	log.Info("Unloaded ...")
}

// AddTasks ...
func (a *taskManager) AddTasks() {

	const perPage int64 = 500
	var page int64 = 0
LOOP:
	tasks, _, err := a.adaptors.Task.List(context.Background(), perPage, page*perPage, "", "", true)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, task := range tasks {
		a.addTask(task)
	}
	if len(tasks) != 0 {
		page++
		goto LOOP
	}
}

func (a *taskManager) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.CommandEnableTask:
		go a.updateTask(v.Id)
	case events.CommandDisableTask:
		go a.unloadTask(v.Id)

	case events.EventUpdatedTaskModel:
		go a.updateTask(v.Id)
	case events.EventCreatedTaskModel:
		go a.updateTask(v.Id)
	case events.EventRemovedTaskModel:
		go a.unloadTask(v.Id)
	}
}

// addTask ...
func (a *taskManager) addTask(model *m.Task) {
	task := NewTask(a.eventBus, a.scriptService, model)
	a.taskCount.Inc()
	log.Infof("add task name(%s) id(%d)", task.Name(), task.Id())
	a.Lock()
	a.tasks[model.Id] = task
	a.Unlock()
	task.Start()
}

func (a *taskManager) unloadTask(id int64) {
	a.Lock()
	defer a.Unlock()

	task, ok := a.tasks[id]
	if !ok {
		return
	}
	log.Infof("unload task %d", id)

	task.Stop()
	delete(a.tasks, id)
	a.taskCount.Dec()

}

func (a *taskManager) updateTask(id int64) {
	a.unloadTask(id)

	task, err := a.adaptors.Task.GetById(context.Background(), id)
	if err != nil {
		return
	}

	a.addTask(task)
}

func (a *taskManager) IsLoaded(id int64) (loaded bool) {
	a.Lock()
	_, loaded = a.tasks[id]
	a.Unlock()
	return
}

func (a *taskManager) TaskTelemetry(id int64) telemetry.Telemetry {
	a.Lock()
	defer a.Unlock()
	if task, ok := a.tasks[id]; ok {
		return task.Telemetry()
	}
	return nil
}
