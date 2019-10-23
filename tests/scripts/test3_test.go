package scripts

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test3(t *testing.T) {

	var state string
	store = func(i interface{}) {
		state = fmt.Sprintf("%v", i)
	}

	Convey("eval script", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			storeRegisterCallback(scriptService)

			// engine
			// ------------------------------------------------
			script := &m.Script{
				Lang: common.ScriptLangJavascript,
			}
			engine, err := scriptService.NewEngine(script)
			So(err, ShouldBeNil)

			// scripts
			// ------------------------------------------------
			scripts := GetScripts(ctx, scriptService, adaptors, 3)

			// execute script
			// ------------------------------------------------
			err = engine.EvalString(scripts["script3"].Compiled)
			So(err, ShouldBeNil)

			_, err = engine.DoCustom("on_enter")
			So(err, ShouldBeNil)
			So(state, ShouldEqual, "on_enter")

			_, err = engine.DoCustom("on_exit")
			So(err, ShouldBeNil)
			So(state, ShouldEqual, "on_exit")
		})
	})
}
