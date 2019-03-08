package use_case

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"errors"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
)

func AddRole(roleParams models.NewRole, adaptors *adaptors.Adaptors) (result *models.RoleModel, errs []*validation.Error, err error) {

	role := &m.Role{
		Name:        roleParams.Name,
		Description: roleParams.Description,
	}

	// validation
	_, errs = role.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.Role.Add(role); err != nil {
		return
	}

	result = &models.RoleModel{}
	err = common.Copy(&result, &role)

	return
}

func GetRoleByName(name string, adaptors *adaptors.Adaptors) (role *m.Role, err error) {

	role, err = adaptors.Role.GetByName(name)

	return
}

func UpdateRole(roleParams *models.UpdateRole, adaptors *adaptors.Adaptors) (ok bool, errs []*validation.Error, err error) {

	role, err := adaptors.Role.GetByName(roleParams.Name)
	if err != nil {
		return
	}

	if roleParams.Parent.Name == "" {
		role.Parent = nil
	} else {
		role.Parent = &m.Role{
			Name: roleParams.Parent.Name,
		}
	}

	// validation
	ok, errs = role.Valid()
	if len(errs) > 0 || !ok {
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

func SearchRole(query string, limit, offset int, adaptors *adaptors.Adaptors) (roles []*m.Role, total int64, err error) {

	roles, total, err = adaptors.Role.Search(query, limit, offset)

	return
}

func GetAccessList(roleName string,
	adaptors *adaptors.Adaptors,
	accessListService *access_list.AccessListService) (accessList access_list.AccessList, err error) {

	var role *m.Role
	if role, err = adaptors.Role.GetByName(roleName); err != nil {
		return
	}

	accessList, err = accessListService.GetFullAccessList(role)

	return
}

func UpdateAccessList(roleName string, accessListDif map[string]map[string]bool, adaptors *adaptors.Adaptors) (err error) {

	var role *m.Role
	if role, err = adaptors.Role.GetByName(roleName); err != nil {
		return
	}

	addPerms := make([]*m.Permission, 0)
	delPerms := make([]string, 0)
	for packName, pack := range accessListDif {
		for levelName, dir := range pack {
			if dir {
				addPerms = append(addPerms, &m.Permission{
					RoleName:    role.Name,
					PackageName: packName,
					LevelName:   levelName,
				})
			} else {
				delPerms = append(delPerms, levelName)
			}

			if len(delPerms) > 0 {
				if err = adaptors.Permission.Delete(packName, delPerms); err != nil {
					return
				}
				delPerms = []string{}
			}
		}
	}

	if len(addPerms) == 0 {
		return
	}

	for _, perm := range addPerms {
		adaptors.Permission.Add(perm)
	}

	return
}
