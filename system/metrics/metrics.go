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

package metrics

import (
	"fmt"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/op/go-logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
)

var (
	log = logging.MustGetLogger("metrics")
)

type MetricManager struct {
	*Publisher
	cfg               *MetricConfig
	prometheusHandler http.Handler
	Cpu               *CpuManager
	Disk              *DiskManager
	Uptime            *UptimeManager
	Memory            *MemoryManager
	Gate              *GateManager
	Workflow          *WorkflowManager
	Node              *NodeManager
	Device            *DeviceManager
	MapElement        *MapElementManager
	Flow              *FlowManager
	graceful          *graceful_service.GracefulService
}

func NewMetricManager(cfg *MetricConfig,
	graceful *graceful_service.GracefulService) *MetricManager {
	metric := &MetricManager{
		Publisher:         NewPublisher(),
		cfg:               cfg,
		prometheusHandler: promhttp.Handler(),
		Disk:              NewDiskManager(),
		Uptime:            NewUptimeManager(),
		Memory:            NewMemoryManager(),
		graceful:          graceful,
	}

	metric.Gate = NewGateManager(metric)
	metric.Workflow = NewWorkflowManager(metric)
	metric.Node = NewNodeManager(metric)
	metric.Device = NewDeviceManager(metric)
	metric.MapElement = NewMapElementManager(metric)
	metric.Flow = NewFlowManager(metric)
	metric.Cpu = NewCpuManager(metric)

	return metric
}

func (m *MetricManager) Start() {

	m.Cpu.start(5)
	m.Disk.start(60)
	m.Uptime.start(15)
	m.Memory.start(3)

	m.graceful.Subscribe(m)

	if m.cfg.RunMode != config.DebugMode {
		return
	}

	log.Infof("Serving metric server at http://[::]:%d", m.cfg.Port)

	r := http.NewServeMux()

	// prometheus
	r.HandleFunc("/metrics", m.prometheusHandler.ServeHTTP)

	// Регистрация pprof-обработчиков
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port), r); err != nil {
		log.Fatal(err)
	}
}

func (m MetricManager) Shutdown() {

	m.Cpu.stop()
	m.Disk.stop()
	m.Uptime.stop()
	m.Memory.stop()
}

func (m *MetricManager) Update(t interface{}) {
	m.Workflow.update(t)
	m.Gate.update(t)
	m.Node.update(t)
	m.Device.update(t)
	m.MapElement.update(t)
	m.Flow.update(t)
}
