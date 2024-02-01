// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"embed"
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"runtime/debug"

	"github.com/dop251/goja"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/system/scripts/eventloop"
)

//go:embed typescript.js
//go:embed coffeescript.js
var scriptsAsset embed.FS

// Javascript ...
type Javascript struct {
	engine       *Engine
	compiler     string
	vm           *goja.Runtime
	loop         *eventloop.EventLoop
	program      *goja.Program
	lockPrograms sync.Mutex
	programs     map[string]*goja.Program
}

// NewJavascript ...
func NewJavascript(engine *Engine) *Javascript {
	return &Javascript{
		engine: engine,

		programs: make(map[string]*goja.Program),
	}
}

// Init ...
func (j *Javascript) Init() (err error) {

	j.vm = goja.New()
	j.vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	//j.vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	j.loop = eventloop.NewEventLoop(j.vm)

	j.bind()

	if j.engine.model.Compiled == "" {
		return
	}

	if j.program, err = goja.Compile("", j.engine.model.Compiled, false); err != nil {
		log.Error(err.Error())
	}

	return
}

// Compile ...
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

// GetCompiler ...
func (j *Javascript) GetCompiler() error {

	switch j.engine.model.Lang {
	case ScriptLangTs:
		data, err := scriptsAsset.ReadFile("typescript.js")
		if err != nil {
			log.Error(err.Error())
		}
		j.compiler = string(data)

	case ScriptLangCoffee:
		data, err := scriptsAsset.ReadFile("coffeescript.js")
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

	// prepare script to inline
	doc := strings.Join(strings.Split(j.engine.model.Source, "\n"), `\n`)
	doc = strings.Replace(doc, `"`, `\"`, -1)

	var SRC = fmt.Sprintf(`
function compileTypeScriptString(tsCode) {
  const compilerOptions = {
    target: 'es6',
    newLine: 1,
    module: ts.ModuleKind.CommonJS,
  };

  const result = ts.transpileModule(tsCode, { compilerOptions });
  return result.outputText;
}

try {
  compileTypeScriptString("%s");
} catch (error) {
  console.error(error.message);
}
`, doc)

	// compile from typescript to native script
	var program *goja.Program
	if program, err = goja.Compile("", SRC, true); err != nil {
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
	if program, err = goja.Compile("", SRC, true); err != nil {
		return
	}

	result, err = j.vm.RunProgram(program)

	return
}

// Do ...
func (j *Javascript) Do() (result string, err error) {
	result, err = j.unsafeRun(j.program)
	if err != nil {
		if j.engine.ScriptId() != 0 {
			err = errors.Wrapf(err, "script id:%d ", j.engine.ScriptId())
		}
	}
	return
}

// AssertFunction ...
func (j *Javascript) AssertFunction(f string, args ...interface{}) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Recovered script id: %d, %s", j.engine.model.Id, f)
			log.Debug(j.vm.Get(f).String())
			debug.PrintStack()
		}
	}()
	if assertFunc, ok := goja.AssertFunction(j.vm.Get(f)); ok {
		var value goja.Value
		var gojaArgs []goja.Value
		for _, arg := range args {
			gojaArgs = append(gojaArgs, j.vm.ToValue(arg))
		}
		if value, err = assertFunc(goja.Undefined(), gojaArgs...); err != nil {
			return
		}
		result = value.String()
	}
	return
}

// PushStruct ...
func (j *Javascript) PushStruct(name string, s interface{}) {
	_ = j.vm.Set(name, s)
}

// PushFunction ...
func (j *Javascript) PushFunction(name string, s interface{}) {
	_ = j.vm.Set(name, s)
}

// EvalString ...
func (j *Javascript) EvalString(src string) (result string, err error) {

	var program *goja.Program
	if program, err = goja.Compile("", src, false); err != nil {
		return
	}

	result, err = j.unsafeRun(program)

	return
}

func (j *Javascript) bind() {

	//
	// print()
	// console()
	// hex2arr()
	// marshal(obj)
	// unmarshal(json)
	//

	_ = j.vm.Set("print", log.Info)

	_, _ = j.vm.RunString(`

    console = {log:print,warn:print,error:print,info:print},
	hex2arr = function (hexString) {
	   var result = [];
	   while (hexString.length >= 2) {
		   result.push(parseInt(hexString.substring(0, 2), 16));
		   hexString = hexString.substring(2, hexString.length);
	   }
	   return result;
	};
	unmarshal = function(j) { return JSON.parse(j); }
	marshal = function(obj) { return JSON.stringify(obj); }
	`)

	j.engine.functions.Range(func(key, value interface{}) bool {
		_ = j.vm.Set(key.(string), value)
		return true
	})

	j.engine.structures.Range(func(key, value interface{}) bool {
		_ = j.vm.Set(key.(string), value)
		return true
	})
}

// CreateProgram ...
func (j *Javascript) CreateProgram(name, source string) (err error) {
	j.lockPrograms.Lock()
	j.programs[name], err = goja.Compile("", source, false)
	j.lockPrograms.Unlock()
	return
}

// RunProgram ...
func (j *Javascript) RunProgram(name string) (result string, err error) {
	j.lockPrograms.Lock()
	defer j.lockPrograms.Unlock()

	program, ok := j.programs[name]
	if !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("name \"%s\"", name))
		return
	}

	result, err = j.unsafeRun(program)

	return
}

func (j *Javascript) unsafeRun(program *goja.Program) (result string, err error) {

	if program == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Warn("Recovered script: ", j.engine.model.Id)
			debug.PrintStack()
		}
	}()

	var value goja.Value

	wg := sync.WaitGroup{}
	wg.Add(1)

	j.loop.Run(func(vm *goja.Runtime) {
		defer wg.Done()
		value, err = vm.RunProgram(program)
	})

	wg.Wait()

	if err != nil {
		err = errors.Wrap(err, "unsafeRun")
		return
	}

	result = value.String()

	return
}
