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

	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/scripts"

	"go.uber.org/atomic"
)

// NewCondition ...
func NewCondition(scriptService scripts.ScriptService,
	model *m.Condition) (condition *Condition, err error) {

	var scriptEngine scripts.EngineWatcher

	if model.Script != nil {
		if scriptEngine, err = scriptService.NewEngineWatcher(model.Script); err != nil {
			return
		}

		scriptEngine.PushStruct("Condition", NewConditionBind(condition))
		scriptEngine.Spawn(func(engine scripts.Engine) {
			//if _, err = engine.Do(); err != nil {
			//	return
			//}
		})
	}

	condition = &Condition{
		model:        model,
		lastStatus:   atomic.NewBool(false),
		scriptEngine: scriptEngine,
	}

	return
}

// Condition ...
type Condition struct {
	model        *m.Condition
	lastStatus   *atomic.Bool
	scriptEngine scripts.EngineWatcher
}

func (r *Condition) Stop() {
	if r.scriptEngine != nil {
		r.scriptEngine.Stop()
	}
}

// Check ...
func (r *Condition) Check(ctx context.Context) (result string, err error) {

	if r.scriptEngine != nil && r.scriptEngine.Engine() != nil {
		if result, err = r.scriptEngine.Engine().AssertFunction(ConditionFunc, ctx.Value("entityId")); err != nil {
			log.Error(err.Error())
		}
	}

	state := result == "true"
	r.lastStatus.Store(state)

	return
}

// Status ...
func (r *Condition) Status() bool {
	return r.lastStatus.Load()
}
