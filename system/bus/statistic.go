// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package bus

import (
	"go.uber.org/atomic"
	"strings"
	"time"
)

type Statistic struct {
	min *atomic.Duration
	max *atomic.Duration
	avg *atomic.Duration
	rps *RPSCounter
}

func NewStatistic() *Statistic {
	return &Statistic{
		min: atomic.NewDuration(0),
		max: atomic.NewDuration(0),
		avg: atomic.NewDuration(0),
		rps: startRPSCounter(),
	}
}

func (s *Statistic) setTime(t time.Duration) {
	if s.min.Load() == 0 {
		s.min.Store(t)
	}
	if s.min.Load() > t {
		s.min.Store(t)
	}
	if s.max.Load() == 0 {
		s.max.Store(t)
	}
	if t > s.max.Load() {
		s.max.Store(t)
	}
	s.avg.Store((s.max.Load() + s.min.Load()) / 2)
}

// StatItem ...
type StatItem struct {
	Topic       string
	Subscribers int
	Min         time.Duration
	Max         time.Duration
	Avg         time.Duration
	Rps         float64
}

// Stats ...
type Stats []StatItem

func (s Stats) Len() int           { return len(s) }
func (s Stats) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Stats) Less(i, j int) bool { return strings.Compare(s[i].Topic, s[j].Topic) == -1 }

type RPSCounter struct {
	count     *atomic.Int64
	value     *atomic.Float64
	isRunning *atomic.Bool
}

func startRPSCounter() *RPSCounter {
	counter := &RPSCounter{
		count:     atomic.NewInt64(0),
		value:     atomic.NewFloat64(0),
		isRunning: atomic.NewBool(true),
	}

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()
		for counter.isRunning.Load() {
			select {
			case <-ticker.C:
				counter.value.Store(float64(counter.count.Load()) / 5)
				counter.count.Store(0)
			}
		}
	}()

	return counter
}

func (c *RPSCounter) Inc() {
	c.count.Inc()
}

func (c *RPSCounter) Value() float64 {
	return c.value.Load()
}

func (c *RPSCounter) Stop() {
	c.isRunning.Store(false)
}
