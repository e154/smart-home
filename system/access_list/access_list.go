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

type AccessListService struct {
	List     *AccessList
	adaptors *adaptors.Adaptors
}

func NewAccessListService(adaptors *adaptors.Adaptors) *AccessListService {
	accessList := &AccessListService{
		adaptors: adaptors,
	}
	accessList.ReadConfig()
	return accessList
}

func (a *AccessListService) ReadConfig() (err error) {

	//var file []byte
	//file, err = ioutil.ReadFile(path)
	//if err != nil {
	//	log.Fatal("Error reading config file")
	//	return
	//}

	a.List = &AccessList{}
	err = json.Unmarshal([]byte(DATA), a.List)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	return
}

func (a *AccessListService) GetFullAccessList(role *m.Role) (accessList AccessList, err error) {

	var permissions []*m.Permission
	if permissions, err = a.adaptors.Permission.GetAllPermissions(role.Name); err != nil {
		return
	}

	accessList = make(AccessList)
	var item AccessItem
	var levels AccessLevels
	var ok bool
	list := *a.List
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

func (a *AccessListService) GetShotAccessList(role *m.Role) (err error) {

	err = a.adaptors.Role.GetAccessList(role)
	return
}
