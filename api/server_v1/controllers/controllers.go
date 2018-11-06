package controllers

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
)

type Controllers struct {
	Index *ControllerIndex
	Node  *ControllerNode
}

func NewControllers(adaptors *adaptors.Adaptors,
	core *core.Core) *Controllers {
	common := NewControllerCommon(adaptors, core)
	return &Controllers{
		Index: NewControllerIndex(common),
		Node:  NewControllerNode(common),
	}
}
