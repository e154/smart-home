package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/uuid"
)

type Connection struct {
	table *db.Connections
	db    *gorm.DB
}

func GetConnectionAdaptor(d *gorm.DB) *Connection {
	return &Connection{
		table: &db.Connections{Db: d},
		db:    d,
	}
}

func (n *Connection) Add(con *m.Connection) (id uuid.UUID, err error) {
	dbConnection := n.toDb(con)
	if id, err = n.table.Add(dbConnection); err != nil {
		return
	}

	return
}

func (n *Connection) GetById(conId uuid.UUID) (con *m.Connection, err error) {

	var dbConnection *db.Connection
	if dbConnection, err = n.table.GetById(conId); err != nil {
		return
	}

	con = n.fromDb(dbConnection)

	return
}

func (n *Connection) Update(con *m.Connection) (err error) {
	dbConnection := n.toDb(con)
	err = n.table.Update(dbConnection)
	return
}

func (n *Connection) Delete(conIds []uuid.UUID) (err error) {
	err = n.table.Delete(conIds)
	return
}

func (n *Connection) List(limit, offset int64, orderBy, sort string) (list []*m.Connection, total int64, err error) {
	var dbList []*db.Connection
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Connection, 0)
	for _, dbConnection := range dbList {
		con := n.fromDb(dbConnection)
		list = append(list, con)
	}

	return
}

func (n *Connection) AddOrUpdateConnection(connection *m.Connection) (err error) {

	if connection.Uuid.String() == "00000000-0000-0000-0000-000000000000" {
		_, err = n.Add(connection)
		return
	}

	if _, err = n.table.GetById(connection.Uuid); err != nil {
		_, err = n.Add(connection)
		return
	}

	err = n.Update(connection)

	return
}

func (n *Connection) fromDb(dbConnection *db.Connection) (con *m.Connection) {
	con = &m.Connection{
		Uuid:          dbConnection.Uuid,
		Name:          dbConnection.Name,
		FlowId:        dbConnection.FlowId,
		GraphSettings: dbConnection.GraphSettings,
		Direction:     dbConnection.Direction,
		ElementFrom:   dbConnection.ElementFrom,
		ElementTo:     dbConnection.ElementTo,
		PointFrom:     dbConnection.PointFrom,
		PointTo:       dbConnection.PointTo,
		CreatedAt:     dbConnection.CreatedAt,
		UpdatedAt:     dbConnection.UpdatedAt,
	}
	return
}

func (n *Connection) toDb(con *m.Connection) (dbConnection *db.Connection) {
	dbConnection = &db.Connection{
		Uuid:          con.Uuid,
		Name:          con.Name,
		PointTo:       con.PointTo,
		PointFrom:     con.PointFrom,
		ElementTo:     con.ElementTo,
		ElementFrom:   con.ElementFrom,
		Direction:     con.Direction,
		GraphSettings: con.GraphSettings,
		FlowId:        con.FlowId,
	}
	if dbConnection.GraphSettings == nil {
		dbConnection.GraphSettings.UnmarshalJSON([]byte("{}"))
	}
	return
}
