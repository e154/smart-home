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

package scripts

import (
	"errors"
	"fmt"
	. "github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/scripts/bind"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("scripts")
)

type ScriptService struct {
	cfg  *config.AppConfig
	pull *Pull
}

func NewScriptService(cfg *config.AppConfig) (service *ScriptService) {

	pull := &Pull{
		functions:  make(map[string]interface{}),
		structures: make(map[string]interface{}),
	}

	service = &ScriptService{
		cfg:  cfg,
		pull: pull,
	}

	service.PushStruct("Log", &bind.LogBind{})
	service.PushFunctions("ExecuteSync", bind.ExecuteSync)
	service.PushFunctions("ExecuteAsync", bind.ExecuteAsync)
	return service
}

func (service ScriptService) NewEngine(s *m.Script) (engine *Engine, err error) {

	engine = &Engine{
		model: s,
		buf:   make([]string, 0),
		pull:  service.pull,
	}

	switch s.Lang {
	case ScriptLangTs, ScriptLangCoffee, ScriptLangJavascript:
		engine.script = &Javascript{engine: engine}
	default:
		err = errors.New(fmt.Sprintf("undefined language %s", s.Lang))
		return
	}

	//if err == nil {
	//	log.Infof("Add script: %s (%s)", s.Name, s.Lang)
	//}

	engine.script.Init()

	return
}

func (service *ScriptService) PushStruct(name string, s interface{}) {
	service.pull.Lock()
	defer service.pull.Unlock()

	service.pull.structures[name] = s
}

func (service *ScriptService) PushFunctions(name string, s interface{}) {
	service.pull.Lock()
	defer service.pull.Unlock()

	//fmt.Println("push function")

	service.pull.functions[name] = s
}
