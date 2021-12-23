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
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Value ...
type Value interface{}

var (
	ErrDestinationPointerIsNil = errors.New("destination pointer is nil")
	ErrDestinationNotAPointer  = errors.New("destination not a pointer")
)

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	} else {
		c := make([]byte, len(b))
		copy(c, b)
		return c
	}
}

func convertAssign(dest, src interface{}) error {

	//switch v := src.(type) {
	//case string:
	//	if v == "null" {
	//		dest = nil
	//	}
	//}

	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		return ErrDestinationNotAPointer
	}
	if dpv.IsNil() {
		return ErrDestinationPointerIsNil
	}

	var sv reflect.Value

	dv := reflect.Indirect(dpv)
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		switch b := src.(type) {
		case []byte:
			dv.Set(reflect.ValueOf(cloneBytes(b)))
		default:
			dv.Set(sv)
		}
		return nil
	}

	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	switch dv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv := reflect.ValueOf(src)
		s := strconv.FormatInt(rv.Int(), 10)
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetInt(i64)
		return nil
	}

	return nil
}
