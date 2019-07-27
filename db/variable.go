package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Variables struct {
	Db *gorm.DB
}

type Variable struct {
	Name      string `gorm:"primary_key"`
	Value     string
	Autoload  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *Variable) TableName() string {
	return "variables"
}

func (n Variables) Add(variable *Variable) (err error) {
	err = n.Db.Create(&variable).Error
	return
}

func (n Variables) GetByName(name string) (variable *Variable, err error) {
	variable = &Variable{}
	err = n.Db.Model(variable).
		Where("name = ?", name).
		First(&variable).
		Error
	return
}

func (n Variables) GetAllEnabled() (list []*Variable, err error) {
	list = make([]*Variable, 0)
	err = n.Db.Where("autoload = ?", true).
		Find(&list).Error
	return
}

func (n Variables) Update(m *Variable) (err error) {
	err = n.Db.Model(&Variable{Name: m.Name}).Updates(map[string]interface{}{
		"value":    m.Value,
		"autoload": m.Autoload,
	}).Error
	return
}

func (n Variables) Delete(name string) (err error) {
	err = n.Db.Delete(&Variable{Name: name}).Error
	return
}

func (n *Variables) List(limit, offset int64, orderBy, sort string) (list []*Variable, total int64, err error) {

	if err = n.Db.Model(Variable{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Variable, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
