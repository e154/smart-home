package use_case

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"errors"
	"github.com/e154/smart-home/system/access_list"
)

func AddRole(role *m.Role, adaptors *adaptors.Adaptors) (ok bool, errs []*validation.Error, err error) {

	// validation
	ok, errs = role.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.Role.Add(role); err != nil {
		return
	}

	return
}

func GetRoleByName(name string, adaptors *adaptors.Adaptors) (role *m.Role, err error) {

	role, err = adaptors.Role.GetByName(name)

	return
}

func UpdateRole(role *m.Role, adaptors *adaptors.Adaptors) (ok bool, errs []*validation.Error, err error) {

	// validation
	ok, errs = role.Valid()
	if len(errs) > 0 {
		return
	}

	err = adaptors.Role.Update(role)

	return
}

func GetRoleList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (items []*m.Role, total int64, err error) {

	items, total, err = adaptors.Role.List(limit, offset, order, sortBy)

	return
}

func DeleteRoleByName(name string, adaptors *adaptors.Adaptors) (err error) {

	if name == "admin" {
		err = errors.New("admin is base role")
		return
	}
	err = adaptors.Role.Delete(name)

	return
}

func Search(query string, limit, offset int, adaptors *adaptors.Adaptors) (roles []*m.Role, total int64, err error) {

	roles, total, err = adaptors.Role.Search(query, limit, offset)

	return
}

func GetAccessList(name string,
	adaptors *adaptors.Adaptors,
	accessListService *access_list.AccessListService) (accessList access_list.AccessList, err error) {

	var role *m.Role
	if role, err = adaptors.Role.GetByName(name); err != nil {
		return
	}

	var permissions []*m.Permission
	if permissions, err = adaptors.Permission.GetAllPermissions(role.Name); err != nil {
		return
	}

	accessList = make(access_list.AccessList)
	var item access_list.AccessItem
	var levels access_list.AccessLevels
	var ok bool
	list := *accessListService.List
	for _, perm := range permissions {

		if levels, ok = list[perm.PackageName]; !ok {
			continue
		}

		if accessList[perm.PackageName] == nil {
			accessList[perm.PackageName] = access_list.NewAccessLevels()
		}

		if item, ok = levels[perm.LevelName]; !ok {
			continue
		}

		item.RoleName = perm.RoleName
		accessList[perm.PackageName][perm.LevelName] = item
	}

	return
}
