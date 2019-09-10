package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type ZoneTag struct {
	table *db.ZoneTags
	db    *gorm.DB
}

func GetZoneTagAdaptor(d *gorm.DB) *ZoneTag {
	return &ZoneTag{
		table: &db.ZoneTags{Db: d},
		db:    d,
	}
}

func (n *ZoneTag) Add(tag *m.ZoneTag) (id int64, err error) {

	dbTag := n.toDb(tag)
	id, err = n.table.Add(dbTag)

	return
}

func (n *ZoneTag) Delete(name string) (err error) {

	err = n.table.Delete(name)

	return
}

func (n *ZoneTag) Search(query string, limit, offset int) (list []*m.ZoneTag, total int64, err error) {
	var dbList []*db.ZoneTag
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.ZoneTag, 0)
	for _, dbTag := range dbList {
		node := n.fromDb(dbTag)
		list = append(list, node)
	}

	return
}

func (n *ZoneTag) toDb(tag *m.ZoneTag) *db.ZoneTag {
	return &db.ZoneTag{
		Id:   tag.Id,
		Name: tag.Name,
	}
}

func (n *ZoneTag) fromDb(tag *db.ZoneTag) *m.ZoneTag {
	return &m.ZoneTag{
		Id:   tag.Id,
		Name: tag.Name,
	}
}
