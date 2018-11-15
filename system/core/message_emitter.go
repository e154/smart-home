package core

//ActionPrototypes
type MessageEmitter struct{}

func (m *MessageEmitter) After(message *Message, flow *Flow) (err error) {
	//log.Info("MessageEmitter.after: ", message)
	return
}

func (m *MessageEmitter) Run(message *Message, flow *Flow) (err error) {
	//log.Info("MessageEmitter.run: ", message)
	return
}

func (m *MessageEmitter) Before(message *Message, flow *Flow) (err error) {
	//log.Info("MessageEmitter.before: ", message)
	return
}

func (m *MessageEmitter) Type() string {
	return "MessageEmitter"
}
