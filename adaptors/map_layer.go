package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"sort"
)

type MapLayer struct {
	table *db.MapLayers
	db    *gorm.DB
}

func GetMapLayerAdaptor(d *gorm.DB) *MapLayer {
	return &MapLayer{
		table: &db.MapLayers{Db: d},
		db:    d,
	}
}

func (n *MapLayer) Add(ver *m.MapLayer) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *MapLayer) GetById(mapId int64) (ver *m.MapLayer, err error) {

	var dbVer *db.MapLayer
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *MapLayer) Update(ver *m.MapLayer) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *MapLayer) Sort(ver *m.MapLayer) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Sort(dbVer)
	return
}

func (n *MapLayer) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

func (n *MapLayer) List(limit, offset int64, orderBy, sort string) (list []*m.MapLayer, total int64, err error) {
	var dbList []*db.MapLayer
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.MapLayer, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *MapLayer) fromDb(dbVer *db.MapLayer) (ver *m.MapLayer) {
	ver = &m.MapLayer{
		Id:          dbVer.Id,
		Name:        dbVer.Name,
		MapId:       dbVer.MapId,
		Status:      dbVer.Status,
		Weight:      dbVer.Weight,
		Description: dbVer.Description,
		Elements:    make([]*m.MapElement, 0),
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}

	// elements
	mapElementAdaptor := GetMapElementAdaptor(n.db)
	for _, dbElement := range dbVer.Elements {
		element := mapElementAdaptor.fromDb(dbElement)
		ver.Elements = append(ver.Elements, element)
	}

	// map
	if dbVer.Map != nil {
		mapAdaptor := GetMapAdaptor(n.db)
		ver.Map = mapAdaptor.fromDb(dbVer.Map)
	}

	sort.Sort(m.SortMapElementByWeight(ver.Elements))

	return
}

func (n *MapLayer) toDb(ver *m.MapLayer) (dbVer *db.MapLayer) {
	dbVer = &db.MapLayer{
		Id:          ver.Id,
		Name:        ver.Name,
		MapId:       ver.MapId,
		Status:      ver.Status,
		Weight:      ver.Weight,
		Description: ver.Description,
	}
	return
}
