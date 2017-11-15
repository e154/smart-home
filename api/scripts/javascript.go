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
	return j.DoCustom("main")
}

func (j *Javascript) DoCustom(f string) (result string, err error) {

	//j.ctx.PushGlobalObject()
	//if err = j.ctx.PushTimers(); err != nil {
	//	return
	//}

	if b := j.ctx.GetGlobalString(f); !b {
		err = errors.New("main function not found!")
		return
	}

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

	if res := j.ctx.SafeToString(-1); res != "undefined" {
		result = res
	}

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

	// base
	j.ctx.PevalString(fmt.Sprintf(`
	IC = {};
	IC.Runmode = '%s';
	IC.hex2arr = function (hexString) {
	   var result = [];
	   while (hexString.length >= 2) {
		   result.push(parseInt(hexString.substring(0, 2), 16));
		   hexString = hexString.substring(2, hexString.length);
	   }
	   return result;
	};

	IC.CurrentNode = function(){

		var action, flow, node;
		node = null;

		if ((typeof IC !== "undefined" && IC !== null ? IC.Flow : void 0) != null) {
			flow = IC.Flow();
			node = flow.node();
		}
		if (!node && ((typeof IC !== "undefined" && IC !== null ? IC.Action : void 0) != null)) {
			action = IC.Action();
			node = action.node();
		}

		if (!node) {
			//IC.warn('node not found');
			return null;
		}

		return node;
	};

	IC.CurrentDevice = function(){

		var action, dev;
		dev = null;

		if (!dev && ((typeof IC !== "undefined" && IC !== null ? IC.Action : void 0) != null)) {
			action = IC.Action();
			dev = action.device();
		}

		if (!dev) {
			//IC.warn('device not found');
			return null;
		}

		return dev;
	};

	`, beego.BConfig.RunMode))

	// push structures
	for name, structure := range pull.GetStruct() {
		if b := j.ctx.GetGlobalString("IC"); !b {
			return
		}
		j.ctx.PushObject()
		j.ctx.PushStruct(structure)
		j.ctx.PutPropString(-3, name)
		j.ctx.Pop()
	}

	// push functions
	for name, function := range pull.Getfunctions() {
		if b := j.ctx.GetGlobalString("IC"); !b {
			return
		}
		j.ctx.PushObject()
		j.ctx.PushGoFunction(function)
		j.ctx.PutPropString(-3, name)
		j.ctx.Pop()
	}

	// print
	j.PushFunction("print", func(a ...interface{}){
		j.engine.Print(a...)
	})

	return
}