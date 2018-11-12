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

func (n *Script) Add(script *m.Script) (id int64, err error) {

	var dbScript *db.Script
	if dbScript, err = n.toDb(script); err != nil {
		return
	}

	id, err = n.table.Add(dbScript)
	return
}

func (n *Script) GetById(scriptId int64) (script *m.Script, err error) {

	var dbScript *db.Script
	if dbScript, err = n.table.GetById(scriptId); err != nil {
		return
	}

	script, _ = n.fromDb(dbScript)

	return
}

func (n *Script) Update(script *m.Script) (err error) {
	var dbScript *db.Script
	if dbScript, err = n.toDb(script); err != nil {
		return
	}
	err = n.table.Update(dbScript)
	return
}

func (n *Script) Delete(scriptId int64) (err error) {
	err = n.table.Delete(scriptId)
	return
}

func (n *Script) List(limit, offset int64, orderBy, sort string) (list []*m.Script, total int64, err error) {
	var dbList []*db.Script
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Script, 0)
	for _, dbScript := range dbList {
		script, _ := n.fromDb(dbScript)
		list = append(list, script)
	}

	return
}

func (n *Script) fromDb(dbScript *db.Script) (script *m.Script, err error) {
	script = &m.Script{}
	err = copier.Copy(&script, &dbScript)
	return
}

func (n *Script) toDb(script *m.Script) (dbScript *db.Script, err error) {
	dbScript = &db.Script{}
	err = copier.Copy(&dbScript, &script)
	return
}
