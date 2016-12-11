package core

//ActionPrototypes
type MessageEmitter struct {}

func (m *MessageEmitter) After(message *Message, flow *Flow) (err error) {
	//log.Println("MessageEmitter.after: ", message)
	return
}

func (m *MessageEmitter) Run(message *Message, flow *Flow) (err error) {
	//log.Println("MessageEmitter.run: ", message)
	return
}

func (m *MessageEmitter) Before(message *Message, flow *Flow) (err error) {
	//log.Println("MessageEmitter.before: ", message)
	return
}

func (m *MessageEmitter) Type() string {
	return  "MessageEmitter"
}