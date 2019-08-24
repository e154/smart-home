package endpoint

import (
	"context"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/gin-gonic/gin"
)

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
		d.gate.Close()
		d.gate.Connect()
	} else {
		d.gate.Close()
	}

	return
}

func (d *GateEndpoint) GetMobileList(ctx *gin.Context) (list *gate_client.MobileList, err error) {
	list, err = d.gate.GetMobileList(ctx)
	return
}

func (d *GateEndpoint) DeleteMobile(token string, ctx context.Context) (list *gate_client.MobileList, err error) {
	list, err = d.gate.DeleteMobile(token, ctx)
	return
}

func (d *GateEndpoint) AddMobile(ctx context.Context) (list *gate_client.MobileList, err error) {
	list, err = d.gate.AddMobile(ctx)
	return
}
