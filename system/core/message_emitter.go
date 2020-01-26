package core

//ActionPrototypes
type MessageEmitter struct{}

func (m *MessageEmitter) After(flow *Flow) (err error) {
	//log.Infof("MessageEmitter.after: %v", message)
	return
}

func (m *MessageEmitter) Run(flow *Flow) (err error) {
	//log.Infof("MessageEmitter.run: %v", message)
	return
}

func (m *MessageEmitter) Before(flow *Flow) (err error) {
	//log.Infof("MessageEmitter.before: %v", message)
	return
}

func (m *MessageEmitter) Type() string {
	return "MessageEmitter"
}
