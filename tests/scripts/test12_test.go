// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/storage"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test12(t *testing.T) {

	counter := 0

	pool := []string{
		"",
		`{"bar":"bar"}`,
		"bar",
		`{"bar":"foo"}`,
		"foo",
		`map[]`,
		`map[foo:{"bar":"foo"}]`,
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

	Convey("check db storage", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			storageService *storage.Storage,
			scriptService scripts.ScriptService) {

			// clear database
			// ------------------------------------------------
			migrations.Purge()

			initCallback(ctx)

			storeRegisterCallback(scriptService)

			script1 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test12",
				Source:      coffeeScripts["coffeeScript27"],
				Description: "test12",
			}

			engine, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)

			err = engine.Compile()
			So(err, ShouldBeNil)

			_, err = engine.Do()
			So(err, ShouldBeNil)

			_, err = adaptors.Storage.GetByName("foo")
			So(err, ShouldNotBeNil)

			storageService.Serialize()
			storage, err := adaptors.Storage.GetByName("foo")
			So(err, ShouldBeNil)
			So(string(storage.Value), ShouldEqual, `{"bar": "foo"}`)

			err = adaptors.Storage.CreateOrUpdate(m.Storage{
				Name: "foo2",
				Value: []byte(`{"foo":"bar"}`),
			})
			So(err, ShouldBeNil)

			storage, err = adaptors.Storage.GetByName("foo2")
			So(err, ShouldBeNil)
			So(string(storage.Value), ShouldEqual, `{"foo": "bar"}`)
		})
	})
}
