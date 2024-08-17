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

package console

import (
	"github.com/dop251/goja"

	"github.com/e154/smart-home/system/scripts/require"
	"github.com/e154/smart-home/system/scripts/util"
)

const ModuleName = "console"

type Console struct {
	runtime *goja.Runtime
	util    *goja.Object
	printer Printer
}

type Printer interface {
	Log(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Debug(v ...interface{})
}

func (c *Console) log(p func(v ...interface{})) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if format, ok := goja.AssertFunction(c.util.Get("format")); ok {
			ret, err := format(c.util, call.Arguments...)
			if err != nil {
				panic(err)
			}

			p(ret.String())
		} else {
			panic(c.runtime.NewTypeError("util.format is not a function"))
		}

		return nil
	}
}

func Require(runtime *goja.Runtime, module *goja.Object) {
	requireWithPrinter(defaultStdPrinter)(runtime, module)
}

func RequireWithPrinter(printer Printer) require.ModuleLoader {
	return requireWithPrinter(printer)
}

func requireWithPrinter(printer Printer) require.ModuleLoader {
	return func(runtime *goja.Runtime, module *goja.Object) {
		c := &Console{
			runtime: runtime,
			printer: printer,
		}

		c.util = require.Require(runtime, util.ModuleName).(*goja.Object)

		o := module.Get("exports").(*goja.Object)
		o.Set("log", c.log(c.printer.Log))
		o.Set("error", c.log(c.printer.Error))
		o.Set("warn", c.log(c.printer.Warn))
		o.Set("info", c.log(c.printer.Log))
		o.Set("debug", c.log(c.printer.Log))
	}
}

func Enable(runtime *goja.Runtime) {
	runtime.Set("console", require.Require(runtime, ModuleName))
}

func init() {
	require.RegisterCoreModule(ModuleName, Require)
}
