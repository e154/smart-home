package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common/debug"
)

type UserMeta struct {
	table *db.UserMetas
	db    *gorm.DB
}

func GetUserMetaAdaptor(d *gorm.DB) *UserMeta {
	return &UserMeta{
		table: &db.UserMetas{Db: d},
		db:    d,
	}
}

func (n *UserMeta) UpdateOrCreate(meta *m.UserMeta) (id int64, err error) {

	dbMeta := n.toDb(meta)
	debug.Println(dbMeta)
	if id, err = n.table.UpdateOrCreate(dbMeta); err != nil {
		return
	}

	return
}

func (n *UserMeta) fromDb(dbMeta *db.UserMeta) (meta *m.UserMeta) {
	meta = &m.UserMeta{
		Id:     dbMeta.Id,
		Key:    dbMeta.Key,
		UserId: dbMeta.UserId,
		Value:  dbMeta.Value,
	}
	return
}

func (n *UserMeta) toDb(meta *m.UserMeta) (dbMeta *db.UserMeta) {
	dbMeta = &db.UserMeta{
		Id:     meta.Id,
		Key:    meta.Key,
		UserId: meta.UserId,
		Value:  meta.Value,
	}
	return
}
