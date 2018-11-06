package db

import (
	"time"
	"github.com/jinzhu/gorm"
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
	UpdateAt    *time.Time
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
