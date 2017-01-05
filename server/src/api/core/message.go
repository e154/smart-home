package core

import (
	"../models"
)

type Message struct {
	Device       *models.Device
	Flow         *models.Flow
	Node         *models.Node
	Result       string
	Error        string
	Device_state func(state string)
}