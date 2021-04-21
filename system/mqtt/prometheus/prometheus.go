// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package prometheus

import (
	"context"
	"github.com/e154/smart-home/common"
	"net/http"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/persistence/subscription"
)

var (
	log = common.MustGetLogger("prometheus")
)

const name = "prometheus"
const metricPrefix = "gmqtt_"

// Prometheus served as a prometheus exporter that exposes gmqtt metrics.
type Prometheus struct {
	statsManager gmqtt.StatsManager
	httpServer   *http.Server
	path         string
}

// New ...
func New(httpSever *http.Server, path string) *Prometheus {
	p := &Prometheus{
		httpServer: httpSever,
		path:       path,
	}
	return p
}

// Load ...
func (p *Prometheus) Load(service gmqtt.Server) error {
	p.statsManager = service.GetStatsManager()
	r := prometheus.NewPedanticRegistry()
	r.MustRegister(p)
	mu := http.NewServeMux()
	mu.Handle(p.path, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	p.httpServer.Handler = mu
	go func() {
		err := p.httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err.Error())
		}
	}()
	return nil
}

// Unload ...
func (p *Prometheus) Unload() error {
	return p.httpServer.Shutdown(context.Background())
}

// HookWrapper ...
func (p *Prometheus) HookWrapper() gmqtt.HookWrapper {
	return gmqtt.HookWrapper{}
}

// Name ...
func (p *Prometheus) Name() string {
	return name
}

// Describe ...
func (p *Prometheus) Describe(desc chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(p, desc)
}

// Collect ...
func (p *Prometheus) Collect(m chan<- prometheus.Metric) {
	log.Debug("metrics collected")
	st := p.statsManager.GetStats()
	collectPacketsStats(st.PacketStats, m)
	collectClientStats(st.ClientStats, m)
	collectSubscriptionStats(st.SubscriptionStats, m)
	collectMessageStats(st.MessageStats, m)
}

func collectPacketsStats(ps *gmqtt.PacketStats, m chan<- prometheus.Metric) {
	bytesReceivedMetricName := metricPrefix + "packets_received_bytes_total"
	ReceivedCounterMetricName := metricPrefix + "packets_received_total"
	bytesSentMetricName := metricPrefix + "packets_sent_bytes_total"
	sentCounterMetricName := metricPrefix + "packets_sent_total"

	collectPacketsStatsBytes(bytesReceivedMetricName, ps.BytesReceived, m)
	collectPacketsStatsBytes(bytesSentMetricName, ps.BytesSent, m)

	collectPacketsStatsCounter(ReceivedCounterMetricName, ps.ReceivedTotal, m)
	collectPacketsStatsCounter(sentCounterMetricName, ps.SentTotal, m)
}
func collectPacketsStatsBytes(metricName string, pb *gmqtt.PacketBytes, m chan<- prometheus.Metric) {
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Connect)),
		"CONNECT",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Connack)),
		"CONNACK",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Disconnect)),
		"DISCONNECT",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Pingreq)),
		"PINGREQ",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Pingresp)),
		"PINGRESP",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Puback)),
		"PUBACK",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Pubcomp)),
		"PUBCOMP",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Publish)),
		"PUBLISH",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Pubrec)),
		"PUBREC",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Pubrel)),
		"PUBREL",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Suback)),
		"SUBACK",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Subscribe)),
		"SUBSCRIBE",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Unsuback)),
		"UNSUBACK",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pb.Unsubscribe)),
		"UNSUBSCRIBE",
	)
}
func collectPacketsStatsCounter(metricName string, pc *gmqtt.PacketCount, m chan<- prometheus.Metric) {
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Connect)),
		"CONNECT",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Connack)),
		"CONNACK",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Disconnect)),
		"DISCONNECT",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Pingreq)),
		"PINGREQ",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Pingresp)),
		"PINGRESP",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Puback)),
		"PUBACK",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Pubcomp)),
		"PUBCOMP",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Publish)),
		"PUBLISH",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Pubrec)),
		"PUBREC",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Pubrel)),
		"PUBREL",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Suback)),
		"SUBACK",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Subscribe)),
		"SUBSCRIBE",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Unsuback)),
		"UNSUBACK",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"type"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&pc.Unsubscribe)),
		"UNSUBSCRIBE",
	)
}

func collectClientStats(c *gmqtt.ClientStats, m chan<- prometheus.Metric) {
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricPrefix+"clients_connected_total", "", nil, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&c.ConnectedTotal)),
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricPrefix+"sessions_active_current", "", nil, nil),
		prometheus.GaugeValue,
		float64(atomic.LoadUint64(&c.ActiveCurrent)),
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricPrefix+"sessions_inactive_current", "", nil, nil),
		prometheus.GaugeValue,
		float64(atomic.LoadUint64(&c.InactiveCurrent)),
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricPrefix+"clients_disconnected_total", "", nil, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&c.DisconnectedTotal)),
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricPrefix+"sessions_expired_total", "", nil, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&c.ExpiredTotal)),
	)
}
func collectMessageStats(ms *gmqtt.MessageStats, m chan<- prometheus.Metric) {
	collectMessageStatsDropped(ms, m)
	collectMessageStatsQueued(ms, m)
	collectMessageStatsReceived(ms, m)
	collectMessageStatsSent(ms, m)
}
func collectMessageStatsDropped(ms *gmqtt.MessageStats, m chan<- prometheus.Metric) {
	metricName := metricPrefix + "messages_dropped_total"
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"qos"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&ms.Qos0.DroppedTotal)), "0",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"qos"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&ms.Qos1.DroppedTotal)), "1",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"qos"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&ms.Qos2.DroppedTotal)), "2",
	)
}

func collectMessageStatsQueued(ms *gmqtt.MessageStats, m chan<- prometheus.Metric) {
	metricName := "messages_queued_current"
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", nil, nil),
		prometheus.GaugeValue,
		float64(atomic.LoadUint64(&ms.QueuedCurrent)),
	)
}
func collectMessageStatsReceived(ms *gmqtt.MessageStats, m chan<- prometheus.Metric) {
	metricName := "messages_received_total"
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"qos"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&ms.Qos0.ReceivedTotal)), "0",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"qos"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&ms.Qos1.ReceivedTotal)), "1",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"qos"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&ms.Qos2.ReceivedTotal)), "2",
	)
}
func collectMessageStatsSent(ms *gmqtt.MessageStats, m chan<- prometheus.Metric) {
	metricName := "messages_sent_total"
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"qos"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&ms.Qos0.SentTotal)), "0",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"qos"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&ms.Qos1.SentTotal)), "1",
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", []string{"qos"}, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&ms.Qos2.SentTotal)), "2",
	)
}

func collectSubscriptionStats(s *subscription.Stats, m chan<- prometheus.Metric) {
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricPrefix+"subscriptions_total", "", nil, nil),
		prometheus.CounterValue,
		float64(atomic.LoadUint64(&s.SubscriptionsTotal)),
	)
	m <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricPrefix+"subscriptions_current", "", nil, nil),
		prometheus.GaugeValue,
		float64(atomic.LoadUint64(&s.SubscriptionsCurrent)),
	)
}
