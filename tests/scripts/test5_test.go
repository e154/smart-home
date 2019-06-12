package scripts

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"fmt"
)

func Test5(t *testing.T) {

	s := ""

	Convey("javascript PushStruct", t, func(ctx C) {
		container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			storeRegisterCallback(scriptService)

			script1 := &m.Script{
				Source: s,
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

			So(fmt.Sprintf("%v", store), ShouldEqual, fmt.Sprintf("[42 84 21 63]"))
		})
	})
}
