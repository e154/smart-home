package db

import (
	"time"
	"github.com/jinzhu/gorm"
	"fmt"
	"database/sql"
)

type Roles struct {
	Db *gorm.DB
}

type Role struct {
	Name        string `gorm:"primary_key"`
	Description string
	Role        *Role
	RoleName    sql.NullString `gorm:"column:parent"`
	Children    []*Role
	Permissions []*Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *Role) TableName() string {
	return "roles"
}

func (n Roles) Add(role *Role) (err error) {
	if err = n.Db.Create(&role).Error; err != nil {
		return
	}
	return
}

func (n Roles) GetByName(name string) (role *Role, err error) {

	role = &Role{Name: name}
	err = n.Db.First(&role).Error
	if err != nil {
		return
	}

	err = n.RelData(role)

	return
}

func (n Roles) Update(m *Role) (err error) {
	err = n.Db.Model(&Role{Name: m.Name}).Updates(map[string]interface{}{
		"description": m.Description,
		"parent":      m.RoleName,
	}).Error
	return
}

func (n Roles) Delete(name string) (err error) {
	err = n.Db.Delete(&Role{Name: name}).Error
	return
}

func (n *Roles) List(limit, offset int64, orderBy, sort string) (list []*Role, total int64, err error) {

	if err = n.Db.Model(Role{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Role, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	if err != nil {
		return
	}

	for _, role := range list {
		n.RelData(role)
	}

	return
}

func (n *Roles) Search(query string, limit, offset int) (list []*Role, total int64, err error) {

	fmt.Println(query)
	q := n.Db.Model(&Role{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Role, 0)
	err = q.Find(&list).Error

	return
}

func (n *Roles) RelData(role *Role) (err error) {

	// get parent
	if role.RoleName.Valid {
		role.Role = &Role{}
		err = n.Db.Model(role).
			Where("name = ?", role.RoleName.String).
			Find(&role.Role).
			Error
	}

	// get children
	role.Children = make([]*Role, 0)
	err = n.Db.Model(role).
		Where("parent = ?", role.Name).
		Find(&role.Children).
		Error

	return
}