package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type MapZones struct {
	Db *gorm.DB
}

type MapZone struct {
	Id   int64 `gorm:"primary_key"`
	Name string
}

func (d *MapZone) TableName() string {
	return "map_zones"
}

func (n MapZones) Add(zone *MapZone) (id int64, err error) {
	if err = n.Db.Create(&zone).Error; err != nil {
		return
	}
	id = zone.Id
	return
}

func (n MapZones) GetByName(zoneName string) (zone *MapZone, err error) {

	zone = &MapZone{}
	err = n.Db.Model(zone).
		Where("name = ?", zoneName).
		First(&zone).
		Error

	return
}

func (n *MapZones) Search(query string, limit, offset int) (list []*MapZone, total int64, err error) {

	q := n.Db.Model(&MapZone{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapZone, 0)
	err = q.Find(&list).Error

	return
}

func (n MapZones) Delete(name string) (err error) {
	if name == "" {
		err = fmt.Errorf("zero name")
		return
	}

	err = n.Db.Delete(&MapZone{}, "name = ?", name).Error
	return
}
