package scripts

import (
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
)

func Test1(t *testing.T) {

	var state string
	store = func(i interface{}) {
		state = fmt.Sprintf("%v", i)
	}

	var script1 *m.Script
	Convey("scripts run syn command", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			storeRegisterCallback(scriptService)

			script1 = &m.Script{
				Lang:        "coffeescript",
				Name:        "test1",
				Source:      coffeeScript1,
				Description: "test1",
			}

			engine1, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)
			err = engine1.Compile()
			So(err, ShouldBeNil)

			_, err = engine1.Do()
			So(err, ShouldBeNil)

			So(state, ShouldEqual, "ok")
		})
	})
}
