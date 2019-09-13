package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type MapZone struct {
	table *db.MapZones
	db    *gorm.DB
}

func GetMapZoneAdaptor(d *gorm.DB) *MapZone {
	return &MapZone{
		table: &db.MapZones{Db: d},
		db:    d,
	}
}

func (n *MapZone) Add(tag *m.MapZone) (id int64, err error) {

	dbTag := n.toDb(tag)
	id, err = n.table.Add(dbTag)

	return
}


func (n *MapZone) GetByName(zoneName string) (ver *m.MapZone, err error) {

	var dbVer *db.MapZone
	if dbVer, err = n.table.GetByName(zoneName); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *MapZone) Delete(name string) (err error) {

	err = n.table.Delete(name)

	return
}

func (n *MapZone) Search(query string, limit, offset int) (list []*m.MapZone, total int64, err error) {
	var dbList []*db.MapZone
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.MapZone, 0)
	for _, dbTag := range dbList {
		node := n.fromDb(dbTag)
		list = append(list, node)
	}

	return
}

func (n *MapZone) toDb(tag *m.MapZone) *db.MapZone {
	return &db.MapZone{
		Id:   tag.Id,
		Name: tag.Name,
	}
}

func (n *MapZone) fromDb(tag *db.MapZone) *m.MapZone {
	return &m.MapZone{
		Id:   tag.Id,
		Name: tag.Name,
	}
}
