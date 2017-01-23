package core

//ActionPrototypes
type Gateway struct {}

func (m *Gateway) After(message *Message, flow *Flow) (err error) {
	//log.Info("Gateway.after: ", message)
	return
}

func (m *Gateway) Run(message *Message, flow *Flow) (err error) {
	//log.Info("Gateway.run: ", message)
	return
}

func (m *Gateway) Before(message *Message, flow *Flow) (err error) {
	//log.Info("Gateway.before: ", message)
	return
}

func (m *Gateway) Type() string {
	return  "Gateway"
}