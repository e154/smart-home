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

func Test4(t *testing.T) {

	var state string
	store = func(i interface{}) {
		state = fmt.Sprintf("%v", i)
	}

	Convey("javascript PushStruct", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migratxions *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			storeRegisterCallback(scriptService)

			script1 := &m.Script{
				Lang:   common.ScriptLangCoffee,
			}
			engine, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)

			for i:=0;i<2;i++ {
				_, err = engine.PushStruct("test", &MyStruct{
					Int:     42,
					Float64: 21.0,
					Nested:  &MyStruct{Int: 21},
				})
				So(err, ShouldBeNil)
			}

			err = engine.EvalString(`IC.store([
				test.int,
				test.multiply(2),
				test.nested.int,
				test.nested.multiply(3)
			])`)
			So(err, ShouldBeNil)

			So(state, ShouldEqual, fmt.Sprintf("[42 84 21 63]"))
		})
	})
}
