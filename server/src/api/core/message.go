package core

import (
	"../models"
)

type Message struct {
	Device		*models.Device
	Flow		*models.Flow
	Node		*models.Node
	Result		[]byte
	ResultStr	string
}