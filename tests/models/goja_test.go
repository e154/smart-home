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

package models

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/e154/smart-home/adaptors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type Obj struct {
	Value string
}

func MyObject(call goja.ConstructorCall) *goja.Object {
	// call.This contains the newly created object as per http://www.ecma-international.org/ecma-262/5.1/index.html#sec-13.2.2
	// call.Arguments contain arguments passed to the function

	value := call.Arguments
	fmt.Println("---", value)


	o := &Obj{Value:"KKK"}
	call.This.Set("value", o.Value)


	//...

	// If return value is a non-nil *Object, it will be used instead of call.This
	// This way it is possible to return a Go struct or a map converted
	// into goja.Value using runtime.ToValue(), however in this case
	// instanceof will not work as expected.
	return call.This
}

func IC(call goja.ConstructorCall) *goja.Object {

	call.This.Set("Runmode", "debug")
	call.This.Set("hex2arr", nil)
	call.This.Set("CurrentNode", nil)
	call.This.Set("CurrentDevice", nil)

	return call.This
}

func TestGoja(t *testing.T) {

	const SCRIPT = `


var main = function() {
	print('main');
}

`

	Convey("goja", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors) {

			vm := goja.New()
			vm.Set("MyObject", MyObject)
			vm.Set("IC", IC)
			obj := &Obj{
				Value: "XXX",
			}
			vm.Set("MyObject2", obj)
			vm.Set("print", func(arg interface{}) {fmt.Println(arg)})

			p, err := goja.Compile("", SCRIPT, false)
			ctx.So(err, ShouldBeNil)

			//obj3 := &Obj{
			//	Value: "YYY",
			//}
			//vm.Set("MyObject3", obj3)

			result, err := vm.RunProgram(p)
			ctx.So(err, ShouldBeNil)


			fmt.Println(result)
			//
			//result, err = vm.RunProgram(p)
			//ctx.So(err, ShouldBeNil)

			//fmt.Println(result)

			//fmt.Println(">>>>>")
			//fmt.Println(obj.Value)
			//fmt.Println(">>>>>")

			//_, err = vm.RunString("print('???');")
			//_, err = vm.RunString("print(arg);")
			//ctx.So(err, ShouldBeNil)

			for i:=0;i<1000000;i++ {
				_, err = vm.RunString("main();")
				ctx.So(err, ShouldBeNil)
				time.Sleep(time.Microsecond * 1000)
			}


		})
	})
}
