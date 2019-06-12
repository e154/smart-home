package db

import (
	"time"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/e154/smart-home/common"
)

type Logs struct {
	Db *gorm.DB
}

type Log struct {
	Id        int64 `gorm:"primary_key"`
	Body      string
	Level     common.LogLevel
	CreatedAt time.Time
}

type LogQuery struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Levels    []string   `json:"levels"`
}

func (m *Log) TableName() string {
	return "logs"
}

func (n Logs) Add(v *Log) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		return
	}
	id = v.Id
	return
}

func (n Logs) GetById(mapId int64) (v *Log, err error) {
	v = &Log{Id: mapId}
	err = n.Db.First(&v).Error
	return
}

func (n Logs) Delete(mapId int64) (err error) {
	err = n.Db.Delete(&Log{Id: mapId}).Error
	return
}

func (n *Logs) List(limit, offset int64, orderBy, sort string, queryObj *LogQuery) (list []*Log, total int64, err error) {

	if err = n.Db.Model(Log{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Log, 0)
	q := n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy))

	if queryObj != nil {
		if queryObj.StartDate != nil {
			q = q.Where("created_at >= ?", &queryObj.StartDate)
		}
		if queryObj.EndDate != nil {
			q = q.Where("created_at <= ?", &queryObj.EndDate)
		}
		if len(queryObj.Levels) > 0 {
			q = q.Where("level in (?)", queryObj.Levels)
		}
	}

	err = q.
		Find(&list).
		Error

	if err != nil {
		log.Error(err.Error())
	}

	return
}

func (n *Logs) Search(query string, limit, offset int) (list []*Log, total int64, err error) {

	q := n.Db.Model(&Log{}).
		Where("body LIKE ?", "%"+query+"%").
		Order("body ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Log, 0)
	err = q.Find(&list).Error

	return
}
