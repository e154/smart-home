// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package scripts

import (
	"fmt"
	"sync"

	"github.com/e154/smart-home/pkg/events"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/scripts"

	"github.com/pkg/errors"

	"github.com/e154/bus"
)

var _ scripts.EnginesWatcher = (*EnginesWatcher)(nil)

type EnginesWatcher struct {
	eventBus      bus.Bus
	scriptService *scriptService
	f             func(engine scripts.Engine)
	fBefore       func(engine scripts.Engine)
	structures    *Pull
	functions     *Pull
	mx            *sync.Mutex
	engine        scripts.Engine
	scripts       []*m.Script
}

func NewEnginesWatcher(scripts []*m.Script, s *scriptService, eventBus bus.Bus) *EnginesWatcher {
	w := &EnginesWatcher{
		eventBus:      eventBus,
		scriptService: s,
		mx:            &sync.Mutex{},
		scripts:       scripts,
		structures:    NewPull(),
		functions:     NewPull(),
	}

	for _, script := range scripts {
		_ = eventBus.Subscribe(fmt.Sprintf("system/models/scripts/%d", script.Id), w.eventHandler)
	}

	return w
}

func (w *EnginesWatcher) Stop() {
	w.mx.Lock()
	defer w.mx.Unlock()
	for _, script := range w.scripts {
		if script.Id != 0 {
			_ = w.eventBus.Unsubscribe(fmt.Sprintf("system/models/scripts/%d", script.Id), w.eventHandler)
		}
	}
}

func (w *EnginesWatcher) Spawn(f func(engine scripts.Engine)) {
	w.mx.Lock()
	defer w.mx.Unlock()

	w.engine, _ = w.scriptService.NewEngine(nil)
	w.structures.Range(func(key, value interface{}) bool {
		w.engine.PushStruct(key.(string), value)
		return true
	})
	w.functions.Range(func(key, value interface{}) bool {
		w.engine.PushFunction(key.(string), value)
		return true
	})

	if w.fBefore != nil {
		w.fBefore(w.engine)
	}

	for _, script := range w.scripts {
		if _, err := w.engine.EvalScript(script); err != nil {
			if script.Id != 0 {
				log.Errorf("script id: %d, %s", script.Id, err.Error())
			}
			log.Error(err.Error())
		}
	}

	if f != nil {
		w.f = f
		w.f(w.engine)
	}
}

func (w *EnginesWatcher) BeforeSpawn(f func(engine scripts.Engine)) {
	if f == nil {
		return
	}
	w.mx.Lock()
	defer w.mx.Unlock()
	w.fBefore = f
}

func (w *EnginesWatcher) Engine() scripts.Engine {
	w.mx.Lock()
	defer w.mx.Unlock()
	return w.engine
}

func (w *EnginesWatcher) AssertFunction(f string, arg ...interface{}) (result string, err error) {
	w.mx.Lock()
	defer w.mx.Unlock()
	if w.engine == nil {
		return
	}
	result, err = w.engine.AssertFunction(f, arg...)
	if err != nil {
		ids := []int64{}
		for _, script := range w.scripts {
			ids = append(ids, script.Id)
		}
		err = errors.Wrapf(err, "see scripts: %v ", ids)
	}
	return
}

// eventHandler ...
func (w *EnginesWatcher) eventHandler(_ string, message interface{}) {

	switch msg := message.(type) {
	case events.EventUpdatedScriptModel:
		go w.eventUpdatedScript(msg)
	case events.EventRemovedScriptModel:
		go w.eventScriptDeleted(msg)
	}
}

func (w *EnginesWatcher) eventUpdatedScript(msg events.EventUpdatedScriptModel) {

	if msg.Script == nil {
		return
	}

	for s, script := range w.scripts {
		if script.Id == msg.ScriptId {
			w.scripts[s] = msg.Script
			break
		}
	}

	w.Spawn(w.f)

	log.Infof("script '%s' (%d) updated", msg.Script.Name, msg.ScriptId)
}

func (w *EnginesWatcher) eventScriptDeleted(msg events.EventRemovedScriptModel) {
	if w.engine.Script() != nil {
		_ = w.eventBus.Unsubscribe(fmt.Sprintf("system/models/scripts/%d", msg.ScriptId), w.eventHandler)
	}

	var scriptName string

	var indexesToDelete []int

	for i, script := range w.scripts {
		if script.Id == msg.ScriptId {
			scriptName = script.Name
			indexesToDelete = append(indexesToDelete, i)
			break
		}
	}

	// remove script
	for i := len(indexesToDelete) - 1; i >= 0; i-- {
		index := indexesToDelete[i]
		w.scripts = append(w.scripts[:index], w.scripts[index+1:]...)
	}

	log.Infof("script '%s' (%d) deleted", scriptName, msg.ScriptId)

	w.Spawn(w.f)
}

func (w *EnginesWatcher) PushStruct(name string, str interface{}) {
	w.structures.Push(name, str)
	if w.engine != nil {
		w.engine.PushStruct(name, str)
	}
}

func (w *EnginesWatcher) PopStruct(name string) {
	w.structures.Pop(name)
}

func (w *EnginesWatcher) PushFunction(name string, f interface{}) {
	w.functions.Push(name, f)
	if w.engine != nil {
		w.engine.PushFunction(name, f)
	}
}

func (w *EnginesWatcher) PopFunction(name string) {
	w.functions.Pop(name)
}
