package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type Node struct {
	table *db.Nodes
	db    *gorm.DB
}

func GetNodeAdaptor(d *gorm.DB) *Node {
	return &Node{
		table: &db.Nodes{Db: d},
		db:    d,
	}
}

func (n *Node) Add(node *m.Node) (id int64, err error) {

	dbNode := n.toDb(node)
	if id, err = n.table.Add(dbNode); err != nil {
		return
	}

	return
}

func (n *Node) GetAllEnabled() (list []*m.Node, err error) {

	var dbList []*db.Node
	if dbList, err = n.table.GetAllEnabled(); err != nil {
		return
	}

	list = make([]*m.Node, 0)
	for _, dbNode := range dbList {
		node := n.fromDb(dbNode)
		list = append(list, node)
	}

	return
}

func (n *Node) GetById(nodeId int64) (node *m.Node, err error) {

	var dbNode *db.Node
	if dbNode, err = n.table.GetById(nodeId); err != nil {
		return
	}

	node = n.fromDb(dbNode)

	return
}

func (n *Node) Update(node *m.Node) (err error) {
	dbNode := n.toDb(node)
	err = n.table.Update(dbNode)
	return
}

func (n *Node) Delete(nodeId int64) (err error) {
	err = n.table.Delete(nodeId)
	return
}

func (n *Node) List(limit, offset int64, orderBy, sort string) (list []*m.Node, total int64, err error) {
	var dbList []*db.Node
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Node, 0)
	for _, dbNode := range dbList {
		node := n.fromDb(dbNode)
		list = append(list, node)
	}

	return
}

func (n *Node) Search(query string, limit, offset int) (list []*m.Node, total int64, err error) {
	var dbList []*db.Node
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Node, 0)
	for _, dbNode := range dbList {
		node := n.fromDb(dbNode)
		list = append(list, node)
	}

	return
}

func (n *Node) fromDb(dbNode *db.Node) (node *m.Node) {
	node = &m.Node{
		Id: dbNode.Id,
		Name: dbNode.Name,
		Ip: dbNode.Ip,
		Port: dbNode.Port,
		Status: dbNode.Status,
		Description: dbNode.Description,
		CreatedAt: dbNode.CreatedAt,
		UpdatedAt: dbNode.UpdatedAt,
	}

	return
}

func (n *Node) toDb(node *m.Node) (dbNode *db.Node) {
	dbNode = &db.Node{
		Id:          node.Id,
		Name:        node.Name,
		Ip:          node.Ip,
		Port:        node.Port,
		Status:      node.Status,
		Description: node.Description,
	}
	return
}
