package mqtt

import (
	"github.com/e154/smart-home/system/mqtt/gmqtt"
	"github.com/e154/smart-home/system/config"
	"time"
)

type MqttConfig struct {
	Port                       int
	RetryInterval              time.Duration
	RetryCheckInterval         time.Duration
	SessionExpiryInterval      time.Duration
	SessionExpireCheckInterval time.Duration
	QueueQos0Messages          bool
	MaxInflight                int
	MaxAwaitRel                int
	MaxMsgQueue                int
	DeliverMode                gmqtt.DeliverMode
}

func NewMqttConfig(cfg *config.AppConfig) *MqttConfig {
	return &MqttConfig{
		Port:                       cfg.MqttPort,
		RetryInterval:              cfg.MqttRetryInterval,
		RetryCheckInterval:         cfg.MqttRetryCheckInterval,
		SessionExpiryInterval:      cfg.MqttSessionExpiryInterval,
		SessionExpireCheckInterval: cfg.MqttSessionExpireCheckInterval,
		QueueQos0Messages:          cfg.MqttQueueQos0Messages,
		MaxInflight:                cfg.MqttMaxInflight,
		MaxAwaitRel:                cfg.MqttMaxAwaitRel,
		MaxMsgQueue:                cfg.MqttMaxMsgQueue,
		DeliverMode:                gmqtt.DeliverMode(cfg.MqttDeliverMode),
	}
}
