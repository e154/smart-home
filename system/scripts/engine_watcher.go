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

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
)

type EngineWatcher struct {
	eventBus      bus.Bus
	scriptService *scriptService
	f             func(engine *Engine)
	fBefore       func(engine *Engine)
	structures    *Pull
	functions     *Pull
	mx            *sync.RWMutex
	script        *m.Script
	engine        *Engine
}

func NewEngineWatcher(script *m.Script, s *scriptService, eventBus bus.Bus) *EngineWatcher {
	w := &EngineWatcher{
		eventBus:      eventBus,
		scriptService: s,
		mx:            &sync.RWMutex{},
		script:        script,
		structures:    NewPull(),
		functions:     NewPull(),
	}

	w.engine, _ = w.scriptService.NewEngine(nil)
	w.structures.Range(func(key, value interface{}) bool {
		w.engine.PushStruct(key.(string), value)
		return true
	})
	w.functions.Range(func(key, value interface{}) bool {
		w.engine.PushFunction(key.(string), value)
		return true
	})

	if script.Id != 0 {
		_ = eventBus.Subscribe(fmt.Sprintf("system/models/scripts/%d", script.Id), w.eventHandler)
	}

	return w
}

func (w *EngineWatcher) Stop() {
	if w.engine.model != nil && w.engine.model.Id != 0 {
		_ = w.eventBus.Unsubscribe(fmt.Sprintf("system/models/scripts/%d", w.engine.model.Id), w.eventHandler)
	}
}

func (w *EngineWatcher) Spawn(f func(engine *Engine)) {
	w.mx.RLock()
	defer w.mx.RUnlock()

	w.engine, _ = w.scriptService.NewEngine(&m.Script{
		Id:   w.script.Id,
		Lang: common.ScriptLangJavascript,
	})
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

	if _, err := w.engine.EvalScript(w.script); err != nil {
		if w.script.Id != 0 {
			log.Errorf("script id: %d, %s", w.script.Id, err.Error())
		}
		log.Error(err.Error())
	}

	if f != nil {
		w.f = f
		w.f(w.engine)
	}
}

func (w *EngineWatcher) BeforeSpawn(f func(engine *Engine)) {
	if f == nil {
		return
	}
	w.mx.Lock()
	defer w.mx.Unlock()
	w.fBefore = f
}

func (w *EngineWatcher) Engine() *Engine {
	w.mx.RLock()
	defer w.mx.RUnlock()
	return w.engine
}

// eventHandler ...
func (w *EngineWatcher) eventHandler(_ string, message interface{}) {

	switch msg := message.(type) {
	case events.EventUpdatedScriptModel:
		go w.eventUpdatedScript(msg)
	case events.EventRemovedScriptModel:
		go w.eventScriptDeleted(msg)
	}
}

func (w *EngineWatcher) eventUpdatedScript(msg events.EventUpdatedScriptModel) {

	if msg.Script == nil {
		return
	}

	w.script = msg.Script

	w.Spawn(w.f)

	log.Infof("script '%s' (%d) updated", msg.Script.Name, msg.ScriptId)
}

func (w *EngineWatcher) eventScriptDeleted(msg events.EventRemovedScriptModel) {
	if w.engine.model != nil {
		_ = w.eventBus.Unsubscribe(fmt.Sprintf("system/models/scripts/%d", w.script.Id), w.eventHandler)
	}

	var err error
	if w.engine, err = w.scriptService.NewEngine(nil); err != nil {
		log.Error(err.Error())
		return
	}
}

func (w *EngineWatcher) PushStruct(name string, str interface{}) {
	w.structures.Push(name, str)
	if w.engine != nil {
		w.engine.PushStruct(name, str)
	}
}

func (w *EngineWatcher) PopStruct(name string) {
	w.structures.Pop(name)
}

func (w *EngineWatcher) PushFunction(name string, f interface{}) {
	w.functions.Push(name, f)
	if w.engine != nil {
		w.engine.PushFunction(name, f)
	}
}

func (w *EngineWatcher) PopFunction(name string) {
	w.functions.Pop(name)
}
