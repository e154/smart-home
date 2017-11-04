package scripts

import (
	"fmt"
	"strings"
	"errors"
	"github.com/e154/go-candyjs"
	"github.com/astaxie/beego"
	"github.com/e154/smart-home/api/log"
)


type Javascript struct {
	engine			*Engine
	ctx				*candyjs.Context
	compiler		string
}

func (j *Javascript) Init() (err error) {

	j.ctx = candyjs.NewContext()
	if err = j.ctx.PevalString(j.engine.model.Compiled); err != nil {
		return
	}

	j.bind()

	return
}

func (j *Javascript) Close() {
	j.ctx.Pop()
	j.ctx.DestroyHeap()
}

func (j *Javascript) Compile() (err error) {

	j.GetCompiler()

	switch j.engine.model.Lang {
	case "ts":
		var result string
		result, err = j.tsCompile()
		if err != nil {
			return
		}

		j.engine.model.Compiled = result

	case "coffeescript":
		var result string
		result, err = j.coffeeCompile()
		if err != nil {
			return
		}

		j.engine.model.Compiled = result

	case "javascript":
		j.engine.model.Compiled = j.engine.model.Source

	}

	if err = j.ctx.PevalString(j.engine.model.Compiled); err != nil {
		return
	}

	return
}

func (j *Javascript) GetCompiler() error {

	switch j.engine.model.Lang {
	case "ts":
		data, err := Asset("scripts/typescriptServices.js")
		if err != nil {
			log.Error(err.Error())
			return err
		}

		j.compiler = string(data)

	case "coffeescript":
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

func (j *Javascript) tsCompile() (result string, err error) {

	if err = j.ctx.PevalString(j.compiler); err != nil {
		return
	}

	const options = `{ target: ts.ScriptTarget.ES5, newLine: 1 }`

	// prepare script to inline
	doc := strings.Join(strings.Split(j.engine.model.Source, "\n"), `\n`)
	doc = strings.Replace(doc, `"`, `\"`, -1)

	// compile from coffee to native script
	if err = j.ctx.PevalString(fmt.Sprintf(`ts.transpile("%s", %s);`, doc, options)); err != nil {
		return
	}

	// return native javascript code
	result = j.ctx.GetString(-1)

	return
}

func (j *Javascript) coffeeCompile() (result string, err error) {

	if err = j.ctx.PevalString( j.compiler ); err != nil {
		return
	}

	// prepare script to inline
	doc := strings.Join(strings.Split(j.engine.model.Source, "\n"), `\n`)
	doc = strings.Replace(doc, `"`, `\"`, -1)

	// compile from coffee to native script
	if err = j.ctx.PevalString(fmt.Sprintf(`CoffeeScript.compile("%s", {"bare":true})`, doc)); err != nil {
		return
	}

	// return native javascript code
	result = j.ctx.GetString(-1)

	return
}

func (j *Javascript) Do() (result string, err error) {

	j.ctx.PushGlobalObject()
	if b := j.ctx.GetPropString(-1, "main"); !b {
		err = errors.New("main function not found!")
		return
	}

	j.ctx.PushTimers()

	defer func() {
		if r := recover(); r != nil {
			log.Critical("Script: Recovered in f", r)
		}
	}()

	// call(arg)
	// arg = stack - num args
	if r := j.ctx.Pcall(0); r != 0 {
		err = errors.New(j.ctx.SafeToString(-1))
		return
	}

	result = j.ctx.SafeToString(-1)

	return
}

func (j *Javascript) PushStruct(name string, s interface{}) (int, error) {
	return j.ctx.PushGlobalStruct(name, s)
}

func (j *Javascript) PushFunction(name string, s interface{}) (int, error) {
	return j.ctx.PushGlobalGoFunction(name, s)
}

func (j *Javascript) EvalString(str string) error {
	return j.ctx.PevalString(str)
}

func (j *Javascript) Ctx() *candyjs.Context {
	return j.ctx
}

func (j *Javascript) bind() {

	// print
	j.PushFunction("print", func(a ...interface{}){
		j.engine.Print(a...)
	})

	j.EvalString(fmt.Sprintf(`run_mode = '%s'`, beego.BConfig.RunMode))
	j.PushStruct("log", &Log{})
	j.PushStruct("model", &Model{})
	j.PushStruct("node", &Node{})
	j.PushStruct("notifr", &Notifr{})
	j.PushStruct("exec", &Exec{})
	//j.PushFunction("to_time", func(i int64) time.Duration {
	//	return time.Duration(i)
	//})

	// etc
	j.EvalString(`helper = {}`)
	j.EvalString(`helper.hex2arr = function (hexString) {
   var result = [];
   while (hexString.length >= 2) {
       result.push(parseInt(hexString.substring(0, 2), 16));
       hexString = hexString.substring(2, hexString.length);
   }
   return result;
}`)

	return
}