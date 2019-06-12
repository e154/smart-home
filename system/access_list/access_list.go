package access_list

import (
	"io/ioutil"
	"encoding/json"
	"github.com/op/go-logging"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
)

var (
	log = logging.MustGetLogger("access_list")
)

type AccessListService struct {
	List     *AccessList
	adaptors *adaptors.Adaptors
}

func NewAccessListService(adaptors *adaptors.Adaptors) *AccessListService {
	accessList := &AccessListService{
		adaptors: adaptors,
	}
	accessList.ReadConfig("./conf/access_list.json")
	return accessList
}

func (a *AccessListService) ReadConfig(path string) (err error) {

	var file []byte
	file, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading config file")
		return
	}

	a.List = &AccessList{}
	err = json.Unmarshal(file, a.List)
	if err != nil {
		log.Fatal("Error: wrong format of config file")
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
