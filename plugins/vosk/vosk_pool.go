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

//go:build !test
// +build !test

package vosk

import "sync"

type Pool struct {
	min int
	v   *Vosk
	sync.Mutex
	pool []*Worker
}

func NewPool(v *Vosk, min int) *Pool {
	var pool = make([]*Worker, 0, min)
	for i := 0; i < min; i++ {
		worker, _ := NewWorker(v)
		pool = append(pool, worker)
	}
	return &Pool{
		v:    v,
		min:  min,
		pool: pool,
	}
}

func (p *Pool) Release(worker *Worker) {
	p.Lock()
	defer p.Unlock()
	for _, w := range p.pool {
		if w == worker {
			w.InUse = false
		}
	}
}

func (p *Pool) GetWorker() *Worker {
	p.Lock()
	defer p.Unlock()
	for _, w := range p.pool {
		if !w.InUse {
			w.InUse = true
			return w
		}
	}
	worker, _ := NewWorker(p.v)
	p.pool = append(p.pool, worker)
	return worker
}

func (p *Pool) Free() {
	p.Lock()
	defer p.Unlock()
	for _, w := range p.pool {
		w.InUse = false
		w.Rec.Free()
	}
}
