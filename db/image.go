package db

import (
	"time"
	"github.com/jinzhu/gorm"
	"fmt"
)

type Images struct {
	Db *gorm.DB
}

type Image struct {
	Id        int64 `gorm:"primary_key"`
	Thumb     string
	Url       string `gorm:"-"`
	Image     string
	MimeType  string
	Title     string
	Size      int64
	Name      string
	CreatedAt time.Time
}

func (m *Image) TableName() string {
	return "images"
}

func (n Images) Add(v *Image) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n Images) GetById(mapId int64) (v *Image, err error) {
	v = &Image{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

func (n Images) Update(m *Image) (err error) {
	err = n.Db.Model(&Image{Id: m.Id}).Updates(map[string]interface{}{
		"title": m.Title,
		"Name":  m.Name,
	}).Error
	return
}

func (n Images) Delete(mapId int64) (err error) {
	err = n.Db.Delete(&Image{Id: mapId}).Error
	return
}

func (n *Images) List(limit, offset int64, orderBy, sort string) (list []*Image, total int64, err error) {

	if err = n.Db.Model(Image{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Image, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
