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

package models

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/debug"
	"github.com/e154/smart-home/common/encryptor"
)

// Attribute ...
type Attribute struct {
	Name  string               `json:"name"`
	Type  common.AttributeType `json:"type"`
	Value interface{}          `json:"value,omitempty"`
}

// String ...
func (a Attribute) String() string {
	if a.Value == nil {
		return ""
	}
	if value, ok := a.Value.(string); ok {
		return value
	}
	return fmt.Sprintf("%v", a.Value)
}

// Decrypt ...
func (a Attribute) Decrypt() string {
	if a.Value == nil {
		return ""
	}
	if value, ok := a.Value.(string); ok {
		if str, err := encryptor.Decrypt(value); err == nil {
			return str
		}
	}
	return fmt.Sprintf("%v", a.Value)
}

// Int64 ...
func (a Attribute) Int64() int64 {
	if a.Value == nil {
		return 0
	}
	if value, ok := a.Value.(int64); ok {
		return value
	}
	t := reflect.TypeOf(a.Value)
	switch t.Kind() {
	case reflect.Uint64:
		return int64(a.Value.(uint64))
	case reflect.Float32:
		return int64(a.Value.(float32))
	case reflect.Int:
		return int64(a.Value.(int))
	case reflect.Float64:
		return int64(a.Float64())
	case reflect.String:
		i, _ := strconv.Atoi(a.Value.(string))
		return int64(i)
	default:
		log.Warnf("unknown type %s", t.String())
	}
	return 0
}

// Time ...
func (a *Attribute) Time() time.Time {
	if a == nil || a.Value == nil {
		return time.Time{}
	}
	if value, ok := a.Value.(time.Time); ok {
		return value
	}
	if value, ok := a.Value.(string); ok {
		if t, err := time.Parse(time.RFC3339, value); err == nil {
			return t
		}
	}
	return time.Time{}
}

// Bool ...
func (a *Attribute) Bool() bool {
	if a == nil || a.Value == nil {
		return false
	}
	if value, ok := a.Value.(bool); ok {
		return value
	}
	return false
}

// Float64 ...
func (a *Attribute) Float64() float64 {
	if a == nil || a.Value == nil {
		return 0
	}
	if value, ok := a.Value.(float64); ok {
		return value
	}
	return float64(a.Int64())
}

// Map ...
func (a Attribute) Map() Attributes {
	if value, ok := a.Value.(Attributes); ok {
		return value
	}
	return nil
}

// ArrayString ...
func (a Attribute) ArrayString() (result []string) {
	if a.Value == nil {
		return
	}
	switch v := a.Value.(type) {
	case []interface{}:
		for _, val := range a.Value.([]interface{}) {
			result = append(result, fmt.Sprintf("%v", val))
		}
	case []string:
		result = v
	}
	return
}

// Point ...
func (a Attribute) Point() (point Point) {
	if a.Value == nil {
		return
	}
	if value, ok := a.Value.([]interface{}); ok {
		point.Lat, _ = strconv.ParseFloat(fmt.Sprintf("%v", value[1]), 64)
		point.Lon, _ = strconv.ParseFloat(fmt.Sprintf("%v", value[0]), 64)
		return
	}
	return
}

// AttributeValue ...
type AttributeValue map[string]interface{}

// Attributes ...
type Attributes map[string]*Attribute

// Serialize ...
func (a Attributes) Serialize() (to AttributeValue) {

	var serialize func(from Attributes, to AttributeValue)
	serialize = func(from Attributes, to AttributeValue) {

		for keyFrom, valueFromRaw := range from {
			switch valueFromRaw.Type {
			case common.AttributeString:
			case common.AttributeInt:
			case common.AttributeTime:
			case common.AttributeBool:
			case common.AttributeFloat:
			case common.AttributeImage:
			case common.AttributeIcon:
			case common.AttributePoint:
			case common.AttributeEncrypted:
			case common.AttributeArray:

				arr := make([]interface{}, 0)

				if valueRawList, ok := valueFromRaw.Value.([]interface{}); ok {
					for _, valueRaw := range valueRawList {
						switch value := valueRaw.(type) {
						case float64, float32:
						case int64, int32:
						case string:
						case Attributes:
							attr := AttributeValue{}
							serialize(value, attr)
							arr = append(arr, attr)
							continue
						default:
							log.Warnf("unknown type %s", reflect.TypeOf(value).String())
						}
						arr = append(arr, valueRaw)
					}
				}

				to[keyFrom] = arr

				continue
			case common.AttributeMap:
				if attrs, ok := valueFromRaw.Value.(Attributes); ok {
					attr := AttributeValue{}
					serialize(attrs, attr)
					to[keyFrom] = attr
				}
				continue
			default:
				log.Warnf("unknown type %s", valueFromRaw.Type)
				continue
			}

			to[keyFrom] = valueFromRaw.Value
		}
	}

	to = AttributeValue{}
	serialize(a, to)

	return
}

