// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package scripts

import (
	"sync"
)

// Pull ...
type Pull struct {
	heap sync.Map
}

// NewPull ...
func NewPull() *Pull {
	return &Pull{
		heap: sync.Map{},
	}
}

// Get ...
func (p *Pull) Get(name string) (value interface{}, ok bool) {
	value, ok = p.heap.Load(name)
	return
}

// Push ...
func (p *Pull) Push(name string, s interface{}) {
	p.heap.Store(name, s)
}

// Pop ...
func (p *Pull) Pop(name string) {
	p.heap.Delete(name)
}

// Purge ...
func (p *Pull) Purge() {
	p.heap.Range(func(key, value interface{}) bool {
		p.heap.Delete(key)
		return true
	})
}

func (p *Pull) Range(f func(key, value interface{}) bool) {
	p.heap.Range(f)
}
