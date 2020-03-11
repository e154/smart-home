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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
)

var (
	log = common.MustGetLogger("metrics")
)

type MetricManager struct {
	*Publisher
	cfg               *MetricConfig
	adaptors          *adaptors.Adaptors
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
	AppMemory         *AppMemoryManager
	Mqtt              *MqttManager
	Zigbee2Mqtt       *Zigbee2MqttManager
	graceful          *graceful_service.GracefulService
}

func NewMetricManager(cfg *MetricConfig,
	graceful *graceful_service.GracefulService,
	adaptors *adaptors.Adaptors) *MetricManager {
	metric := &MetricManager{
		adaptors:          adaptors,
		Publisher:         NewPublisher(),
		cfg:               cfg,
		prometheusHandler: promhttp.Handler(),
		graceful:          graceful,
	}

	metric.Gate = NewGateManager(metric)
	metric.Workflow = NewWorkflowManager(metric, adaptors)
	metric.Node = NewNodeManager(metric)
	metric.Device = NewDeviceManager(metric)
	metric.MapElement = NewMapElementManager(metric)
	metric.Flow = NewFlowManager(metric, adaptors)
	metric.Cpu = NewCpuManager(metric)
	metric.Disk = NewDiskManager(metric)
	metric.Uptime = NewUptimeManager(metric)
	metric.Memory = NewMemoryManager(metric)
	metric.AppMemory = NewAppMemoryManager(metric)
	metric.Mqtt = NewMqttManager(metric)
	metric.Zigbee2Mqtt = NewZigbee2MqttManager(metric)

	return metric
}

func (m *MetricManager) Start() {

	m.Cpu.start(5)
	m.Disk.start(60)
	m.Uptime.start(15)
	m.Memory.start(5)
	m.AppMemory.start(20)

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
		log.Fatal(err.Error())
	}
}

func (m MetricManager) Shutdown() {

	m.Cpu.stop()
	m.Disk.stop()
	m.Uptime.stop()
	m.Memory.stop()
	m.AppMemory.stop()
}

func (m *MetricManager) Update(t interface{}) {
	m.Workflow.update(t)
	m.Gate.update(t)
	m.Node.update(t)
	m.Device.update(t)
	m.MapElement.update(t)
	m.Flow.update(t)
	m.Mqtt.update(t)
	m.Zigbee2Mqtt.update(t)
}
