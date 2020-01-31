package core

//ActionPrototypes
type MessageHandler struct{}

func (m *MessageHandler) After(flow *Flow) (err error) {
	//log.Infof("MessageHandler.after: %v", message)
	return
}

func (m *MessageHandler) Run(flow *Flow) (err error) {
	//log.Infof("MessageHandler.run: %v", message)
	return
}

func (m *MessageHandler) Before(flow *Flow) (err error) {
	//log.Infof("MessageHandler.before: %v", message)
	return
}

func (m *MessageHandler) Type() string {
	return "MessageHandler"
}
