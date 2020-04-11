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

package storage

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/graceful_service"
	"go.uber.org/atomic"
	"strings"
	"sync"
	"time"
)

var (
	log = common.MustGetLogger("storage")
)

// Storage ...
type Storage struct {
	adaptors  *adaptors.Adaptors
	pool      sync.Map
	quit      chan struct{}
	inProcess *atomic.Bool
}

// NewStorage ...
func NewStorage(
	adaptors *adaptors.Adaptors,
	graceful *graceful_service.GracefulService) *Storage {
	storage := &Storage{
		adaptors:  adaptors,
		pool:      sync.Map{},
		quit:      make(chan struct{}),
		inProcess: atomic.NewBool(false),
	}

	graceful.Subscribe(storage)
	
	go func() {
		ticker := time.NewTicker(time.Minute * 1)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				storage.serialize()

			case _, ok := <-storage.quit:
				if !ok {
					return
				}
				close(storage.quit)
				return
			}
		}
	}()

	return storage
}

func (s *Storage) Shutdown() {
	s.quit <- struct{}{}
	s.serialize()
}

// Search ...
func (s *Storage) Search(name string) (result map[string][]byte) {
	return s.search(name)
}

// Push ...
func (s *Storage) Push(name string, v string) (err error) {
	err = s.push(name, v)
	return
}

// GetByName ...
func (s *Storage) GetByName(name string) (val []byte, err error) {
	return s.getByName(name)
}

// Pop ...
func (s *Storage) Pop(name string) (val []byte, err error) {
	return s.pop(name)
}

func (s *Storage) push(name string, v string) (err error) {
	s.pool.Store(name, m.Storage{
		Name:    name,
		Changed: true,
		Value:   []byte(v),
	})
	return
}

func (s *Storage) getByName(name string) (val []byte, err error) {

	if v, ok := s.pool.Load(name); ok {
		val = v.(m.Storage).Value
		return
	}
	var storage m.Storage
	if storage, err = s.adaptors.Storage.GetByName(name); err != nil {
		return
	}
	val = storage.Value

	return
}

func (s *Storage) pop(name string) (val []byte, err error) {
	val, err = s.getByName(name)
	if err != nil {
		return
	}
	if err = s.adaptors.Storage.Delete(name); err != nil {
		return
	}
	s.pool.Delete(name)
	return
}

func (s *Storage) Serialize() {
	s.serialize()
}

func (s *Storage) serialize() {

	if s.inProcess.Load() {
		return
	}
	s.inProcess.Store(true)

	var data m.Storage
	var ok bool

	s.pool.Range(func(key, val interface{}) bool {
		data, ok = val.(m.Storage)
		if !ok {
			return true
		}

		if !data.Changed {
			return true
		}

		s.pool.Store(key, data)

		if err := s.adaptors.Storage.CreateOrUpdate(data); err != nil {
			log.Error(err.Error())
			return true
		}

		return true
	})

}

func (s *Storage) search(sub string) (result map[string][]byte) {
	result = make(map[string][]byte)
	s.pool.Range(func(key, val interface{}) bool {
		if strings.Contains(key.(string), sub) {
			if data, ok := val.(m.Storage); ok {
				result[data.Name] = data.Value
			}
		}

		return true
	})

	list, _, err := s.adaptors.Storage.Search(sub, 99, 0)
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
