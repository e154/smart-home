package models

import "log"

//ActionPrototypes
type MessageHandler struct {}

func (m *MessageHandler) After(message *Message, flow *Flow) (err error) {
	log.Println("MessageHandler.after: ", message)
	return
}

func (m *MessageHandler) Run(message *Message, flow *Flow) (err error) {
	log.Println("MessageHandler.run: ", message)
	return
}

func (m *MessageHandler) Before(message *Message, flow *Flow) (err error) {
	log.Println("MessageHandler.before: ", message)
	return
}

func (m *MessageHandler) Type() string {
	return  "MessageHandler"
}