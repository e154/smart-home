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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/message_queue"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"sync"
)

var (
	log = common.MustGetLogger("automation")
)

const (
	queueSize = 100
)

type Automation interface {
	Start() (err error)
	Shutdown() (err error)
	Reload()
	AddTask(model *m.Task)
	RemoveTask(model *m.Task)
}

type automation struct {
	eventBus      event_bus.EventBus
	taskLock      *sync.Mutex
	scriptService scripts.ScriptService
	tasks         map[int64]*Task
	taskCount     atomic.Uint64
	msgQueue      message_queue.MessageQueue
	entityManager entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	isStarted     *atomic.Bool
	rawPlugin     triggers.IGetTrigger
	pluginManager common.PluginManager
}

func NewAutomation(lc fx.Lifecycle,
	eventBus event_bus.EventBus,
	scriptService scripts.ScriptService,
	entityManager entity_manager.EntityManager,
	adaptors *adaptors.Adaptors,
	pluginManager common.PluginManager) (auto Automation) {

	auto = &automation{
		eventBus:      eventBus,
		taskLock:      &sync.Mutex{},
		scriptService: scriptService,
		tasks:         make(map[int64]*Task),
		msgQueue:      message_queue.New(queueSize),
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

func (a *automation) Start() (err error) {
	a.load()
	return
}

func (a *automation) Shutdown() (err error) {
	a.unload()
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
	for _, model := range tasks {
		a.AddTask(model)
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

	for id, task := range a.tasks {
		task.Stop()
		delete(a.tasks, id)
	}
	a.taskCount.Store(0)
	a.isStarted.Store(false)

	log.Info("Unloaded ...")
}

func (a *automation) Reload() {
	a.unload()
	a.load()
}

func (a *automation) AddTask(model *m.Task) {
	a.taskLock.Lock()
	defer a.taskLock.Unlock()

	task := NewTask(a, a.scriptService, model, a.entityManager, a.rawPlugin)
	a.taskCount.Inc()
	log.Infof("add task name(%s) id(%d)", task.Name(), task.Id())
	a.tasks[model.Id] = task
	task.Start()
}

func (a *automation) RemoveTask(model *m.Task) {
	a.taskLock.Lock()
	defer a.taskLock.Unlock()

	if task, ok := a.tasks[model.Id]; ok {
		task.Stop()
		delete(a.tasks, model.Id)
		a.taskCount.Dec()
	}
}
