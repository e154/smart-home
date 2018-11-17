package core

import (
	"github.com/e154/smart-home/system/scripts"
	m "github.com/e154/smart-home/models"
)

type Action struct {
	Device  *m.Device
	//Node    *Node
	Script  *scripts.Engine
	Message *Message
}
