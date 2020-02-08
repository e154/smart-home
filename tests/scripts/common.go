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
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
)

var store = func(interface{}) {}

func storeRegisterCallback(scriptService *scripts.ScriptService) {
	scriptService.PushFunctions("store", func(value interface{}) {
		store(value)
	})
}

type MyStruct struct {
	Bool    bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	UInt    uint
	UInt8   uint8
	UInt16  uint16
	UInt32  uint32
	UInt64  uint64
	String  string
	Bytes   []byte
	Float32 float32
	Float64 float64
	Empty   *MyStruct
	Nested  *MyStruct
	Slice   []int
	private int
}

func (m *MyStruct) Multiply(x int) int {
	return m.Int * x
}

func (m *MyStruct) privateMethod() int {
	return 1
}

func GetScripts(ctx C, scriptService *scripts.ScriptService, adaptors *adaptors.Adaptors, args ...int) (scripts map[string]*m.Script) {

	scripts = make(map[string]*m.Script)
	for _, arg := range args {
		script := &m.Script{
			Lang:        "coffeescript",
			Name:        fmt.Sprintf("test%d", arg),
			Source:      coffeeScripts[fmt.Sprintf("coffeeScript%d", arg)],
			Description: "test",
		}

		engine, err := scriptService.NewEngine(script)
		ctx.So(err, ShouldBeNil)
		err = engine.Compile()
		ctx.So(err, ShouldBeNil)
		scriptId, err := adaptors.Script.Add(script)
		ctx.So(err, ShouldBeNil)
		script, err = adaptors.Script.GetById(scriptId)
		ctx.So(err, ShouldBeNil)
		scripts[fmt.Sprintf("script%d", arg)] = script
	}

	return
}
