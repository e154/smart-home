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

package counters

import (
	"github.com/prometheus/client_golang/prometheus"
)

// CounterCollector ...
type CounterCollector struct {
	CounterDesc *prometheus.Desc
	value       float64
}

// Describe ...
func (c *CounterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.CounterDesc
}

// Collect ...
func (c *CounterCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.CounterDesc,
		prometheus.CounterValue,
		c.value,
	)
}

// Set ...
func (c *CounterCollector) Set(val float64) {
	c.value = val
}
