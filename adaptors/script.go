package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/copier"
)

type Script struct {
	table *db.Scripts
	db    *gorm.DB
}

func GetScriptAdaptor(d *gorm.DB) *Script {
	return &Script{
		table: &db.Scripts{Db: d},
		db:    d,
	}
}

func (n *Script) Add(node *m.Script) (id int64, err error) {

	var dbScript *db.Script
	if dbScript, err = n.toDb(node); err != nil {
		return
	}

	id, err = n.table.Add(dbScript)
	return
}

func (n *Script) GetById(nodeId int64) (node *m.Script, err error) {

	var dbScript *db.Script
	if dbScript, err = n.table.GetById(nodeId); err != nil {
		return
	}

	node, _ = n.fromDb(dbScript)

	return
}

func (n *Script) Update(node *m.Script) (err error) {
	var dbScript *db.Script
	if dbScript, err = n.toDb(node); err != nil {
		return
	}
	err = n.table.Update(dbScript)
	return
}

func (n *Script) Delete(nodeId int64) (err error) {
	err = n.table.Delete(nodeId)
	return
}

func (n *Script) List(limit, offset int64, orderBy, sort string) (list []*m.Script, total int64, err error) {
	var dbList []*db.Script
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Script, 0)
	for _, dbScript := range dbList {
		node, _ := n.fromDb(dbScript)
		list = append(list, node)
	}

	return
}

func (n *Script) fromDb(dbScript *db.Script) (node *m.Script, err error) {
	node = &m.Script{}
	err = copier.Copy(&node, &dbScript)
	return
}

func (n *Script) toDb(node *m.Script) (dbScript *db.Script, err error) {
	dbScript = &db.Script{}
	err = copier.Copy(&dbScript, &node)
	return
}
