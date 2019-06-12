package mqtt

import (
	"github.com/e154/smart-home/system/config"
)

type MqttConfig struct {
	SrvKeepAlive        int
	SrvConnectTimeout   int
	SrvSessionsProvider string
	SrvTopicsProvider   string
	SrvPort             int
}

func NewMqttConfig(cfg *config.AppConfig) *MqttConfig {
	return &MqttConfig{
		SrvKeepAlive:        cfg.MqttKeepAlive,
		SrvConnectTimeout:   cfg.MqttConnectTimeout,
		SrvSessionsProvider: cfg.MqttSessionsProvider,
		SrvTopicsProvider:   cfg.MqttTopicsProvider,
		SrvPort:             cfg.MqttPort,
	}
}
