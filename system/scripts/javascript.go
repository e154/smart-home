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
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	. "github.com/e154/smart-home/common"
	"strings"
	"sync"
)

var (
	ErrorProgramNotFound = errors.New("program not found")
)

type Javascript struct {
	engine       *Engine
	compiler     string
	vm           *goja.Runtime
	program      *goja.Program
	lockPrograms sync.Mutex
	programs     map[string]*goja.Program
}

func NewJavascript(engine *Engine) *Javascript {
	return &Javascript{
		engine:   engine,
		programs: make(map[string]*goja.Program),
	}
}

func (j *Javascript) Init() (err error) {

	j.vm = goja.New()
	registry := new(require.Registry) // this can be shared by multiple runtimes
	registry.Enable(j.vm)

	j.bind()

	if j.engine.model.Compiled == "" {
		return
	}

	if j.program, err = goja.Compile("", j.engine.model.Compiled, false); err != nil {
		log.Error(err.Error())
	}

	return
}

func (j *Javascript) Close() {

}

func (j *Javascript) Compile() (err error) {

	if err = j.GetCompiler(); err != nil {
		return
	}

	switch j.engine.model.Lang {
	case ScriptLangTs:
		var result goja.Value
		result, err = j.tsCompile()
		if err != nil {
			return
		}

		j.engine.model.Compiled = result.String()

	case ScriptLangCoffee:
		var result goja.Value
		result, err = j.coffeeCompile()
		if err != nil {
			return
		}

		j.engine.model.Compiled = result.String()

	case ScriptLangJavascript:
		j.engine.model.Compiled = j.engine.model.Source

	}

	j.program, err = goja.Compile("", j.engine.model.Compiled, false)

	return
}

func (j *Javascript) GetCompiler() error {

	switch j.engine.model.Lang {
	case ScriptLangTs:
		data, err := Asset("scripts/typescriptServices.js")
		if err != nil {
			log.Error(err.Error())
			return err
		}

		j.compiler = string(data)

	case ScriptLangCoffee:
		data, err := Asset("scripts/coffee-script.js")
		if err != nil {
			log.Error(err.Error())
			return err
		}

		j.compiler = string(data)

	default:

	}

	return nil
}

func (j *Javascript) tsCompile() (result goja.Value, err error) {

	if _, err = j.EvalString(j.compiler); err != nil {
		return
	}

	const options = `{ target: ts.ScriptTarget.ES5, newLine: 1 }`

	// prepare script to inline
	doc := strings.Join(strings.Split(j.engine.model.Source, "\n"), `\n`)
	doc = strings.Replace(doc, `"`, `\"`, -1)

	var SRC = fmt.Sprintf(`ts.transpile("%s", %s);`, doc, options)

	// compile from typescript to native script
	var program *goja.Program
	if program, err = goja.Compile("", SRC, false); err != nil {
		return
	}

	result, err = j.vm.RunProgram(program)

	return
}

func (j *Javascript) coffeeCompile() (result goja.Value, err error) {

	if _, err = j.EvalString(j.compiler); err != nil {
		return
	}

	// prepare script to inline
	doc := strings.Join(strings.Split(j.engine.model.Source, "\n"), `\n`)
	doc = strings.Replace(doc, `"`, `\"`, -1)

	var SRC = fmt.Sprintf(`CoffeeScript.compile("%s", {"bare":true})`, doc)

	// compile from coffee to native script
	var program *goja.Program
	if program, err = goja.Compile("", SRC, false); err != nil {
		return
	}

	result, err = j.vm.RunProgram(program)

	return
}

func (j *Javascript) Do() (result string, err error) {
	var value goja.Value
	if value, err = j.vm.RunProgram(j.program); err != nil {
		return
	}
	result = value.String()
	return
}

func (j *Javascript) AssertFunction(f string) (result string, err error) {
	if assertFunc, ok := goja.AssertFunction(j.vm.Get(f)); ok {
		var value goja.Value
		if value, err = assertFunc(goja.Undefined(), j.vm.ToValue(4), j.vm.ToValue(10)); err != nil {
			return
		}
		result = value.String()
	}
	return
}

func (j *Javascript) PushStruct(name string, s interface{}) {
	j.vm.Set(name, s)
}

func (j *Javascript) PushFunction(name string, s interface{}) {
	j.vm.Set(name, s)
}

func (j *Javascript) EvalString(src string) (result string, err error) {

	var program *goja.Program
	if program, err = goja.Compile("", src, false); err != nil {
		return
	}

	var value goja.Value
	if value, err = j.vm.RunProgram(program); err != nil {
		return
	}

	result = value.String()

	return
}

func (j *Javascript) bind() {

	//
	// print()
	// hex2arr()
	// CurrentNode()
	// CurrentDevice()
	//

	j.vm.Set("print", fmt.Println)

	_, _ = j.vm.RunString(`
	
	var self = {},
    console = {log:print,warn:print,error:print,info:print},
    global = {};
	
	hex2arr = function (hexString) {
	   var result = [];
	   while (hexString.length >= 2) {
		   result.push(parseInt(hexString.substring(0, 2), 16));
		   hexString = hexString.substring(2, hexString.length);
	   }
	   return result;
	};

	CurrentNode = function(){

		var action, flow, node;
		node = null;

		if (typeof Flow !== "undefined" && Flow !== null) {
			node = Flow.Node();
		}
		if (!node && (typeof Action !== "undefined" && Action !== null)) {
			node = Action.Node();
		}

		if (!node) {
			//warn('node not found');
			return null;
		}

		return node;
	};

	CurrentDevice = function(){

		var action, dev;
		dev = null;

		if (!dev && (typeof Action !== "undefined" && Action !== null)) {
			dev = Action.Device();
		}

		if (!dev) {
			//warn('device not found');
			return null;
		}

		return dev;
	};

	`)

	for name, structure := range j.engine.structures.heap {
		j.vm.Set(name, structure)
	}

	for name, structure := range j.engine.functions.heap {
		j.vm.Set(name, structure)
	}

	return
}

func (j *Javascript) CreateProgram(name, source string) (err error) {
	j.lockPrograms.Lock()
	j.programs[name], err = goja.Compile("", source, false)
	j.lockPrograms.Unlock()
	return
}

func (j *Javascript) RunProgram(name string) (result string, err error) {
	j.lockPrograms.Lock()
	defer j.lockPrograms.Unlock()

	program, ok := j.programs[name]
	if !ok {
		err = ErrorProgramNotFound
		return
	}
	var value goja.Value
	if value, err = j.vm.RunProgram(program); err != nil {
		return
	}

	result = value.String()

	return
}
