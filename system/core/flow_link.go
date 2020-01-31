package core

//ActionPrototypes
type FlowLink struct {}

func (m *FlowLink) After(flow *Flow) (err error) {
	//log.Infof("FlowLink.after: %v", message)
	return
}

func (m *FlowLink) Run(flow *Flow) (err error) {
	//log.Infof("FlowLink.run: %v", message)
	return
}

func (m *FlowLink) Before(flow *Flow) (err error) {
	//log.Infof("FlowLink.before: %v", message)
	return
}

func (m *FlowLink) Type() string {
	return  "FlowLink"
}