// Deserialize ...
func (a Attributes) Deserialize(data AttributeValue) (changed bool, err error) {

	var deserialize func(from AttributeValue, to Attributes)
	deserialize = func(from AttributeValue, to Attributes) {

		for keyFrom, valueFromRaw := range from {
			if valueFromRaw == nil {
				continue
			}
			switch valueFrom := valueFromRaw.(type) {
			case map[string]interface{}:
				if _, ok := to[keyFrom]; ok {
					switch value := to[keyFrom].Value.(type) {
					case map[string]interface{}:
						attrs := make(Attributes)
						for k, values := range value {
							if val, ok := values.(map[string]interface{}); ok {
								attrs[k] = &Attribute{
									Name:  val["name"].(string),
									Type:  common.AttributeType(val["type"].(string)),
									Value: valueFrom[k],
								}
							}
						}
						deserialize(valueFrom, attrs)
						to[keyFrom].Value = attrs

					case Attributes:
						deserialize(valueFrom, value)
					default:
						log.Warnf("unknown type %s (%v)", reflect.TypeOf(to[keyFrom].Value).String(), valueFromRaw)
					}
				}
				continue
			case AttributeValue:
				if _, ok := to[keyFrom]; ok {
					if attrs, ok := to[keyFrom].Value.(Attributes); ok {
						deserialize(valueFrom, attrs)
					}
				}
				continue
			case string:
			case bool:
			case time.Time:
			case float64, float32:
			case int64, int32, int:
			case []string:
			case []interface{}:

				var arr []interface{}

				for i, valueRaw := range valueFrom {

					switch value := valueRaw.(type) {
					case int64, int32:
					case float64, float32:
					case bool:
					case string:
					case time.Time:
					case AttributeValue:
						if items, ok := to[keyFrom].Value.([]Attributes); ok {
							if len(items) == 0 || len(items) < i {
								continue
							}
							deserialize(value, items[i])
							arr = append(arr, items[i])
						}
						continue
					case map[string]interface{}:
						if items, ok := to[keyFrom].Value.([]Attributes); ok {
							if len(items) == 0 || len(items) < i {
								continue
							}
							deserialize(value, items[i])
							arr = append(arr, items[i])
						}
						continue
					default:
						log.Warnf("unknown type %s", reflect.TypeOf(valueRaw).String())
					}
					arr = append(arr, valueRaw)

				}

				value := to[keyFrom]
				value.Value = arr
				to[keyFrom] = value

				continue
			default:
				log.Warnf("unknown type %s (%v)", reflect.TypeOf(valueFromRaw).String(), valueFromRaw)
			}

			if v, ok := to[keyFrom]; ok {
				to[keyFrom] = &Attribute{
					Name:  v.Name,
					Type:  v.Type,
					Value: valueFromRaw,
				}

				if fmt.Sprintf("%v", valueFromRaw) != fmt.Sprintf("%v", v.Value) {
					changed = true
				}
			}
		}
	}

	deserialize(data, a)

	return
}

// Signature ...
func (a Attributes) Signature() (signature Attributes) {

	var serialize func(from, to Attributes)
	serialize = func(from, to Attributes) {

		for kFrom, vFrom := range from {
			switch vFrom.Type {
			case common.AttributeString:
			case common.AttributeInt:
			case common.AttributeTime:
			case common.AttributeBool:
			case common.AttributeFloat:
			case common.AttributeImage:
			case common.AttributeIcon:
			case common.AttributePoint:
			case common.AttributeEncrypted:
			case common.AttributeArray:

				if attrs, ok := vFrom.Value.([]interface{}); ok {
					arr := make([]interface{}, 0)
					for _, attr := range attrs {
						switch t5 := attr.(type) {
						case Attributes:
							t12 := Attributes{}
							serialize(t5, t12)
							arr = append(arr, t12)
						case float64, float32:
						case int64, int32:
						default:
							log.Warnf("unknown type %s", reflect.TypeOf(t5).String())
							continue
						}
					}
					v := from[kFrom]
					v.Value = arr
					to[kFrom] = v
				}

				continue
			case common.AttributeMap:
				if t6, ok := vFrom.Value.(Attributes); ok {
					to2 := Attributes{}
					serialize(t6, to2)
					v := from[kFrom]
					v.Value = to2
					to[kFrom] = v
				}
				continue
			default:
				log.Warnf("unknown type %s", vFrom.Type)
				continue
			}

			v := from[kFrom]
			v.Value = nil
			to[kFrom] = v
		}
	}

	signature = make(Attributes)
	cpy := a.Copy()
	serialize(cpy, signature)

	return
}

// Copy ...
func (a Attributes) Copy() (copy Attributes) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	var err error
	defer func() {
		if err != nil {
			log.Info("============= object with error =============")
			debug.Println(a)
		}
	}()

	if err = enc.Encode(a); err != nil {
		log.Error(err.Error())
		return
	}

	copy = make(Attributes)
	if err = dec.Decode(&copy); err != nil {
		log.Error(err.Error())
	}
	return
}

// Compare ...
func (a Attribute) Compare(b *Attribute) (ident bool) {

	switch a.Type {
	case common.AttributeString:
	case common.AttributeInt:
	case common.AttributeTime:
	case common.AttributeBool:
	case common.AttributeFloat:
	case common.AttributeImage:
	case common.AttributeIcon:
	case common.AttributePoint:
	case common.AttributeEncrypted:
	case common.AttributeArray:
		return
	}

	ident = fmt.Sprintf("%v", a.Value) == fmt.Sprintf("%v", b.Value)

	return
}

func init() {
	gob.Register(time.Time{})
	gob.Register(Attributes{})
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
}
