package counters

import (
	"github.com/prometheus/client_golang/prometheus"
)

type CounterCollector struct {
	CounterDesc *prometheus.Desc
	value float64
}

func (c *CounterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.CounterDesc
}

func (c *CounterCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.CounterDesc,
		prometheus.CounterValue,
		c.value,
	)
}

func (c *CounterCollector) Set(val float64) {
	c.value = val
}

