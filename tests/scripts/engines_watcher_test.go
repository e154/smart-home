// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package scripts

import (
	"fmt"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/system/bus"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
)

func TestEnginesWatcher(t *testing.T) {

	t.Run("init", func(t *testing.T) {
		Convey("", t, func(ctx C) {
			err := container.Invoke(func(adaptors *adaptors.Adaptors,
				migrations *migrations.Migrations,
				scriptService scripts.ScriptService,
				eventBus bus.Bus) {

				script := &m.Script{
					Id:       1,
					Compiled: `function foo() {return 1 + 1}`,
				}

				script2 := &m.Script{
					Id:       2,
					Compiled: `function main() {return foo() + 1}`,
				}

				enginesWatcher, err := scriptService.NewEnginesWatcher([]*m.Script{script, script2})
				So(err, ShouldBeNil)
				So(enginesWatcher.Engine(), ShouldBeNil)

				enginesWatcher.Spawn(nil)

				So(enginesWatcher.Engine(), ShouldNotBeNil)
				result, err := enginesWatcher.Engine().AssertFunction("main")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, "3")

				// update script
				script.Compiled = `function foo() {return 1 + 2}`
				eventBus.Publish("system/models/scripts/1", events.EventUpdatedScriptModel{
					ScriptId: 1,
					Script:   script,
				})

				time.Sleep(time.Millisecond * 500)

				So(enginesWatcher.Engine(), ShouldNotBeNil)
				result, err = enginesWatcher.Engine().AssertFunction("main")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, "4")

				// remove script
				eventBus.Publish("system/models/scripts/1", events.EventRemovedScriptModel{
					ScriptId: 1,
					Script:   script,
				})

				time.Sleep(time.Millisecond * 500)

				So(enginesWatcher.Engine(), ShouldNotBeNil)
				result, err = enginesWatcher.Engine().AssertFunction("main")
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "ReferenceError: foo is not defined at main (<eval>:1:28(2))")
				So(result, ShouldEqual, "")

				enginesWatcher.Stop()
			})
			if err != nil {
				fmt.Println(err.Error())
			}
		})
	})

	t.Run("spawn", func(t *testing.T) {
		Convey("", t, func(ctx C) {
			err := container.Invoke(func(adaptors *adaptors.Adaptors,
				migrations *migrations.Migrations,
				scriptService scripts.ScriptService,
				eventBus bus.Bus) {

				script := &m.Script{
					Id:       1,
					Compiled: `function omega() {return foo() + bar + structObj.val}`,
				}

				script2 := &m.Script{
					Id:       2,
					Compiled: `function main() {return omega() + 1}`,
				}

				enginesWatcher, err := scriptService.NewEnginesWatcher([]*m.Script{script, script2})
				So(err, ShouldBeNil)
				So(enginesWatcher.Engine(), ShouldBeNil)

				enginesWatcher.PushFunction("foo", func() int {
					return 4
				})

				structObj := map[string]int{
					"val": 1,
				}
				enginesWatcher.PushStruct("structObj", structObj)

				enginesWatcher.BeforeSpawn(func(engine *scripts.Engine) {
					if _, err = engine.EvalString(fmt.Sprintf("const bar = %d;", 4)); err != nil {
					}
				})
				enginesWatcher.Spawn(func(engine *scripts.Engine) {

				})

				So(enginesWatcher.Engine(), ShouldNotBeNil)
				result, err := enginesWatcher.Engine().AssertFunction("main")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, "10")

				// update script
				script2.Compiled = `function main() {return omega() + 2}`
				eventBus.Publish("system/models/scripts/1", events.EventUpdatedScriptModel{
					ScriptId: 2,
					Script:   script2,
				})

				time.Sleep(time.Millisecond * 500)

				So(enginesWatcher.Engine(), ShouldNotBeNil)
				result, err = enginesWatcher.Engine().AssertFunction("main")
				So(err, ShouldBeNil)
				So(result, ShouldEqual, "11")

				// remove script
				eventBus.Publish("system/models/scripts/1", events.EventRemovedScriptModel{
					ScriptId: 1,
					Script:   script,
				})

				time.Sleep(time.Millisecond * 500)

				So(enginesWatcher.Engine(), ShouldNotBeNil)
				result, err = enginesWatcher.Engine().AssertFunction("main")
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "ReferenceError: omega is not defined at main (<eval>:1:30(2))")
				So(result, ShouldEqual, "")

				enginesWatcher.Stop()
			})
			if err != nil {
				fmt.Println(err.Error())
			}
		})
	})
}
