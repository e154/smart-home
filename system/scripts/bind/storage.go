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

package bind

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"time"
)

// Javascript Binding
//
// Storage
// 	 .Search(name)
// 	 .Push(name, value)
// 	 .GetByName(name)
// 	 .Pop(name)
//
type StorageBind struct {
	adaptors *adaptors.Adaptors
}

func NewStorageBind(adaptors *adaptors.Adaptors) *StorageBind {
	return &StorageBind{adaptors: adaptors}
}

func (s *StorageBind) Search(name string) (result map[string]string) {
	result = make(map[string]string)
	list, _, err := s.adaptors.Storage.Search(name, 99, 0)
	if err != nil {
		return
	}
	for _, stor := range list {
		result[stor.Name] = string(stor.Value)
	}
	return
}

func (s *StorageBind) Push(name string, v string) (err error) {
	err = s.adaptors.Storage.CreateOrUpdate(&m.Storage{
		Name:      name,
		Value:     []byte(v),
		UpdatedAt: time.Time{},
		CreatedAt: time.Time{},
	})
	return
}

func (s *StorageBind) GetByName(name string) (value string) {
	storage, err := s.adaptors.Storage.GetByName(name)
	if err != nil {
		return
	}
	v, _ := storage.Value.MarshalJSON()
	value = string(v)
	return
}

func (s *StorageBind) Pop(name string) (value string) {
	value = s.GetByName(name)
	s.adaptors.Storage.Delete(name)
	return
}
