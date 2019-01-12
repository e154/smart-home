package db

import (
	"time"
	"github.com/jinzhu/gorm"
	"fmt"
)

type Nodes struct {
	Db *gorm.DB
}

type Node struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Ip          string
	Port        int
	Status      string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *Node) TableName() string {
	return "nodes"
}

func (n Nodes) Add(node *Node) (id int64, err error) {
	if err = n.Db.Create(&node).Error; err != nil {
		return
	}
	id = node.Id
	return
}

func (n Nodes) GetAllEnabled() (list []*Node, err error) {
	list = make([]*Node, 0)
	err = n.Db.Where("status = ?", "enabled").
		Find(&list).Error
	return
}

func (n Nodes) GetById(nodeId int64) (node *Node, err error) {
	node = &Node{Id: nodeId}
	err = n.Db.First(&node).Error
	return
}

func (n Nodes) Update(m *Node) (err error) {
	err = n.Db.Model(&Node{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"ip":          m.Ip,
		"status":      m.Status,
		"port":        m.Port,
	}).Error
	return
}

func (n Nodes) Delete(nodeId int64) (err error) {
	err = n.Db.Delete(&Node{Id: nodeId}).Error
	return
}

func (n *Nodes) List(limit, offset int64, orderBy, sort string) (list []*Node, total int64, err error) {

	if err = n.Db.Model(Node{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Node, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}

func (n *Nodes) Search(query string, limit, offset int) (list []*Node, total int64, err error) {

	fmt.Println(query)
	q := n.Db.Model(&Node{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Node, 0)
	err = q.Find(&list).Error

	return
}