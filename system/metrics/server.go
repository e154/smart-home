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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
)

type MetricServer struct {
	cfg               *MetricConfig
	prometheusHandler http.Handler
}

func NewMetricServer(cfg *MetricConfig,
	graceful *graceful_service.GracefulService) *MetricServer {
	metric := &MetricServer{
		cfg:               cfg,
		prometheusHandler: promhttp.Handler(),
	}

	//graceful.Subscribe(metric)
	return metric
}

func (m MetricServer) Start() {

	if m.cfg.RunMode == config.ReleaseMode {
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

