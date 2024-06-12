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

package storage

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/e154/bus"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
)

var (
	log = logger.MustGetLogger("storage")
)

// Storage ...
type Storage struct {
	adaptors  *adaptors.Adaptors
	pool      sync.Map
	quit      chan struct{}
	inProcess *atomic.Bool
	isStarted *atomic.Bool
	eventBus  bus.Bus
}

// NewStorage ...
func NewStorage(
	adaptors *adaptors.Adaptors,
	eventBus bus.Bus) *Storage {
	storage := &Storage{
		adaptors:  adaptors,
		pool:      sync.Map{},
		quit:      make(chan struct{}),
		inProcess: atomic.NewBool(false),
		isStarted: atomic.NewBool(true),
		eventBus:  eventBus,
	}

	go func() {
		ticker := time.NewTicker(time.Minute * 1)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				storage.serialize()
			case <-storage.quit:
				return
			}
		}
	}()

	return storage
}

// Shutdown ...
func (s *Storage) Shutdown() {
	if !s.isStarted.CompareAndSwap(true, false) {
		return
	}
	close(s.quit)
	s.serialize()
}

// Search ...
func (s *Storage) Search(name string) (result map[string]string) {
	return s.search(name)
}

// Push ...
func (s *Storage) Push(name string, v string) (err error) {
	err = s.push(name, v)
	return
}

// GetByName ...
func (s *Storage) GetByName(name string) (val string, err error) {
	return s.getByName(name)
}

// Pop ...
func (s *Storage) Pop(name string) (val string, err error) {
	return s.pop(name)
}

func (s *Storage) push(name string, v string) (err error) {
	s.pool.Store(name, m.Variable{
		Name:    name,
		Changed: true,
		Value:   v,
	})
	return
}

func (s *Storage) getByName(name string) (val string, err error) {

	if v, ok := s.pool.Load(name); ok {
		val = v.(m.Variable).Value
		return
	}
	var storage m.Variable
	if storage, err = s.adaptors.Variable.GetByName(context.Background(), name); err != nil {
		return
	}
	val = storage.Value

	return
}

func (s *Storage) pop(name string) (val string, err error) {
	val, err = s.getByName(name)
	if err != nil {
		return
	}
	if err = s.adaptors.Variable.Delete(context.Background(), name); err != nil {
		return
	}
	s.pool.Delete(name)
	return
}

// Serialize ...
func (s *Storage) Serialize() {
	s.serialize()
}

func (s *Storage) serialize() {

	if !s.inProcess.CompareAndSwap(false, true) {
		return
	}
	defer s.inProcess.Store(false)

	var data m.Variable
	var ok bool

	s.pool.Range(func(key, val interface{}) bool {
		data, ok = val.(m.Variable)
		if !ok {
			return true
		}

		if !data.Changed {
			return true
		}

		data.Changed = false

		s.pool.Store(key, data)

		if err := s.adaptors.Variable.CreateOrUpdate(context.Background(), data); err != nil {
			log.Error(err.Error())
			return true
		}
		s.eventBus.Publish(fmt.Sprintf("system/models/variables/%s", data.Name), events.EventUpdatedVariableModel{
			Name:  data.Name,
			Value: data.Value,
		})

		return true
	})

}

func (s *Storage) search(sub string) (result map[string]string) {
	result = make(map[string]string)
	s.pool.Range(func(key, val interface{}) bool {
		if strings.Contains(key.(string), sub) {
			if data, ok := val.(m.Variable); ok {
				result[data.Name] = data.Value
			}
		}

		return true
	})

	list, _, err := s.adaptors.Variable.Search(context.Background(), sub, 99, 0)
	if err != nil {
		return
	}
	for _, fromDb := range list {
		if _, ok := result[fromDb.Name]; ok {
			continue
		}
		result[fromDb.Name] = fromDb.Value
		s.pool.Store(fromDb.Name, fromDb)
	}
	return
}
