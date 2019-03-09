package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type Log struct {
	table *db.Logs
	db    *gorm.DB
}

func GetLogAdaptor(d *gorm.DB) *Log {
	return &Log{
		table: &db.Logs{Db: d},
		db:    d,
	}
}

func (n *Log) Add(ver *m.Log) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *Log) GetById(verId int64) (ver *m.Log, err error) {

	var dbVer *db.Log
	if dbVer, err = n.table.GetById(verId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Log) Delete(verId int64) (err error) {
	err = n.table.Delete(verId)
	return
}

func (n *Log) List(limit, offset int64, orderBy, sort string, queryObj *m.LogQuery) (list []*m.Log, total int64, err error) {

	var dbList []*db.Log
	var dbQueryObj *db.LogQuery

	if queryObj != nil {
		dbQueryObj = &db.LogQuery{
			StartDate: queryObj.StartDate,
			EndDate:   queryObj.EndDate,
			Levels:    queryObj.Levels,
		}
	}

	if dbList, total, err = n.table.List(limit, offset, orderBy, sort, dbQueryObj); err != nil {
		return
	}

	list = make([]*m.Log, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Log) Search(query string, limit, offset int) (list []*m.Log, total int64, err error) {
	var dbList []*db.Log
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Log, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Log) fromDb(dbVer *db.Log) (ver *m.Log) {
	ver = &m.Log{
		Id:        dbVer.Id,
		Body:      dbVer.Body,
		Level:     dbVer.Level,
		CreatedAt: dbVer.CreatedAt,
	}

	return
}

func (n *Log) toDb(ver *m.Log) (dbVer *db.Log) {
	dbVer = &db.Log{
		Body:  ver.Body,
		Level: ver.Level,
	}
	return
}
