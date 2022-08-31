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

package scripts

import (
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/common/web"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts/bind"
	"github.com/e154/smart-home/system/storage"
)

var (
	log = logger.MustGetLogger("scripts")
)

// ScriptService ...
type ScriptService interface {
	NewEngine(s *m.Script) (*Engine, error)
	PushStruct(name string, s interface{})
	PushFunctions(name string, s interface{})
	Purge()
}

// scriptService ...
type scriptService struct {
	cfg        *m.AppConfig
	functions  *Pull
	structures *Pull
	storage    *storage.Storage
	crawler    web.Crawler
}

// NewScriptService ...
func NewScriptService(cfg *m.AppConfig,
	storage *storage.Storage,
	crawler web.Crawler) ScriptService {

	s := &scriptService{
		cfg:        cfg,
		functions:  NewPull(),
		structures: NewPull(),
		storage:    storage,
		crawler:    crawler,
	}

	s.bind()

	return s
}

// NewEngine ...
func (s *scriptService) NewEngine(scr *m.Script) (*Engine, error) {
	return NewEngine(scr, s.structures, s.functions)
}

// PushStruct ...
func (s *scriptService) PushStruct(name string, str interface{}) {
	log.Infof("register structure: '%s'", name)
	s.structures.Add(name, str)
}

// PushFunctions ...
func (s *scriptService) PushFunctions(name string, f interface{}) {
	log.Infof("register function: '%s'", name)
	s.functions.Add(name, f)
}

// Purge ...
func (s *scriptService) Purge() {
	s.functions.Purge()
	s.structures.Purge()
	s.bind()
}

func (s *scriptService) bind() {
	s.PushStruct("Log", &bind.LogBind{})
	s.PushFunctions("ExecuteSync", bind.ExecuteSync)
	s.PushFunctions("ExecuteAsync", bind.ExecuteAsync)
	s.PushStruct("Storage", bind.NewStorageBind(s.storage))
	s.PushStruct("http", bind.NewHttpBind(s.crawler))
	s.PushStruct("HTTP", bind.NewHttpBind(s.crawler))
}
