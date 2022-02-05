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
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	. "github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/hashicorp/go-multierror"
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
	Close()
	CreateProgram(name, source string) (err error)
	RunProgram(name string) (result string, err error)
}

// Engine ...
type Engine struct {
	model      *m.Script
	script     IScript
	buf        []string
	IsRun      bool //todo fix
	functions  *Pull
	structures *Pull
}

// NewEngine ...
func NewEngine(s *m.Script, functions, structures *Pull) (engine *Engine, err error) {

	if s == nil {
		s = &m.Script{
			Lang: ScriptLangJavascript,
		}
	}

	engine = &Engine{
		model:      s,
		buf:        make([]string, 0),
		functions:  functions,
		structures: structures,
	}

	switch s.Lang {
	case ScriptLangTs, ScriptLangCoffee, ScriptLangJavascript:
		engine.script = NewJavascript(engine)
	default:
		err = errors.Wrap(ErrNotFound, string(s.Lang))
		return
	}

	err = engine.script.Init()

	return
}

// Compile ...
func (s *Engine) Compile() error {
	return s.script.Compile()
}

// PushStruct ...
func (s *Engine) PushStruct(name string, i interface{}) {
	s.script.PushStruct(name, i)
}

// PushFunction ...
func (s *Engine) PushFunction(name string, i interface{}) {
	s.script.PushFunction(name, i)
}

// EvalString ...
func (s *Engine) EvalString(str ...string) (result string, errs error) {
	var err error
	if len(str) == 0 {
		if result, err = s.script.Do(); err != nil {
			err = multierror.Append(err, errs)
		}
		return
	}
	for _, st := range str {
		if result, err = s.script.EvalString(st); err != nil {
			err = multierror.Append(err, errs)
		}
	}
	return
}

// EvalScript ...
func (s *Engine) EvalScript(script *m.Script) (result string, err error) {
	programName := strconv.Itoa(int(script.Id))
	if result, err = s.script.RunProgram(programName); err == nil {
		return
	}
	if errors.Is(err, ErrNotFound) {
		if err = s.script.CreateProgram(programName, script.Compiled); err != nil {
			err = errors.Wrap(ErrInternal, err.Error())
			return
		}
		result, err = s.script.RunProgram(programName)
	}
	return
}

// Close ...
func (s *Engine) Close() {
	s.script.Close()
}

// DoFull ...
func (s *Engine) DoFull() (res string, err error) {
	if s.IsRun {
		return
	}

	s.IsRun = true
	var result string
	result, err = s.script.Do()
	if err != nil {
		err = errors.Wrap(err, "do full")
		return
	}
	for _, r := range s.buf {
		res += r + "\n"
	}

	res += result + "\n"

	// reset buffer
	s.buf = []string{}
	s.IsRun = false

	return
}

// Do ...
func (s *Engine) Do() (string, error) {
	return s.script.Do()
}

// AssertFunction ...
func (s *Engine) AssertFunction(f string, arg ...interface{}) (result string, err error) {

	if s.IsRun {
		return
	}

	s.IsRun = true
	result, err = s.script.AssertFunction(f, arg...)

	// reset buffer
	s.buf = []string{}
	s.IsRun = false

	return
}

// Print ...
func (s *Engine) Print(v ...interface{}) {
	fmt.Println(v...)
	s.buf = append(s.buf, fmt.Sprint(v...))
}

// Get ...
func (s *Engine) Get() IScript {
	return s.script
}

// File ...
func (s *Engine) File(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}
