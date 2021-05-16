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
)

// ErrNotFound represents a not found error.
var ErrNotFound = errors.New("not found")

type quickList struct {
	index map[string]*list.Element
	rows  *list.List
}

func newQuickList() *quickList {
	return &quickList{
		index: make(map[string]*list.Element),
		rows:  list.New(),
	}
}

func (q *quickList) set(id string, value interface{}) {
	if e, ok := q.index[id]; ok {
		e.Value = value
	} else {
		elem := q.rows.PushBack(value)
		q.index[id] = elem
	}
}
func (q *quickList) remove(id string) *list.Element {
	elem := q.index[id]
	if elem != nil {
		q.rows.Remove(elem)
	}
	delete(q.index, id)
	return elem
}
func (q *quickList) getByID(id string) (*list.Element, error) {
	if i, ok := q.index[id]; ok {
		return i, nil
	}
	return nil, ErrNotFound
}
func (q *quickList) iterate(fn func(elem *list.Element), offset, n int) error {
	if offset < 0 || n < 0 {
		return errors.New("invalid offset or n")
	}
	if q.rows.Len() <= offset {
		return errors.New("invalid offset")
	}
	var i int
	for e := q.rows.Front(); e != nil; e = e.Next() {
		if i >= offset && i < offset+n {
			fn(e)
		}
		if i == offset+n {
			break
		}
		i++
	}
	return nil
}
