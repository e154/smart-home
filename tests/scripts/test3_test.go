package scripts

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/common"
	"fmt"
)

type Test3Message struct {
	Value     string `json:"value"`
	Error     string `json:"error"`
	Direction bool   `json:"direction"`
}

func Test3(t *testing.T) {

	message := &Test3Message{}
	s := ""

	Convey("input <==> output", t, func(ctx C) {
		container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			storeRegisterCallback(scriptService)

			script1 := &m.Script{
				Source: s,
				Lang:   common.ScriptLangCoffee,
			}
			engine, err := scriptService.NewEngine(script1)
			engine.PushStruct("message", message)
			So(err, ShouldBeNil)
			err = engine.Compile()
			So(err, ShouldBeNil)
			_, err = engine.Do()
			So(err, ShouldBeNil)

			script3 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test3",
				Source:      coffeeScript3,
				Description: "test3",
			}

			engine3, err := scriptService.NewEngine(script3)
			So(err, ShouldBeNil)
			err = engine3.Compile()
			So(err, ShouldBeNil)

			err = engine.EvalString(script3.Compiled)
			So(err, ShouldBeNil)

			err = engine.EvalString(`IC.store(message);`)
			So(err, ShouldBeNil)

			value, ok := store.(map[string]interface{})
			So(ok, ShouldEqual, true)

			fmt.Println(value["error"].(string))
			fmt.Println(store)

			fmt.Println(message)
		})
	})
}
