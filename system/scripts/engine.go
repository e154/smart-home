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
	"fmt"
	m "github.com/e154/smart-home/models"
	"io/ioutil"
)

type Magic interface {
	Init() error
	Do() (string, error)
	DoCustom(string) (string, error)
	Compile() error
	PushStruct(string, interface{}) (int, error)
	PushGlobalProxy(string, interface{}) int
	PushFunction(string, interface{}) (int, error)
	EvalString(string) (error)
	Close()
	Gc()
}

type Engine struct {
	model  *m.Script
	script Magic
	buf    []string
	IsRun  bool
	pull   *Pull
}

func (s *Engine) Compile() error {
	return s.script.Compile()
}

func (s *Engine) PushStruct(name string, i interface{}) (int, error) {
	return s.script.PushStruct(name, i)
}

func (s *Engine) Gc() {
	s.script.Gc()
}

func (s *Engine) PushGlobalProxy(name string, i interface{}) int {
	return s.script.PushGlobalProxy(name, i)
}

func (s *Engine) PushFunction(name string, i interface{}) (int, error) {
	return s.script.PushFunction(name, i)
}

func (s *Engine) EvalString(str string) error {
	return s.script.EvalString(str)
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

func (s *Engine) DoCustom(f string) (result string, err error) {

	if s.IsRun {
		return
	}

	s.IsRun = true
	result, err = s.script.DoCustom(f)

	// reset buffer
	s.buf = []string{}
	s.IsRun = false

	return
}

func (s *Engine) Print(v ...interface{}) {
	fmt.Println(v...)
	s.buf = append(s.buf, fmt.Sprint(v...))
}

func (s *Engine) Get() Magic {
	return s.script
}

func (s *Engine) File(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}
