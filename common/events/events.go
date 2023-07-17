package events

import (
	"reflect"

	"github.com/iancoleman/strcase"
)

func EventName(event interface{}) string {
	if t := reflect.TypeOf(event); t.Kind() == reflect.Ptr {
		return strcase.ToSnake(t.Elem().Name())
	} else {
		return strcase.ToSnake(t.Name())
	}
}
