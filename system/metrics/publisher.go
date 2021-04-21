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
	"sync"
)

// ISubscriber ...
type ISubscriber interface {
	Broadcast(interface{})
}

// IPublisher ...
type IPublisher interface {
	Broadcast(interface{})
}

// Publisher ...
type Publisher struct {
	sync.Mutex
	subscribers map[string]ISubscriber
}

// NewPublisher ...
func NewPublisher() (t *Publisher) {
	t = &Publisher{
		subscribers: make(map[string]ISubscriber),
	}

	return
}

// Subscribe ...
func (p *Publisher) Subscribe(command string, f ISubscriber) {
	p.Lock()
	defer p.Unlock()
	if p.subscribers[command] != nil {
		delete(p.subscribers, command)
	}
	p.subscribers[command] = f
}

// UnSubscribe ...
func (p *Publisher) UnSubscribe(command string) {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.subscribers[command]; ok {
		delete(p.subscribers, command)
	}
}

// Broadcast ...
func (p *Publisher) Broadcast(param interface{}) {
	p.Lock()
	defer p.Unlock()
	for _, f := range p.subscribers {
		f.Broadcast(param)
	}
}
