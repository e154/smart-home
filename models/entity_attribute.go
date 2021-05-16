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

package models

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/e154/smart-home/common"
	"reflect"
	"time"
)

type EntityAttribute struct {
	Name  string                     `json:"name"`
	Type  common.EntityAttributeType `json:"type"`
	Value interface{}                `json:"value,omitempty"`
}

func (a EntityAttribute) String() string {
	if value, ok := a.Value.(string); ok {
		return value
	}
	return ""
}

func (a EntityAttribute) Int64() int64 {
	if value, ok := a.Value.(int64); ok {
		return value
	}
	t := reflect.TypeOf(a.Value)
	switch t.Kind() {
	case reflect.Float64:
		return int64(a.Float64())
	}
	return 0
}

func (a EntityAttribute) Time() time.Time {
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

func (a EntityAttribute) Bool() bool {
	if value, ok := a.Value.(bool); ok {
		return value
	}
	return false
}

func (a EntityAttribute) Float64() float64 {
	if value, ok := a.Value.(float64); ok {
		return value
	}
	return 0
}

func (a EntityAttribute) Map() EntityAttributes {
	if value, ok := a.Value.(EntityAttributes); ok {
		return value
	}
	return nil
}

type EntityAttributeValue map[string]interface{}

type EntityAttributes map[string]*EntityAttribute

func (a EntityAttributes) Serialize() (to EntityAttributeValue) {

	var serialize func(from EntityAttributes, to EntityAttributeValue)
	serialize = func(from EntityAttributes, to EntityAttributeValue) {

		for kFrom, vFrom := range from {
			switch vFrom.Type {
			case common.EntityAttributeString:
			case common.EntityAttributeInt:
			case common.EntityAttributeTime:
			case common.EntityAttributeBool:
			case common.EntityAttributeFloat:
			case common.EntityAttributeArray:

				arr := make([]interface{}, 0)

				if attrs, ok := vFrom.Value.([]interface{}); ok {
					for _, attr := range attrs {
						switch t5 := attr.(type) {
						case float64, float32:
						case int64, int32:
						case EntityAttributes:
							t12 := EntityAttributeValue{}
							serialize(t5, t12)
							arr = append(arr, t12)
							continue
						default:
							log.Warnf("unknown type %s", reflect.TypeOf(t5).String())
						}
						arr = append(arr, attr)
					}
				}

				to[kFrom] = arr

				continue
			case common.EntityAttributeMap:
				if t6, ok := vFrom.Value.(EntityAttributes); ok {
					to2 := EntityAttributeValue{}
					serialize(t6, to2)
					to[kFrom] = to2
				}
				continue
			default:
				log.Warnf("unknown type %s", vFrom.Type)
				continue
			}

			to[kFrom] = vFrom.Value
		}
	}

	to = EntityAttributeValue{}
	serialize(a, to)

	return
}

func (a EntityAttributes) Deserialize(data EntityAttributeValue) (changed bool, err error) {

	var deserialize func(from EntityAttributeValue, to EntityAttributes)
	deserialize = func(from EntityAttributeValue, to EntityAttributes) {

		for kFrom, vFrom := range from {
			if vFrom == nil {
				continue
			}
			switch vFromCasted := vFrom.(type) {
			case map[string]interface{}:
				if _, ok := to[kFrom]; ok {
					switch value := to[kFrom].Value.(type) {
					case map[string]interface{}:
						attrs := make(EntityAttributes)
						for k, values := range value {
							if val, ok := values.(map[string]interface{}); ok {
								attrs[k] = &EntityAttribute{
									Name:  val["name"].(string),
									Type:  common.EntityAttributeType(val["type"].(string)),
									Value: vFromCasted[k],
								}
							}
						}
						deserialize(vFromCasted, attrs)
						to[kFrom].Value = attrs

					case EntityAttributes:
						deserialize(vFromCasted, value)
					default:
						log.Warnf("unknown type %s (%v)", reflect.TypeOf(to[kFrom].Value).String(), vFrom)
					}
				}
				continue
			case EntityAttributeValue:
				if _, ok := to[kFrom]; ok {
					if attrs, ok := to[kFrom].Value.(EntityAttributes); ok {
						deserialize(vFromCasted, attrs)
					}
				}
				continue
			case string:
			case bool:
			case time.Time:
			case float64, float32:
			case int64, int32, int:
			case []interface{}:

				var arr []interface{}

				for i, t5 := range vFromCasted {

					switch t4 := t5.(type) {
					case int64, int32:
					case float64, float32:
					case bool:
					case time.Time:
					case EntityAttributeValue:
						if t2, ok := to[kFrom].Value.([]EntityAttributes); ok {
							if len(t2) == 0 || len(t2) < i {
								continue
							}
							deserialize(t4, t2[i])
							arr = append(arr, t2[i])
						}
						continue
					case map[string]interface{}:
						if t2, ok := to[kFrom].Value.([]EntityAttributes); ok {
							if len(t2) == 0 || len(t2) < i {
								continue
							}
							deserialize(t4, t2[i])
							arr = append(arr, t2[i])
						}
						continue
					default:
						log.Warnf("unknown type %s", reflect.TypeOf(t5).String())
					}
					arr = append(arr, t5)

				}

				t3 := to[kFrom]
				t3.Value = arr
				to[kFrom] = t3

				continue
			default:
				log.Warnf("unknown type %s (%v)", reflect.TypeOf(vFrom).String(), vFrom)
			}

			if v, ok := to[kFrom]; ok {
				to[kFrom] = &EntityAttribute{
					Name:  v.Name,
					Type:  v.Type,
					Value: vFrom,
				}

				if fmt.Sprintf("%v", vFrom) != fmt.Sprintf("%v", v.Value) {
					changed = true
				}
			}
		}
	}

	deserialize(data, a)

	return
}

func (a EntityAttributes) Signature() (signature EntityAttributes) {

	var serialize func(from, to EntityAttributes)
	serialize = func(from, to EntityAttributes) {

		for kFrom, vFrom := range from {
			switch vFrom.Type {
			case common.EntityAttributeString:
			case common.EntityAttributeInt:
			case common.EntityAttributeTime:
			case common.EntityAttributeBool:
			case common.EntityAttributeFloat:
			case common.EntityAttributeArray:

				if attrs, ok := vFrom.Value.([]interface{}); ok {
					arr := make([]interface{}, 0)
					for _, attr := range attrs {
						switch t5 := attr.(type) {
						case EntityAttributes:
							t12 := EntityAttributes{}
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
			case common.EntityAttributeMap:
				if t6, ok := vFrom.Value.(EntityAttributes); ok {
					to2 := EntityAttributes{}
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

	signature = make(EntityAttributes)
	cpy := a.Copy()
	serialize(cpy, signature)

	return
}

func (e EntityAttributes) Copy() (copy EntityAttributes) {

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	if err := enc.Encode(e); err != nil {
		log.Error(err.Error())
		return
	}

	copy = make(EntityAttributes)
	if err := dec.Decode(&copy); err != nil {
		log.Error(err.Error())
	}
	return
}

func init() {
	gob.Register(time.Time{})
	gob.Register(EntityAttributes{})
}
