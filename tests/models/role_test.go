package models

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"strings"
)

func TestRole(t *testing.T) {

	Convey("add role", t, func(ctx C) {
		container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			accessList *access_list.AccessListService) {

			// clear database
			migrations.Purge()

			demoRole := &m.Role{
				Name: "demo",
			}
			err := adaptors.Role.Add(demoRole)
			So(err, ShouldBeNil)

			userRole := &m.Role{
				Name:   "user",
				Parent: demoRole,
			}
			err = adaptors.Role.Add(userRole)
			So(err, ShouldBeNil)

			//debug.Println(accessList.List)
			//fmt.Println("----")

			var counter int
			for pack, item := range *accessList.List {
				for right := range item {
					if strings.Contains(right, "read") ||
						strings.Contains(right, "view") {
						permission := &m.Permission{
							RoleName:    demoRole.Name,
							PackageName: pack,
							LevelName:   right,
						}

						counter++
						_, err = adaptors.Permission.Add(permission)
						So(err, ShouldBeNil)
					}
				}
			}

			for pack, item := range *accessList.List {
				for right := range item {
					if !strings.Contains(right, "read") &&
						!strings.Contains(right, "view") {
						permission := &m.Permission{
							RoleName:    userRole.Name,
							PackageName: pack,
							LevelName:   right,
						}

						counter++
						_, err = adaptors.Permission.Add(permission)
						So(err, ShouldBeNil)
					}
				}
			}

			adminRole := &m.Role{
				Name:   "admin",
				Parent: userRole,
			}
			err = adaptors.Role.Add(adminRole)
			So(err, ShouldBeNil)

			// user
			userRole, err = adaptors.Role.GetByName("user")
			So(err, ShouldBeNil)
			So(userRole.Name, ShouldEqual, "user")
			So(userRole.Parent, ShouldNotBeNil)
			So(userRole.Parent.Name, ShouldEqual, "demo")
			So(len(userRole.Children), ShouldEqual, 1)
			So(userRole.Children[0].Name, ShouldEqual, "admin")

			// demo
			demoRole, err = adaptors.Role.GetByName("demo")
			So(err, ShouldBeNil)
			So(demoRole.Parent, ShouldBeNil)
			So(demoRole.Name, ShouldEqual, "demo")
			So(len(demoRole.Children), ShouldEqual, 1)
			So(demoRole.Children[0].Name, ShouldEqual, "user")

			// admin
			adminRole, err = adaptors.Role.GetByName("admin")
			So(err, ShouldBeNil)
			So(adminRole.Parent, ShouldNotBeNil)
			So(adminRole.Parent.Name, ShouldEqual, "user")
			So(adminRole.Name, ShouldEqual, "admin")
			So(len(adminRole.Children), ShouldBeZeroValue)

			permissions, err := adaptors.Permission.GetAllPermissions("user")
			So(err, ShouldBeNil)

			So(len(permissions), ShouldEqual, counter)

			err = adaptors.Role.Delete("demo")
			So(err, ShouldNotBeNil)
			err = adaptors.Role.Delete("user")
			So(err, ShouldNotBeNil)
			err = adaptors.Role.Delete("admin")
			So(err, ShouldBeNil)
			err = adaptors.Role.Delete("user")
			So(err, ShouldBeNil)
			err = adaptors.Role.Delete("demo")
			So(err, ShouldBeNil)
		})
	})
}