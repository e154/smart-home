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

package admin

import (
	"container/list"
	"errors"
	"github.com/DrmagicE/gmqtt/server"
	"time"
)

// ErrNotFound represents a not found error.
var ErrNotFound = errors.New("not found")

// Indexer provides a index for a ordered list that supports queries in O(1).
// All methods are not concurrency-safe.
type Indexer struct {
	index map[string]*list.Element
	rows  *list.List
}

// NewIndexer is the constructor of Indexer.
func NewIndexer() *Indexer {
	return &Indexer{
		index: make(map[string]*list.Element),
		rows:  list.New(),
	}
}

// Set sets the value for the id.
func (i *Indexer) Set(id string, value interface{}) {
	if e, ok := i.index[id]; ok {
		e.Value = value
	} else {
		elem := i.rows.PushBack(value)
		i.index[id] = elem
	}
}

// Remove removes and returns the value for the given id.
// Return nil if not found.
func (i *Indexer) Remove(id string) *list.Element {
	elem := i.index[id]
	if elem != nil {
		i.rows.Remove(elem)
	}
	delete(i.index, id)
	return elem
}

// GetByID returns the value for the given id.
// Return nil if not found.
// Notice: Any access to the return *list.Element also require the mutex,
// because the Set method can modify the Value for *list.Element when updating the Value for the same id.
// If the caller needs the Value in *list.Element, it must get the Value before the next Set is called.
func (i *Indexer) GetByID(id string) *list.Element {
	return i.index[id]
}

// Iterate iterates at most n elements in the list begin from offset.
// Notice: Any access to the  *list.Element in fn also require the mutex,
// because the Set method can modify the Value for *list.Element when updating the Value for the same id.
// If the caller needs the Value in *list.Element, it must get the Value before the next Set is called.
func (i *Indexer) Iterate(fn func(elem *list.Element), offset, n uint) {
	if i.rows.Len() < int(offset) {
		return
	}
	var j uint
	for e := i.rows.Front(); e != nil; e = e.Next() {
		if j >= offset && j < offset+n {
			fn(e)
		}
		if j == offset+n {
			break
		}
		j++
	}
}

// Len returns the length of list.
func (i *Indexer) Len() int {
	return i.rows.Len()
}

// GetOffsetN ...
func GetOffsetN(page, pageSize uint) (offset, n uint) {
	offset = (page - 1) * pageSize
	n = pageSize
	return
}

const (
	// Online ...
	Online = "online"
	// Offline ...
	Offline = "offline"
)

func statusText(client server.Client) string {
	if client.SessionInfo().IsExpired(time.Now()) {
		return Online
	} else {
		return Offline
	}
}
