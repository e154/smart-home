package core

//ActionPrototypes
type FlowLink struct {}

func (m *FlowLink) After(message *Message, flow *Flow) (err error) {
	//log.Infof("FlowLink.after: %v", message)
	return
}

func (m *FlowLink) Run(message *Message, flow *Flow) (err error) {
	//log.Infof("FlowLink.run: %v", message)
	return
}

func (m *FlowLink) Before(message *Message, flow *Flow) (err error) {
	//log.Infof("FlowLink.before: %v", message)
	return
}

func (m *FlowLink) Type() string {
	return  "FlowLink"
}