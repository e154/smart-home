package db

import (
	"time"
	"github.com/jinzhu/gorm"
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/common"
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

func (n Maps) GetFullById(mapId int64) (v *Map, err error) {

	v = &Map{}
	err = n.Db.Model(v).
		Where("id = ?", mapId).
		Preload("Layers").
		Preload("Layers.Elements").
		Preload("Layers.Map").
		Find(&v).
		Error

	if err != nil {
		return
	}

	imageIds := make([]int64, 0)
	textIds := make([]int64, 0)
	deviceIds := make([]int64, 0)
	for _, l := range v.Layers {
		for _, e := range l.Elements {
			switch e.PrototypeType {
			case common.PrototypeTypeText:
				textIds = append(textIds, e.Id)
			case common.PrototypeTypeImage:
				imageIds = append(imageIds, e.Id)
			case common.PrototypeTypeDevice:
				deviceIds = append(deviceIds, e.Id)
			}
		}
	}

	images := make([]*MapImage, 0)
	texts := make([]*MapText, 0)
	devices := make([]*MapDevice, 0)

	if len(imageIds) > 0 {
		err = n.Db.Model(&MapImage{}).
			Where("id in (?)", imageIds).
			Find(&images).
			Error
	}

	if len(textIds) > 0 {
		err = n.Db.Model(&MapText{}).
			Where("id in (?)", textIds).
			Find(&texts).
			Error
	}

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
	}

	//debug.Println(devices)
	//fmt.Println("---")

	for _, l := range v.Layers {
		for _, e := range l.Elements {
			switch e.PrototypeType {
			case common.PrototypeTypeText:
				for _, text := range texts {
					if text.Id == e.PrototypeId {
						e.Prototype = Prototype{
							MapText: text,
						}
						continue
					}
				}
			case common.PrototypeTypeImage:
				for _, image := range images {
					if image.Id == e.PrototypeId {
						e.Prototype = Prototype{
							MapImage: image,
						}
						continue
					}
				}
			case common.PrototypeTypeDevice:
				for _, device := range devices {
					if device.Id == e.PrototypeId {
						e.Prototype = Prototype{
							MapDevice: device,
						}
						continue
					}
				}
			}
		}
	}

	//debug.Println(v)
	//fmt.Println("---")

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
