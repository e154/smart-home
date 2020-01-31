package scripts

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	"github.com/jinzhu/copier"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test11(t *testing.T) {

	type Foo struct {
		Bar string
		Foo *Foo
	}

	counter := 0

	pool := []string{
		"bar",
		"&{foo <nil>}",
		"foo",
		"<nil>",
		"bar_new",
		"&{foo_new <nil>}",
		"foo_new",
		"<nil>",
	}

	initCallback := func(ctx C) {
		store = func(i interface{}) {
			v := fmt.Sprintf("%v", i)
			//fmt.Println("v:", v)

			if counter >= len(pool) {
				fmt.Println("========= WARNING =========")
				fmt.Printf("counter(%d), v: %v\n", counter, v)
				return
			}

			switch counter {
			default:
				ctx.So(v, ShouldEqual, pool[counter])
			}

			counter++
		}
	}

	Convey("require external library", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			initCallback(ctx)

			storeRegisterCallback(scriptService)

			script1 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test11",
				Source:      coffeeScripts["coffeeScript26"],
				Description: "test11",
			}

			bar2 := &Foo{
				Bar: "bar",
				Foo: &Foo{
					Bar: "foo",
				},
			}

			engine, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)


			engine.PushGlobalProxy("bar2", bar2)

			err = engine.Compile()
		 	So(err, ShouldBeNil)

			_, err = engine.Do()
			So(err, ShouldBeNil)

			newBar2 := &Foo{
				Bar: "bar_new",
				Foo: &Foo{
					Bar: "foo_new",
				},
			}

			err = copier.Copy(&bar2, &newBar2)
			So(err, ShouldBeNil)

			_, err = engine.Do()
			So(err, ShouldBeNil)

			time.Sleep(time.Second * 2)
		})
	})
}
