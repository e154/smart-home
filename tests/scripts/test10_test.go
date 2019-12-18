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

func Test10(t *testing.T) {

	type Foo struct {
		Bar string
		Foo *Foo
	}

	var script1 *m.Script
	Convey("require external library", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			scriptService.PushFunctions("So", func(actual interface{}, assert string, expected interface{}) {
				//fmt.Printf("actual(%v), expected(%v)\n", actual, expected)
				switch assert {
				case "ShouldEqual":
					So(fmt.Sprintf("%v", actual), ShouldEqual, expected)
				}

			})

			script1 = &m.Script{
				Lang:        "coffeescript",
				Name:        "test10",
				Source:      coffeeScript25,
				Description: "test10",
			}

			foo := &Foo{
				Bar: "foo",
			}
			bar := &Foo{
				Bar: "bar",
				Foo: foo,
			}

			scriptService.PushStruct("bar2", bar)

			scriptService.PushFunctions("external", func(varName string, f *Foo) {
				//fmt.Printf("varName: %v\n", varName)
				switch varName {
				case "bar":
					So(f.Foo.Bar, ShouldEqual, "")
				case "bar2":
					So(f.Foo.Bar, ShouldEqual, "foo")
				case "IC.bar2":
					So(f.Foo.Bar, ShouldEqual, "foo")
				}
			})

			engine, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)

			_, err = engine.PushStruct("bar", bar)
			So(err, ShouldBeNil)

			engine.PushGlobalProxy("bar2", bar)

			err = engine.Compile()
			So(err, ShouldBeNil)

			_, err = engine.Do()
			So(err, ShouldBeNil)
		})
	})
}
