// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package system

import (
	"context"
	"testing"

	"github.com/e154/smart-home/system/migrations"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/e154/smart-home/adaptors"
	localMigrations "github.com/e154/smart-home/system/initial/local_migrations"
)

func TestLocalMigration(t *testing.T) {

	t.Run("local migration", func(t *testing.T) {
		Convey("case 1", t, func(ctx C) {
			err := container.Invoke(func(adaptors *adaptors.Adaptors,
				localMigrations *localMigrations.Migrations,
				migrations *migrations.Migrations) {

				err := migrations.Purge()
				ctx.So(err, ShouldBeNil)

				oldVersion := ""
				currentVersion, err := localMigrations.Up(context.TODO(), nil, oldVersion)
				ctx.So(err, ShouldBeNil)

				//fmt.Println(currentVersion)
				ctx.So(currentVersion, ShouldEqual, "*local_migrations.MigrationZones")
			})
			ctx.So(err, ShouldBeNil)
		})

		Convey("case 2", t, func(ctx C) {
			err := container.Invoke(func(adaptors *adaptors.Adaptors,
				localMigrations *localMigrations.Migrations,
				migrations *migrations.Migrations) {

				oldVersion := "*local_migrations.MigrationZones"
				currentVersion, err := localMigrations.Up(context.TODO(), nil, oldVersion)
				ctx.So(err, ShouldBeNil)

				//fmt.Println(currentVersion)
				ctx.So(currentVersion, ShouldEqual, "*local_migrations.MigrationZones")
			})
			ctx.So(err, ShouldBeNil)
		})

	})
}
