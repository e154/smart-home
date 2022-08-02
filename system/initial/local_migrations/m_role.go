package local_migrations

import (
	"context"
	"strings"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/validation"
)

type MigrationRoles struct {
	adaptors   *adaptors.Adaptors
	accessList access_list.AccessListService
	validation *validation.Validate
}

func NewMigrationRoles(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	validation *validation.Validate) *MigrationRoles {
	return &MigrationRoles{
		adaptors:   adaptors,
		accessList: accessList,
		validation: validation,
	}
}

func (r *MigrationRoles) Up(ctx context.Context, adaptors *adaptors.Adaptors) (err error) {
	if adaptors != nil {
		r.adaptors = adaptors
	}

	if _, err = r.addAdmin(); err != nil {
		return
	}
	var demo *m.Role
	if demo, err = r.addDemo(); err != nil {
		return
	}
	_, err = r.addUser(demo)

	return
}

func (r *MigrationRoles) addAdmin() (adminRole *m.Role, err error) {

	if adminRole, err = r.adaptors.Role.GetByName("admin"); err != nil {
		adminRole = &m.Role{
			Name: "admin",
		}
		err = r.adaptors.Role.Add(adminRole)
		So(err, ShouldBeNil)
	}
	if _, err = r.adaptors.User.GetByNickname("admin"); err == nil {
		return
	}

	// add admin
	adminUser := &m.User{
		Nickname: "admin",
		RoleName: adminRole.Name,
		Email:    "admin@e154.ru",
		Lang:     "en",
		Status:   "active",
	}
	err = adminUser.SetPass("admin")
	So(err, ShouldBeNil)

	for pack, item := range *r.accessList.List() {
		for right := range item {
			permission := &m.Permission{
				RoleName:    adminUser.Nickname,
				PackageName: pack,
				LevelName:   right,
			}

			_, err = r.adaptors.Permission.Add(permission)
			So(err, ShouldBeNil)
		}
	}

	ok, _ := r.validation.Valid(adminUser)
	So(ok, ShouldEqual, true)

	adminUser.Id, err = r.adaptors.User.Add(adminUser)
	So(err, ShouldBeNil)

	return
}

func (r *MigrationRoles) addUser(demoRole *m.Role) (userRole *m.Role, err error) {

	if userRole, err = r.adaptors.Role.GetByName("user"); err != nil {
		userRole = &m.Role{
			Name:   "user",
			Parent: demoRole,
		}
		err := r.adaptors.Role.Add(userRole)
		So(err, ShouldBeNil)

		for pack, item := range *r.accessList.List() {
			for right := range item {
				if strings.Contains(right, "create") ||
					strings.Contains(right, "update") ||
					strings.Contains(right, "delete") {
					permission := &m.Permission{
						RoleName:    userRole.Name,
						PackageName: pack,
						LevelName:   right,
					}

					_, err = r.adaptors.Permission.Add(permission)
					So(err, ShouldBeNil)
				}
			}
		}
	}

	if _, err = r.adaptors.User.GetByNickname("user"); err == nil {
		return
	}

	// add base user
	baseUser := &m.User{
		Nickname: "user",
		RoleName: userRole.Name,
		Email:    "user@e154.ru",
		Lang:     "en",
		Status:   "active",
	}
	err = baseUser.SetPass("user")
	So(err, ShouldBeNil)

	ok, _ := r.validation.Valid(baseUser)
	So(ok, ShouldEqual, true)

	baseUser.Id, err = r.adaptors.User.Add(baseUser)
	So(err, ShouldBeNil)

	return
}

func (r *MigrationRoles) addDemo() (demoRole *m.Role, err error) {

	if demoRole, err = r.adaptors.Role.GetByName("demo"); err != nil {

		demoRole = &m.Role{
			Name: "demo",
		}
		err = r.adaptors.Role.Add(demoRole)
		So(err, ShouldBeNil)

		for pack, item := range *r.accessList.List() {
			for right := range item {
				if strings.Contains(right, "read") ||
					strings.Contains(right, "view") ||
					strings.Contains(right, "preview") {
					permission := &m.Permission{
						RoleName:    demoRole.Name,
						PackageName: pack,
						LevelName:   right,
					}

					_, err = r.adaptors.Permission.Add(permission)
					So(err, ShouldBeNil)
				}
			}
		}
	}

	if _, err = r.adaptors.User.GetByNickname("demo"); err == nil {
		return
	}

	// add demo user
	demoUser := &m.User{
		Nickname: "demo",
		RoleName: demoRole.Name,
		Email:    "demo@e154.ru",
		Lang:     "en",
		Status:   "active",
	}
	err = demoUser.SetPass("demo")
	So(err, ShouldBeNil)

	ok, _ := r.validation.Valid(demoUser)
	So(ok, ShouldEqual, true)

	demoUser.Id, err = r.adaptors.User.Add(demoUser)
	So(err, ShouldBeNil)

	return
}
