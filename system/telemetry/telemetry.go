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

package telemetry

import (
	"github.com/op/go-logging"
	"sync"
)

var (
	log = logging.MustGetLogger("telemetry")
)

type Telemetry struct {
	sync.Mutex
	subscribers map[string]ITelemetry
}

func NewTelemetry() (t2 *Telemetry, t1 ITelemetry) {
	t2 = &Telemetry{
		subscribers: make(map[string]ITelemetry),
	}
	t1 = t2
	return
}

func (s *Telemetry) Subscribe(command string, f ITelemetry) {
	log.Infof("subscribe %s", command)
	s.Lock()
	defer s.Unlock()
	if s.subscribers[command] != nil {
		delete(s.subscribers, command)
	}
	s.subscribers[command] = f
}

func (s *Telemetry) UnSubscribe(command string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.subscribers[command]; ok {
		delete(s.subscribers, command)
	}
}

func (t *Telemetry) Broadcast(param interface{}) {
	t.Lock()
	defer t.Unlock()
	for _, f := range t.subscribers {
		f.Broadcast(param)
	}
}

func (t *Telemetry) BroadcastOne(param interface{}) {
	t.Lock()
	defer t.Unlock()
	for _, f := range t.subscribers {
		f.BroadcastOne(param)
	}
}
