package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type MapImage struct {
	table *db.MapImages
	db    *gorm.DB
}

func GetMapImageAdaptor(d *gorm.DB) *MapImage {
	return &MapImage{
		table: &db.MapImages{Db: d},
		db:    d,
	}
}

func (n *MapImage) Add(ver *m.MapImage) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *MapImage) GetById(mapId int64) (ver *m.MapImage, err error) {

	var dbVer *db.MapImage
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *MapImage) Update(ver *m.MapImage) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *MapImage) Sort(ver *m.MapImage) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Sort(dbVer)
	return
}

func (n *MapImage) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

func (n *MapImage) List(limit, offset int64, orderBy, sort string) (list []*m.MapImage, total int64, err error) {
	var dbList []*db.MapImage
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.MapImage, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *MapImage) fromDb(dbVer *db.MapImage) (ver *m.MapImage) {
	ver = &m.MapImage{
		Id:        dbVer.Id,
		ImageId:   dbVer.ImageId,
		Style:     dbVer.Style,
	}

	return
}

func (n *MapImage) toDb(ver *m.MapImage) (dbVer *db.MapImage) {
	dbVer = &db.MapImage{
		Id:      ver.Id,
		ImageId: ver.ImageId,
		Style:   ver.Style,
	}
	return
}
