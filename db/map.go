package db

import (
	"time"
	"github.com/jinzhu/gorm"
	"encoding/json"
	"fmt"
)

type Maps struct {
	Db *gorm.DB
}

type Map struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Options     json.RawMessage `gorm:"type:jsonb;not null"`
	Layers      []*MapLayer
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *Map) TableName() string {
	return "maps"
}

func (n Maps) Add(v *Map) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n Maps) GetById(mapId int64) (v *Map, err error) {
	v = &Map{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

func (n Maps) Update(m *Map) (err error) {
	err = n.Db.Model(&Map{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"options":     m.Options,
	}).Error
	return
}

func (n Maps) Delete(mapId int64) (err error) {
	err = n.Db.Delete(&Map{Id: mapId}).Error
	return
}

func (n *Maps) List(limit, offset int64, orderBy, sort string) (list []*Map, total int64, err error) {

	if err = n.Db.Model(Map{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Map, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}

func (n *Maps) Search(query string, limit, offset int) (list []*Map, total int64, err error) {

	q := n.Db.Model(&Map{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Map, 0)
	err = q.Find(&list).Error

	return
}
