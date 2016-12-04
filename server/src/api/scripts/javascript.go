package scripts

import (
	"fmt"
	"gopkg.in/olebedev/go-duktape.v2"
	"strings"
	"errors"
)

const (
	coffescript = "coffee-script.js"
)


type Javascript struct {
	engine	*Engine
	ctx	*duktape.Context
}

func (j *Javascript) Init() (err error) {

	j.ctx = duktape.New()
	if err = j.ctx.PevalString(j.engine.model.Compiled); err != nil {
		return
	}

	j.Reg()

	return
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

	if err = j.ctx.PevalFile(coffescript); err != nil {
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
		err = errors.New("main not found")
		return
	}

	// call(arg)
	// arg = stack - num args
	if r := j.ctx.Pcall(0); r != 0 {
		err = errors.New(j.ctx.SafeToString(-1))
		return
	}

	result = j.ctx.SafeToString(-1)

	j.ctx.Pop()

	return
}

func (j *Javascript) Reg() (err error) {

	p := func(ctx *duktape.Context) int {
		j.engine.Print(ctx.SafeToString(-1))
		return 1
	}

	// override default print function
	if _, err = j.ctx.PushGlobalGoFunction("print", p); err != nil {
		return
	}

	return
}