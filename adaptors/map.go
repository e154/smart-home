package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type Map struct {
	table *db.Maps
	db    *gorm.DB
}

func GetMapAdaptor(d *gorm.DB) *Map {
	return &Map{
		table: &db.Maps{Db: d},
		db:    d,
	}
}

func (n *Map) Add(ver *m.Map) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *Map) GetById(mapId int64) (ver *m.Map, err error) {

	var dbVer *db.Map
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *Map) Update(ver *m.Map) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *Map) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

func (n *Map) List(limit, offset int64, orderBy, sort string) (list []*m.Map, total int64, err error) {
	var dbList []*db.Map
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Map, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Map) Search(query string, limit, offset int) (list []*m.Map, total int64, err error) {
	var dbList []*db.Map
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Map, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Map) fromDb(dbVer *db.Map) (ver *m.Map) {
	ver = &m.Map{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		Description: dbVer.Description,
		Options:     dbVer.Options,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}

	return
}

func (n *Map) toDb(ver *m.Map) (dbVer *db.Map) {
	dbVer = &db.Map{
		Id:          ver.Id,
		Name:        ver.Name,
		Description: ver.Description,
		Options:     ver.Options,
	}
	return
}
