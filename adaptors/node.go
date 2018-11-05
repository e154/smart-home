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

func (n *Node) fromDb(dbNode *db.Node) (node *m.Node) {
	node = m.NewNode()
	node.Id = dbNode.Id
	node.Name = dbNode.Name
	node.Ip = dbNode.Ip
	node.Port = dbNode.Port
	node.Status = dbNode.Status
	node.Description = dbNode.Description
	node.CreatedAt = dbNode.CreatedAt
	node.UpdateAt = dbNode.UpdateAt
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
		CreatedAt:   node.CreatedAt,
		UpdateAt:    node.UpdateAt,
	}
	return
}
