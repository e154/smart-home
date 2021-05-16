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

package null

import (
	"database/sql/driver"
	"strconv"
)

// Int64 ...
type Int64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// NewInt64 ...
func NewInt64(value interface{}) (i Int64) {
	i = Int64{}
	i.Scan(value)
	return
}

// Scan ...
func (n *Int64) Scan(value interface{}) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Int64, value)
}

// Value ...
func (n Int64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

// String ...
func (n Int64) String() (value string) {
	if !n.Valid {
		value = "null"
		return
	}
	value = strconv.FormatInt(n.Int64, 10)
	return
}

// MarshalJSON ...
func (n Int64) MarshalJSON() ([]byte, error) {
	return []byte(n.String()), nil
}

// UnmarshalJSON ...
func (n *Int64) UnmarshalJSON(data []byte) error {
	i64, err := strconv.ParseInt(string(data), 10, 0)
	if err != nil {
		return nil
	}
	n.Valid = true
	return n.Scan(i64)
}
