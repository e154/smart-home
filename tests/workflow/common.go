package workflow

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
)

var store = func(interface{}) {}
var store2 = func(interface{}) {}

func storeRegisterCallback(scriptService *scripts.ScriptService) {
	scriptService.PushFunctions("store", func(value interface{}) {
		store(value)
	})
	scriptService.PushFunctions("store2", func(value interface{}) {
		store2(value)
	})
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
