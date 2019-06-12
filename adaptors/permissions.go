package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type Permission struct {
	table *db.Permissions
	db    *gorm.DB
}

func GetPermissionAdaptor(d *gorm.DB) *Permission {
	return &Permission{
		table: &db.Permissions{Db: d},
		db:    d,
	}
}

func (n *Permission) Add(permission *m.Permission) (id int64, err error) {

	dbPermission := n.toDb(permission)
	if id, err = n.table.Add(dbPermission); err != nil {
		return
	}

	return
}

func (n *Permission) Delete(packageName string, levelName []string) (err error) {

	err = n.table.Delete(packageName, levelName)

	return
}

func (n *Permission) GetAllPermissions(roleName string) (permissions []*m.Permission, err error) {

	var dbPermissions []*db.Permission
	if dbPermissions, err = n.table.GetAllPermissions(roleName); err != nil {
		return
	}

	for _, dbVer := range dbPermissions {
		ver := n.fromDb(dbVer)
		permissions = append(permissions, ver)
	}

	return
}

func (n *Permission) fromDb(dbPermission *db.Permission) (permission *m.Permission) {
	permission = &m.Permission{
		Id:          dbPermission.Id,
		RoleName:    dbPermission.RoleName,
		PackageName: dbPermission.PackageName,
		LevelName:   dbPermission.LevelName,
	}

	return
}

func (n *Permission) toDb(permission *m.Permission) (dbPermission *db.Permission) {
	dbPermission = &db.Permission{
		RoleName:    permission.RoleName,
		LevelName:   permission.LevelName,
		PackageName: permission.PackageName,
	}
	return
}
