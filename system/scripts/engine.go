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
	"io/ioutil"
	"strconv"
)

type IScript interface {
	Init() error
	Do() (string, error)
	AssertFunction(string) (string, error)
	Compile() error
	PushStruct(string, interface{})
	PushFunction(string, interface{})
	EvalString(string) (string, error)
	Close()
	CreateProgram(name, source string) (err error)
	RunProgram(name string) (result string, err error)
}

type Engine struct {
	model      *m.Script
	script     IScript
	buf        []string
	IsRun      bool
	functions  *Pull
	structures *Pull
}

func NewEngine(s *m.Script, functions, structures *Pull) (engine *Engine, err error) {

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
		err = errors.New(fmt.Sprintf("undefined language %s", s.Lang))
		return
	}

	err = engine.script.Init()

	return
}

func (s *Engine) Compile() error {
	return s.script.Compile()
}

func (s *Engine) PushStruct(name string, i interface{}) {
	s.script.PushStruct(name, i)
}

func (s *Engine) PushFunction(name string, i interface{}) {
	s.script.PushFunction(name, i)
}

func (s *Engine) EvalString(str string) (result string, err error) {
	result, err = s.script.EvalString(str)
	return
}

func (s *Engine) EvalScript(script *m.Script) (result string, err error) {
	programName := strconv.Itoa(int(script.Id))
	if result, err = s.script.RunProgram(programName); err == nil {
		return
	}

	if err == ErrorProgramNotFound {
		if err = s.script.CreateProgram(programName, script.Compiled); err != nil {
			return
		}
		result, err = s.script.RunProgram(programName)
	}
	return
}

func (s *Engine) Close() {
	s.script.Close()
}

func (s *Engine) DoFull() (res string, err error) {
	if s.IsRun {
		return
	}

	s.IsRun = true
	var result string
	result, err = s.script.Do()
	for _, r := range s.buf {
		res += r + "\n"
	}

	res += result + "\n"

	// reset buffer
	s.buf = []string{}
	s.IsRun = false

	return
}

func (s *Engine) Do() (string, error) {
	return s.script.Do()
}

func (s *Engine) AssertFunction(f string) (result string, err error) {

	if s.IsRun {
		return
	}

	s.IsRun = true
	result, err = s.script.AssertFunction(f)

	// reset buffer
	s.buf = []string{}
	s.IsRun = false

	return
}

func (s *Engine) Print(v ...interface{}) {
	fmt.Println(v...)
	s.buf = append(s.buf, fmt.Sprint(v...))
}

func (s *Engine) Get() IScript {
	return s.script
}

func (s *Engine) File(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}
