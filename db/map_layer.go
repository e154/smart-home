package db

import (
	"time"
	"github.com/jinzhu/gorm"
	"fmt"
)

type MapLayers struct {
	Db *gorm.DB
}

type MapLayer struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Map         *Map
	MapId       int64
	Status      string
	Weight      int64
	Elements    []*MapElement
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *MapLayer) TableName() string {
	return "map_layers"
}

func (n MapLayers) Add(v *MapLayer) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n MapLayers) GetById(mapId int64) (v *MapLayer, err error) {
	v = &MapLayer{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

func (n MapLayers) Update(m *MapLayer) (err error) {
	err = n.Db.Model(&MapLayer{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"status":      m.Status,
		"weight":      m.Weight,
		"map_id":      m.MapId,
	}).Error
	return
}

func (n MapLayers) Sort(m *MapLayer) (err error) {
	err = n.Db.Model(&MapLayer{Id: m.Id}).Updates(map[string]interface{}{
		"weight": m.Weight,
	}).Error
	return
}

func (n MapLayers) Delete(mapId int64) (err error) {
	err = n.Db.Delete(&MapLayer{Id: mapId}).Error
	return
}

func (n *MapLayers) List(limit, offset int64, orderBy, sort string) (list []*MapLayer, total int64, err error) {

	if err = n.Db.Model(MapLayer{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapLayer, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}
