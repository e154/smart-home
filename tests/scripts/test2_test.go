package scripts

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test2(t *testing.T) {

	var state string
	store = func(i interface{}) {
		state = fmt.Sprintf("%v", i)
	}

	var script1 *m.Script
	Convey("require external library", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			storeRegisterCallback(scriptService)

			script1 = &m.Script{
				Lang:        "coffeescript",
				Name:        "test2",
				Source:      coffeeScript2,
				Description: "test2",
			}

			engine1, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)
			err = engine1.Compile()
			So(err, ShouldBeNil)

			_, err = engine1.Do()
			So(err, ShouldBeNil)

			So(state, ShouldEqual, "123-bar-Jan")
		})
	})
}
