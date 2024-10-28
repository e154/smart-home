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

package mqtt

import (
	"sync"
)

// Storage ...
type Storage struct {
	mx   *sync.Mutex
	pull map[string]interface{}
}

// NewStorage ...
func NewStorage() Storage {
	return Storage{
		mx:   &sync.Mutex{},
		pull: make(map[string]interface{}),
	}
}

// GetVar ...
func (s *Storage) GetVar(key string) (value interface{}) {

	s.mx.Lock()
	if v, ok := s.pull[key]; ok {
		value = v
	} else {
		value = nil
	}
	s.mx.Unlock()
	return
}

// SetVar ...
func (s *Storage) SetVar(key string, value interface{}) {

	s.mx.Lock()
	s.pull[key] = value
	s.mx.Unlock()
}

func (s *Storage) copy(newPull map[string]interface{}) {
	s.mx.Lock()
	for key := range s.pull {
		delete(s.pull, key)
	}
	for k, v := range newPull {
		s.pull[k] = v
	}
	s.mx.Unlock()
}
