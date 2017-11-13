package core

import (
	"github.com/e154/smart-home/api/models"
)

type ActionBind struct {
	action *Action
}

func (a *ActionBind) Device() *models.Device {
	return a.action.GetDevice()
}
func (a *ActionBind) Node() *models.Node {
	return a.action.GetNode()
}