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
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/system/bus"
	"sync"
)

type EngineWatcher struct {
	eventBus      bus.Bus
	scriptService *scriptService
	f             func(engine *Engine)
	mx            *sync.RWMutex
	engine        *Engine
}

func NewEngineWatched(engine *Engine, s *scriptService, eventBus bus.Bus) *EngineWatcher {
	w := &EngineWatcher{
		eventBus:      eventBus,
		scriptService: s,
		mx:            &sync.RWMutex{},
		engine:        engine,
	}

	if engine.model != nil && engine.model.Id != 0 {
		_ = eventBus.Subscribe(fmt.Sprintf("system/scripts/%d", engine.model.Id), w.eventHandler)
	}

	return w
}

func (w *EngineWatcher) Stop() {
	if w.engine.model != nil && w.engine.model.Id != 0 {
		_ = w.eventBus.Unsubscribe(fmt.Sprintf("system/scripts/%d", w.engine.model.Id), w.eventHandler)
	}
}

func (w *EngineWatcher) Spawn(f func(engine *Engine)) {
	if f == nil {
		return
	}
	w.mx.RLock()
	defer w.mx.RUnlock()
	w.f = f
	w.f(w.engine)
}

func (w *EngineWatcher) Engine() *Engine {
	w.mx.RLock()
	defer w.mx.RUnlock()
	return w.engine
}

// eventHandler ...
func (w *EngineWatcher) eventHandler(_ string, message interface{}) {

	switch msg := message.(type) {
	case events.EventUpdatedScript:
		go w.eventUpdatedScript(msg)
	case events.EventScriptDeleted:
		go w.eventScriptDeleted(msg)
	}
}

func (w *EngineWatcher) eventUpdatedScript(msg events.EventUpdatedScript) {

	if msg.Script == nil {
		return
	}

	w.mx.Lock()
	defer w.mx.Unlock()

	var err error
	w.engine, err = NewEngine(msg.Script, w.scriptService.functions, w.scriptService.structures)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if w.f == nil {
		return
	}
	w.f(w.engine)
	log.Infof("script '%s' (%d) updated", msg.Script.Name, msg.ScriptId)
}

func (w *EngineWatcher) eventScriptDeleted(msg events.EventScriptDeleted) {
	if w.engine.model != nil {
		_ = w.eventBus.Unsubscribe(fmt.Sprintf("system/scripts/%d", w.engine.model.Id), w.eventHandler)
	}

	var err error
	if w.engine, err = w.scriptService.NewEngine(nil); err != nil {
		log.Error(err.Error())
		return
	}
}
