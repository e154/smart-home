package core

//ActionPrototypes
type FlowLink struct {}

func (m *FlowLink) After(message *Message, flow *Flow) (err error) {
	//log.Info("FlowLink.after: ", message)
	return
}

func (m *FlowLink) Run(message *Message, flow *Flow) (err error) {
	//log.Info("FlowLink.run: ", message)
	return
}

func (m *FlowLink) Before(message *Message, flow *Flow) (err error) {
	//log.Info("FlowLink.before: ", message)
	return
}

func (m *FlowLink) Type() string {
	return  "FlowLink"
}