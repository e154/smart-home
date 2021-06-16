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

package access_list

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

var (
	log = common.MustGetLogger("access_list")
)

type AccessListService interface {
	ReadConfig() (err error)
	GetFullAccessList(roleName string) (accessList AccessList, err error)
	GetShotAccessList(role *m.Role) (err error)
	List() *AccessList
}

// accessListService ...
type accessListService struct {
	list     *AccessList
	adaptors *adaptors.Adaptors
}

// NewAccessListService ...
func NewAccessListService(adaptors *adaptors.Adaptors) AccessListService {
	accessList := &accessListService{
		adaptors: adaptors,
	}
	accessList.ReadConfig()
	return accessList
}

// ReadConfig ...
func (a *accessListService) ReadConfig() (err error) {

	//var file []byte
	//file, err = ioutil.ReadFile(path)
	//if err != nil {
	//	log.Fatal("Error reading config file")
	//	return
	//}

	a.list = &AccessList{}
	err = json.Unmarshal([]byte(DATA), a.list)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	return
}

// GetFullAccessList ...
func (a *accessListService) GetFullAccessList(roleName string) (accessList AccessList, err error) {

	var permissions []*m.Permission
	if permissions, err = a.adaptors.Permission.GetAllPermissions(roleName); err != nil {
		return
	}

	accessList = make(AccessList)
	var item AccessItem
	var levels AccessLevels
	var ok bool
	list := *a.list
	for _, perm := range permissions {

		if levels, ok = list[perm.PackageName]; !ok {
			continue
		}

		if accessList[perm.PackageName] == nil {
			accessList[perm.PackageName] = NewAccessLevels()
		}

		if item, ok = levels[perm.LevelName]; !ok {
			continue
		}

		item.RoleName = perm.RoleName
		accessList[perm.PackageName][perm.LevelName] = item
	}

	return
}

// GetShotAccessList ...
func (a *accessListService) GetShotAccessList(role *m.Role) (err error) {

	err = a.adaptors.Role.GetAccessList(role)
	return
}

func (a *accessListService) List() *AccessList {
	return a.list
}
