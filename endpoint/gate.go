package endpoint

import "github.com/e154/smart-home/system/gate_client"

type GateEndpoint struct {
	*CommonEndpoint
}

func NewGateEndpoint(common *CommonEndpoint) *GateEndpoint {
	return &GateEndpoint{
		CommonEndpoint: common,
	}
}

func (d *GateEndpoint) GetSettings() (settings *gate_client.Settings, err error) {
	settings, err = d.gate.GetSettings()
	return
}

func (d *GateEndpoint) UpdateSettings(settings *gate_client.Settings) (err error) {
	if err = d.gate.UpdateSettings(settings); err != nil {
		return
	}

	if settings.Enabled {
		d.gate.Connect()
	} else {
		d.gate.Close()
	}

	return
}

func (d *GateEndpoint) GetMobileList() (list *gate_client.MobileList, err error) {
	list, err = d.gate.GetMobileList()
	return
}

func (d *GateEndpoint) DeleteMobile(token string) (list *gate_client.MobileList, err error) {
	list, err = d.gate.DeleteMobile(token)
	return
}

func (d *GateEndpoint) AddMobile() (list *gate_client.MobileList, err error) {
	list, err = d.gate.AddMobile()
	return
}
