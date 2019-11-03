package endpoint

import (
	"github.com/e154/smart-home/system/notify"
)

type NotifyEndpoint struct {
	*CommonEndpoint
}

func NewNotifyEndpoint(common *CommonEndpoint) *NotifyEndpoint {
	return &NotifyEndpoint{
		CommonEndpoint: common,
	}
}

func (n *NotifyEndpoint) GetSettings() (cfg *notify.NotifyConfig, err error) {
	cfg = n.notify.GetCfg()
	return
}

func (n *NotifyEndpoint) UpdateSettings(cfg *notify.NotifyConfig) (err error) {
	if err = n.notify.UpdateCfg(cfg); err != nil {
		return
	}

	n.notify.Restart()
	return
}
