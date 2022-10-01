// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package bind

import (
	"github.com/e154/smart-home/system/storage"
)

// StorageBind ...
type StorageBind struct {
	storage *storage.Storage
}

// NewStorageBind ...
func NewStorageBind(
	storage *storage.Storage) *StorageBind {
	return &StorageBind{
		storage: storage,
	}
}

// Search ...
func (s *StorageBind) Search(name string) (result map[string]string) {
	result = make(map[string]string)
	storRes := s.storage.Search(name)
	for k, v := range storRes {
		result[k] = v
	}
	return
}

// Push ...
func (s *StorageBind) Push(name string, val string) (err error) {
	return s.storage.Push(name, val)
}

// GetByName ...
func (s *StorageBind) GetByName(name string) string {
	b, _ := s.storage.GetByName(name)
	return b
}

// Pop ...
func (s *StorageBind) Pop(name string) string {
	b, _ := s.storage.Pop(name)
	return b
}
