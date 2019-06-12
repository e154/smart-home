package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type MapText struct {
	table *db.MapTexts
	db    *gorm.DB
}

func GetMapTextAdaptor(d *gorm.DB) *MapText {
	return &MapText{
		table: &db.MapTexts{Db: d},
		db:    d,
	}
}

func (n *MapText) Add(ver *m.MapText) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *MapText) GetById(mapId int64) (ver *m.MapText, err error) {

	var dbVer *db.MapText
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *MapText) Update(ver *m.MapText) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *MapText) Sort(ver *m.MapText) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Sort(dbVer)
	return
}

func (n *MapText) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

func (n *MapText) List(limit, offset int64, orderBy, sort string) (list []*m.MapText, total int64, err error) {
	var dbList []*db.MapText
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.MapText, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *MapText) fromDb(dbVer *db.MapText) (ver *m.MapText) {
	ver = &m.MapText{
		Id:        dbVer.Id,
		Text:      dbVer.Text,
		Style:     dbVer.Style,
	}

	return
}

func (n *MapText) toDb(ver *m.MapText) (dbVer *db.MapText) {
	dbVer = &db.MapText{
		Id:    ver.Id,
		Text:  ver.Text,
		Style: ver.Style,
	}
	return
}
