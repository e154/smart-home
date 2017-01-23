package core

//ActionPrototypes
type MessageHandler struct {}

func (m *MessageHandler) After(message *Message, flow *Flow) (err error) {
	//log.Info("MessageHandler.after: ", message)
	return
}

func (m *MessageHandler) Run(message *Message, flow *Flow) (err error) {
	//log.Info("MessageHandler.run: ", message)
	return
}

func (m *MessageHandler) Before(message *Message, flow *Flow) (err error) {
	//log.Info("MessageHandler.before: ", message)
	return
}

func (m *MessageHandler) Type() string {
	return  "MessageHandler"
}