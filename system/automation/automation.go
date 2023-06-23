// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

//
//         AUTOMATION
// ┌──────┐  ┌──────┐  ┌──────┐
// │ TASK │  │ TASK │  │ TASK │
// └──────┘  └──────┘  └──────┘   ...
//

package automation

import (
	"context"
	"sync"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
	"go.uber.org/fx"
)

var (
	log = logger.MustGetLogger("automation")
)

const (
	queueSize = 100
)

// Automation ...
type Automation interface {
	Start() (err error)
	Shutdown() (err error)
	Reload()
	AddTask(model *m.Task)
	RemoveTask(model *m.Task)
}

type automation struct {
	eventBus      bus.Bus
	taskLock      *sync.Mutex
	scriptService scripts.ScriptService
	tasks         map[int64]*Task
	taskCount     atomic.Uint64
	msgQueue      bus.Bus
	entityManager entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	isStarted     *atomic.Bool
	rawPlugin     triggers.IGetTrigger
	pluginManager common.PluginManager
}

// NewAutomation ...
func NewAutomation(lc fx.Lifecycle,
	eventBus bus.Bus,
	scriptService scripts.ScriptService,
	entityManager entity_manager.EntityManager,
	adaptors *adaptors.Adaptors,
	pluginManager common.PluginManager) (auto Automation) {

	auto = &automation{
		eventBus:      eventBus,
		taskLock:      &sync.Mutex{},
		scriptService: scriptService,
		tasks:         make(map[int64]*Task),
		msgQueue:      bus.NewBus(),
		entityManager: entityManager,
		adaptors:      adaptors,
		isStarted:     atomic.NewBool(false),
		pluginManager: pluginManager,
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return auto.Shutdown()
		},
	})
	return
}

// Start ...
func (a *automation) Start() (err error) {
	a.load()
	_ = a.eventBus.Subscribe(bus.TopicAutomation, a.eventHandler)
	return
}

// Shutdown ...
func (a *automation) Shutdown() (err error) {
	a.unload()
	_ = a.eventBus.Unsubscribe(bus.TopicAutomation, a.eventHandler)
	return
}

func (a *automation) load() {
	if a.isStarted.Load() {
		return
	}

	// load triggers plugin
	plugin, err := a.pluginManager.GetPlugin(triggers.Name)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if rawPlugin, ok := plugin.(triggers.IGetTrigger); ok {
		a.rawPlugin = rawPlugin
	} else {
		log.Fatal("bad static cast triggers.IGetTrigger")
	}

	const perPage int64 = 100
	var page int64 = 0
LOOP:
	tasks, _, err := a.adaptors.Task.List(perPage, page*perPage, "", "", true)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, task := range tasks {
		a.AddTask(task)
	}
	if len(tasks) != 0 {
		page++
		goto LOOP
	}

	a.isStarted.Store(true)

	log.Info("Loaded ...")
}

func (a *automation) unload() {
	if !a.isStarted.Load() {
		return
	}

	for id := range a.tasks {
		a.removeTask(id)
	}
	a.isStarted.Store(false)

	log.Info("Unloaded ...")
}

// Reload ...
func (a *automation) Reload() {
	a.unload()
	a.load()
}

// AddTask ...
func (a *automation) AddTask(model *m.Task) {
	task := NewTask(a, a.scriptService, model, a.entityManager, a.rawPlugin)
	a.taskCount.Inc()
	log.Infof("add task name(%s) id(%d)", task.Name(), task.Id())
	a.taskLock.Lock()
	a.tasks[model.Id] = task
	a.taskLock.Unlock()
	task.Start()
}

// RemoveTask ...
func (a *automation) RemoveTask(model *m.Task) {
	a.removeTask(model.Id)
}

func (a *automation) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventCallTaskAction:
		go a.eventCallAction(v.Id, v.Name)
	case events.EventCallTaskTrigger:
		go a.eventCallTrigger(v.Id, v.Name)

	case events.EventEnableTask:
		go a.updateTask(v.Id)
	case events.EventUpdateTask:
		go a.updateTask(v.Id)
	case events.EventAddedTask:
		go a.updateTask(v.Id)

	case events.EventDisableTask:
		go a.removeTask(v.Id)
	case events.EventRemoveTask:
		go a.removeTask(v.Id)
	}
}

func (a *automation) eventCallTrigger(id int64, name string) {
	a.taskLock.Lock()
	defer a.taskLock.Unlock()

	if task, ok := a.tasks[id]; ok {
		task.CallTrigger(name)
	}
}

func (a *automation) eventCallAction(id int64, name string) {
	a.taskLock.Lock()
	defer a.taskLock.Unlock()

	if task, ok := a.tasks[id]; ok {
		task.CallAction(name)
	}
}

func (a *automation) removeTask(id int64) {
	a.taskLock.Lock()
	defer a.taskLock.Unlock()

	log.Infof("remove task %d", id)

	task, ok := a.tasks[id]
	if !ok {
		return
	}
	task.Stop()
	delete(a.tasks, id)
	a.taskCount.Dec()
}

func (a *automation) updateTask(id int64) {
	log.Infof("reload task %d", id)
	a.removeTask(id)

	task, err := a.adaptors.Task.GetById(id)
	if err != nil {
		return
	}

	a.AddTask(task)
}
