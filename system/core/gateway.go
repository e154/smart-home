package core

//ActionPrototypes
type Gateway struct{}

func (m *Gateway) After(flow *Flow) (err error) {
	//log.Infof("Gateway.after: %v", message)
	return
}

func (m *Gateway) Run(flow *Flow) (err error) {
	//log.Infof("Gateway.run: %v", message)
	return
}

func (m *Gateway) Before(flow *Flow) (err error) {
	//log.Infof("Gateway.before: %v", message)
	return
}

func (m *Gateway) Type() string {
	return "Gateway"
}
