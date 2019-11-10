package models

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/migrations"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {

	Convey("add user", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			accessList *access_list.AccessListService) {

			// clear database
			migrations.Purge()

			// add role
			userRole := &m.Role{
				Name: "user_role",
			}
			err := adaptors.Role.Add(userRole)
			So(err, ShouldBeNil)

			var counter int
			for pack, item := range *accessList.List {
				for right := range item {
					if strings.Contains(right, "read") ||
						strings.Contains(right, "view") ||
						strings.Contains(right, "show") {
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

			const (
				nickname = "user"
				email    = "email@mail.com"
			)

			// add user
			user := &m.User{
				Nickname: nickname,
				RoleName: "user_role",
				Email:    email,
				Lang:     "en",
				Meta: []*m.UserMeta{
					&m.UserMeta{
						Key:   "phone1",
						Value: "+18004001234",
					},
				},
			}
			err = user.SetPass("123456")
			So(err, ShouldBeNil)

			ok, _ := user.Valid()
			So(ok, ShouldEqual, true)

			user.Id, err = adaptors.User.Add(user)
			So(err, ShouldBeNil)

			user, err = adaptors.User.GetById(user.Id)
			So(err, ShouldBeNil)

			//debug.Println(user)
			//fmt.Println("----")

			So(user.Nickname, ShouldEqual, nickname)
			So(user.Email, ShouldEqual, email)
			So(user.Role, ShouldNotBeNil)
			So(user.Role.Name, ShouldEqual, "user_role")
		})
	})
}
