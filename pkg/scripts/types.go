// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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
	m "github.com/e154/smart-home/pkg/models"
)

// IScript ...
type IScript interface {
	Init() error
	Do() (string, error)
	AssertFunction(string, ...interface{}) (string, error)
	Compile() error
	PushStruct(string, interface{})
	PushFunction(string, interface{})
	EvalString(string) (string, error)
	CreateProgram(name, source string) (err error)
	RunProgram(name string) (result string, err error)
}

type Engine interface {
	Compile() (err error)
	PushStruct(name string, i interface{})
	PushFunction(name string, i interface{})
	EvalString(str ...string) (result string, errs error)
	EvalScript(script *m.Script) (result string, err error)
	DoFull() (res string, err error)
	Do() (string, error)
	AssertFunction(f string, arg ...interface{}) (result string, err error)
	Print(v ...interface{})
	Get() IScript
	File(path string) ([]byte, error)
	ScriptId() int64
	Script() *m.Script
}

type EngineWatcher interface {
	Stop()
	Spawn(f func(engine Engine))
	BeforeSpawn(f func(engine Engine))
	Engine() Engine
	PushStruct(name string, str interface{})
	PopStruct(name string)
	PushFunction(name string, f interface{})
	PopFunction(name string)
}

type EnginesWatcher interface {
	Stop()
	Spawn(f func(engine Engine))
	BeforeSpawn(f func(engine Engine))
	Engine() Engine
	AssertFunction(f string, arg ...interface{}) (result string, err error)
	PushStruct(name string, str interface{})
	PopStruct(name string)
	PushFunction(name string, f interface{})
	PopFunction(name string)
}

// ScriptService ...
type ScriptService interface {
	NewEngine(s *m.Script) (Engine, error)
	NewEngineWatcher(*m.Script) (EngineWatcher, error)
	NewEnginesWatcher([]*m.Script) (EnginesWatcher, error)
	PushStruct(name string, s interface{})
	PopStruct(name string)
	PushFunctions(name string, s interface{})
	PopFunction(name string)
	Restart()
}
