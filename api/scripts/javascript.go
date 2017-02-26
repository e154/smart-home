package scripts

import (
	"fmt"
	"github.com/e154/go-candyjs"
	"strings"
	"errors"
	"github.com/e154/smart-home/api/models"
	r "github.com/e154/smart-home/lib/rpc"
	"github.com/e154/smart-home/lib/common"
)

const (
	coffescript = "coffee-script.js"
)


type Javascript struct {
	engine	*Engine
	ctx	*candyjs.Context
}

func (j *Javascript) Init() (err error) {

	j.ctx = candyjs.NewContext()
	if err = j.ctx.PevalString(j.engine.model.Compiled); err != nil {
		return
	}

	j.pushGlobalCandyJSObject()

	return
}

func (j *Javascript) pushGlobalCandyJSObject() {

	// print
	j.ctx.PushGlobalGoFunction("print", func(a ...interface{}){
		j.engine.Print(a...)
	})

	// etc
	j.PushStruct("request", &r.Request{})
	j.PushStruct("device", &models.Device{})
	j.PushStruct("flow", &models.Flow{})
	j.PushStruct("node", &models.Node{})
	j.ctx.PevalString(`SmartJs = {}`)
	j.ctx.PevalString(`SmartJs.hex2arr = function (hexString) {
   var result = [];
   while (hexString.length >= 2) {
       result.push(parseInt(hexString.substring(0, 2), 16));
       hexString = hexString.substring(2, hexString.length);
   }
   return result;
}`)

	j.PushFunction("node_send", func(protocol string, node *models.Node, args *r.Request,) (result r.Result) {
		if args == nil {
			result.Error = "args is nil pointer"
			return
		}

		if node == nil {
			result.Error = "node is nil pointer"
			return
		}

		result = node.Send(protocol, args)

		return
	})
}

func (j *Javascript) Close() {
	j.ctx.Pop()
	j.ctx.DestroyHeap()
}

func (j *Javascript) Compile() (err error) {

	switch j.engine.model.Lang {
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

func (j *Javascript) coffeeCompile() (result string, err error) {

	//TODO cache file
	file := common.GetScript(coffescript)
	if err = j.ctx.PevalString(string(file)); err != nil {
	//if err = j.ctx.PevalFile(coffescript); err != nil {
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