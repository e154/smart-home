package mobile

import (
	"github.com/e154/smart-home/system/config"
)

type MobileServerConfig struct {
	Host    string
	Port    int
	RunMode config.RunMode
}

func NewMobileServerConfig(cfg *config.AppConfig) *MobileServerConfig {
	return &MobileServerConfig{
		Host: cfg.ServerHost,
		Port: cfg.ServerPort + 1,
		RunMode: cfg.Mode,
	}
}
