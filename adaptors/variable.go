package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type Variable struct {
	table *db.Variables
	db    *gorm.DB
}

func GetVariableAdaptor(d *gorm.DB) *Variable {
	return &Variable{
		table: &db.Variables{Db: d},
		db:    d,
	}
}

func (n *Variable) Add(variables *m.Variable) (err error) {

	dbVariable := n.toDb(variables)
	err = n.table.Add(dbVariable)

	return
}

func (n *Variable) GetAllEnabled() (list []*m.Variable, err error) {

	var dbList []*db.Variable
	if dbList, err = n.table.GetAllEnabled(); err != nil {
		return
	}

	list = make([]*m.Variable, 0)
	for _, dbVariable := range dbList {
		variables := n.fromDb(dbVariable)
		list = append(list, variables)
	}

	return
}

func (n *Variable) GetByName(name string) (variables *m.Variable, err error) {

	var dbVariable *db.Variable
	if dbVariable, err = n.table.GetByName(name); err != nil {
		return
	}

	variables = n.fromDb(dbVariable)

	return
}

func (n *Variable) Update(variable *m.Variable) (err error) {
	if _, err = n.table.GetByName(variable.Name); err != nil {
		err = n.Add(variable)
		return
	}
	dbVariable := n.toDb(variable)
	err = n.table.Update(dbVariable)
	return
}

func (n *Variable) Delete(name string) (err error) {
	err = n.table.Delete(name)
	return
}

func (n *Variable) List(limit, offset int64, orderBy, sort string) (list []*m.Variable, total int64, err error) {
	var dbList []*db.Variable
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Variable, 0)
	for _, dbVariable := range dbList {
		variables := n.fromDb(dbVariable)
		list = append(list, variables)
	}

	return
}

func (n *Variable) fromDb(dbVariable *db.Variable) (variables *m.Variable) {
	variables = &m.Variable{
		Name:      dbVariable.Name,
		Value:     dbVariable.Value,
		Autoload:  dbVariable.Autoload,
		CreatedAt: dbVariable.CreatedAt,
		UpdatedAt: dbVariable.UpdatedAt,
	}

	return
}

func (n *Variable) toDb(variables *m.Variable) (dbVariable *db.Variable) {
	dbVariable = &db.Variable{
		Name:     variables.Name,
		Value:    variables.Value,
		Autoload: variables.Autoload,
	}
	return
}
