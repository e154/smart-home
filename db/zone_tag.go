package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type ZoneTags struct {
	Db *gorm.DB
}

type ZoneTag struct {
	Id   int64 `gorm:"primary_key"`
	Name string
}

func (d *ZoneTag) TableName() string {
	return "zone_tags"
}

func (n ZoneTags) Add(tag *ZoneTag) (id int64, err error) {
	if err = n.Db.Create(&tag).Error; err != nil {
		return
	}
	id = tag.Id
	return
}

func (n *ZoneTags) Search(query string, limit, offset int) (list []*ZoneTag, total int64, err error) {

	q := n.Db.Model(&ZoneTag{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*ZoneTag, 0)
	err = q.Find(&list).Error

	return
}

func (n ZoneTags) Delete(name string) (err error) {
	if name == "" {
		err = fmt.Errorf("zero name")
		return
	}
	err = n.Db.Model(&ZoneTag{}).
		Delete("name = ?", name).Error
	return
}
