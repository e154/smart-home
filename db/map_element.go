package db

import (
	"encoding/json"
	"fmt"
	. "github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"time"
)

type MapElements struct {
	Db *gorm.DB
}

type Prototype struct {
	*MapImage
	*MapText
	*MapDevice
}

type MapElement struct {
	Id            int64 `gorm:"primary_key"`
	Name          string
	Description   string
	PrototypeId   int64
	PrototypeType PrototypeType
	Prototype     Prototype
	Map           *Map
	MapId         int64
	MapLayer      *MapLayer
	MapLayerId    int64
	GraphSettings json.RawMessage `gorm:"type:jsonb;not null"`
	Status        StatusType
	Weight        int64
	ZoneId        *int64
	Zone          *ZoneTag
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
	if err = n.Db.
		Preload("Zone").
		First(&v).Error; err != nil {
		return
	}

	if v.PrototypeId == 0 {
		return
	}

	switch v.PrototypeType {
	case PrototypeTypeText:
		text := &MapText{}
		err = n.Db.Model(&MapText{}).
			Where("id = ?", v.PrototypeId).
			First(&text).
			Error

		if err == nil {
			v.Prototype = Prototype{
				MapText: text,
			}
		}
	case PrototypeTypeImage:
		image := &MapImage{}
		err = n.Db.Model(&MapImage{}).
			Where("id = ?", v.PrototypeId).
			First(&image).
			Error
		if err == nil {
			v.Prototype = Prototype{
				MapImage: image,
			}
		}
	case PrototypeTypeDevice:
		device := &MapDevice{}
		err = n.Db.Model(&MapDevice{}).
			Where("id = ?", v.PrototypeId).
			Preload("Image").
			Preload("States").
			Preload("States.Image").
			Preload("States.DeviceState").
			Preload("Actions").
			Preload("Actions.Image").
			Preload("Actions.DeviceAction").
			Preload("Device").
			Preload("Device.States").
			Preload("Device.Actions").
			First(&device).
			Error
		if err == nil {
			v.Prototype = Prototype{
				MapDevice: device,
			}
		}
	}

	return
}

func (n MapElements) Update(m *MapElement) (err error) {
	err = n.Db.Model(&MapElement{Id: m.Id}).Updates(map[string]interface{}{
		"name":           m.Name,
		"description":    m.Description,
		"prototype_id":   m.PrototypeId,
		"prototype_type": m.PrototypeType,
		"map_id":         m.MapId,
		"layer_id":       m.MapLayerId,
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

func (n *MapElements) List(limit, offset int64, orderBy, sort string) (list []*MapElement, total int64, err error) {

	if err = n.Db.Model(MapElement{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapElement, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}

func (n *MapElements) GetActiveElements(limit, offset int64, orderBy, sort string) (list []*MapElement, total int64, err error) {

	list = make([]*MapElement, 0)

	q := n.Db.Model(&MapElement{}).
		Where("prototype_type = 'device'")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	q = n.Db.Model(&MapElement{}).
		Where("prototype_type = 'device'").
		Preload("Zone").
		Limit(limit).
		Offset(offset)

	if orderBy != "" && sort != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	if err != nil {
		return
	}

	deviceIds := make([]int64, 0)
	for _, e := range list {
		deviceIds = append(deviceIds, e.PrototypeId)
	}

	devices := make([]*MapDevice, 0)
	if len(deviceIds) > 0 {
		err = n.Db.Model(&MapDevice{}).
			Where("id in (?)", deviceIds).
			Preload("Image").
			Preload("States").
			Preload("States.Image").
			Preload("States.DeviceState").
			Preload("Actions").
			Preload("Actions.Image").
			Preload("Actions.DeviceAction").
			Preload("Device").
			Preload("Device.States").
			Preload("Device.Actions").
			Find(&devices).
			Error

		if err != nil {
			return
		}
	}

	for _, e := range list {
		for _, device := range devices {
			if device.Id == e.PrototypeId {
				e.Prototype = Prototype{
					MapDevice: device,
				}
				continue
			}
		}
	}

	return
}
