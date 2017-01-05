package core

import (
	"../models"
)

func NewMessage() *Message {
	return &Message{
		Data: make(map[string]interface{}),
	}
}
//TODO refactor message system
type Message struct {
	Device      	*models.Device
	Flow        	*models.Flow
	Node        	*models.Node
	Error       	string
	Device_state	func(state string)
	Data		map[string]interface{}
}
//TODO refactor message system
func (m *Message) clearError() {
	m.Error = ""
}
//TODO refactor message system
func (m *Message) SetError(err string) {
	m.Error = err
}
