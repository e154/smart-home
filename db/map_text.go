package db

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

type MapTexts struct {
	Db *gorm.DB
}

type MapText struct {
	Id        int64 `gorm:"primary_key"`
	Text      string
	Style     string
}

func (d *MapText) TableName() string {
	return "map_texts"
}

func (n MapTexts) Add(v *MapText) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n MapTexts) GetById(mapId int64) (v *MapText, err error) {
	v = &MapText{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

func (n MapTexts) Update(m *MapText) (err error) {
	err = n.Db.Model(&MapText{Id: m.Id}).Updates(map[string]interface{}{
		"text":  m.Text,
		"style": m.Style,
	}).Error
	return
}

func (n MapTexts) Sort(m *MapText) (err error) {
	err = n.Db.Model(&MapText{Id: m.Id}).Updates(map[string]interface{}{
		"text":  m.Text,
		"style": m.Style,
	}).Error
	return
}

func (n MapTexts) Delete(mapId int64) (err error) {
	err = n.Db.Delete(&MapText{Id: mapId}).Error
	return
}

func (n *MapTexts) List(limit, offset int64, orderBy, sort string) (list []*MapText, total int64, err error) {

	if err = n.Db.Model(MapText{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapText, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
