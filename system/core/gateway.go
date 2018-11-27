package core

//ActionPrototypes
type Gateway struct{}

func (m *Gateway) After(message *Message, flow *Flow) (err error) {
	//log.Infof("Gateway.after: %v", message)
	return
}

func (m *Gateway) Run(message *Message, flow *Flow) (err error) {
	//log.Infof("Gateway.run: %v", message)
	return
}

func (m *Gateway) Before(message *Message, flow *Flow) (err error) {
	//log.Infof("Gateway.before: %v", message)
	return
}

func (m *Gateway) Type() string {
	return "Gateway"
}
