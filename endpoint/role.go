package endpoint

import (
	"errors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/validation"
)

type RoleCommand struct {
	*CommonCommand
}

func NewRoleCommand(common *CommonCommand) *RoleCommand {
	return &RoleCommand{
		CommonCommand: common,
	}
}

func (n *RoleCommand) Add(params *m.Role) (result *m.Role, errs []*validation.Error, err error) {

	// validation
	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Role.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.Role.GetByName(params.Name)

	return
}

func (n *RoleCommand) GetByName(name string) (result *m.Role, err error) {

	result, err = n.adaptors.Role.GetByName(name)

	return
}

func (n *RoleCommand) Update(params *m.Role) (result *m.Role, errs []*validation.Error, err error) {

	role, err := n.adaptors.Role.GetByName(params.Name)
	if err != nil {
		return
	}

	if params.Parent.Name == "" {
		role.Parent = nil
	} else {
		role.Parent = &m.Role{
			Name: params.Parent.Name,
		}
	}

	// validation
	_, errs = role.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Role.Update(role); err != nil {
		return
	}

	role, err = n.adaptors.Role.GetByName(params.Name)

	return
}

func (n *RoleCommand) GetList(limit, offset int64, order, sortBy string) (result []*m.Role, total int64, err error) {

	result, total, err = n.adaptors.Role.List(limit, offset, order, sortBy)

	return
}

func (n *RoleCommand) Delete(name string) (err error) {

	if name == "admin" {
		err = errors.New("admin is base role")
		return
	}
	err = n.adaptors.Role.Delete(name)

	return
}

func (n *RoleCommand) Search(query string, limit, offset int) (result []*m.Role, total int64, err error) {

	result, total, err = n.adaptors.Role.Search(query, limit, offset)

	return
}

func (n *RoleCommand) GetAccessList(roleName string,
	accessListService *access_list.AccessListService) (accessList access_list.AccessList, err error) {

	var role *m.Role
	if role, err = n.adaptors.Role.GetByName(roleName); err != nil {
		return
	}

	accessList, err = accessListService.GetFullAccessList(role)

	return
}

func (n *RoleCommand) UpdateAccessList(roleName string, accessListDif map[string]map[string]bool) (err error) {

	var role *m.Role
	if role, err = n.adaptors.Role.GetByName(roleName); err != nil {
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
				if err = n.adaptors.Permission.Delete(packName, delPerms); err != nil {
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
		n.adaptors.Permission.Add(perm)
	}

	return
}