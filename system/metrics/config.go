package metrics

import "github.com/e154/smart-home/system/config"

type MetricConfig struct {
	RunMode config.RunMode
	Host    string
	Port    int
}

func NewMetricConfig(cfg *config.AppConfig) *MetricConfig {
	return &MetricConfig{
		RunMode: cfg.Mode,
		Host:    "0.0.0.0",
		Port:    cfg.MetricPort,
	}
}
