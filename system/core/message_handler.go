package core

//ActionPrototypes
type MessageHandler struct{}

func (m *MessageHandler) After(message *Message, flow *Flow) (err error) {
	log.Infof("MessageHandler.after: %v", message)
	return
}

func (m *MessageHandler) Run(message *Message, flow *Flow) (err error) {
	log.Infof("MessageHandler.run: %v", message)
	return
}

func (m *MessageHandler) Before(message *Message, flow *Flow) (err error) {
	log.Infof("MessageHandler.before: %v", message)
	return
}

func (m *MessageHandler) Type() string {
	return "MessageHandler"
}
