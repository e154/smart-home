package db

import (
	"github.com/jinzhu/gorm"
	"time"
	"encoding/json"
	. "github.com/e154/smart-home/common"
)

type MapElements struct {
	Db *gorm.DB
}

type MapElement struct {
	Id            int64 `gorm:"primary_key"`
	Name          string
	Description   string
	PrototypeId   int64
	PrototypeType PrototypeType
	Map           *Map
	MapId         int64
	Layer         *MapLayer
	LayerId       int64
	GraphSettings json.RawMessage `gorm:"type:jsonb;not null"`
	Status        StatusType
	Weight        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (d *MapElement) TableName() string {
	return "map_elements"
}

func (n MapElements) Add(v *MapElement) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n MapElements) GetById(mapId int64) (v *MapElement, err error) {
	v = &MapElement{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

func (n MapElements) Update(m *MapElement) (err error) {
	err = n.Db.Model(&MapElement{Id: m.Id}).Updates(map[string]interface{}{
		"name":           m.Name,
		"description":    m.Description,
		"prototype_id":   m.PrototypeId,
		"prototype_type": m.PrototypeType,
		"map_id":         m.MapId,
		"layer_id":       m.LayerId,
		"graph_settings": m.GraphSettings,
		"status":         m.Status,
		"weight":         m.Weight,
	}).Error
	return
}

func (n MapElements) Sort(m *MapElement) (err error) {
	err = n.Db.Model(&MapElement{Id: m.Id}).Updates(map[string]interface{}{
		"weight": m.Weight,
	}).Error
	return
}

func (n MapElements) Delete(mapId int64) (err error) {
	err = n.Db.Delete(&MapElement{Id: mapId}).Error
	return
}