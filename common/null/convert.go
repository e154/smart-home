package null

import (
	"errors"
	"strconv"
	"reflect"
	"fmt"
)

type Value interface{}

var errNilPtr = errors.New("destination pointer is nil")

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
		return errors.New("destination not a pointer")
	}
	if dpv.IsNil() {
		return errNilPtr
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
