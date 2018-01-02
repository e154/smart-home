package core

import "github.com/e154/smart-home/api/models"

type Nodes []*Node

type Node struct {
	model *models.Node
}
