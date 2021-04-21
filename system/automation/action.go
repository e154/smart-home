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

package automation

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
	"sync"
)

type Action struct {
	model         *m.Action
	scriptService *scripts.ScriptService
	entityManager *entity_manager.EntityManager
	scriptEngine  *scripts.Engine
	inProcess     atomic.Bool
	sync.Mutex
}

func NewAction(scriptService *scripts.ScriptService,
	entityManager *entity_manager.EntityManager,
	model *m.Action) (action *Action, err error) {

	action = &Action{
		scriptService: scriptService,
		entityManager: entityManager,
		model:         model,
	}

	if action.scriptEngine, err = scriptService.NewEngine(model.Script); err != nil {
		return
	}

	if _, err = action.scriptEngine.Do(); err != nil {
		return
	}

	action.scriptEngine.PushStruct("Action", NewActionBind(action))

	return
}

func (a *Action) Run(entityId *common.EntityId) (result string, err error) {
	a.Lock()
	defer a.Unlock()

	if result, err = a.scriptEngine.AssertFunction(ActionFunc, entityId); err != nil {
		log.Error(err.Error())
	}

	return
}
