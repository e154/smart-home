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

package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"

	"github.com/e154/smart-home/models"
)

// ReadConfig ...
func ReadConfig(dir, fileName, pref string) (conf *models.AppConfig, _ error) {

	conf = &models.AppConfig{}
	file, err := os.ReadFile(path.Join(dir, fileName))
	if err != nil {
		log.Println("Error reading config file")
	} else {
		err = json.Unmarshal(file, conf)
		if err != nil {
			log.Println("Error: wrong format of config file")
		}
	}
	checkEnv(pref, conf)

	return
}

func checkEnv(pref string, v *models.AppConfig) {

	r := reflect.ValueOf(v)
	t := reflect.TypeOf(*v)

	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)
		f := reflect.Indirect(r).FieldByName(field.Name)
		fieldType := t.Field(i).Type

		if !f.CanSet() || field.Tag.Get("env") == "-" {
			continue
		}

		fieldName := field.Tag.Get("env")
		if pref != "" {
			fieldName = fmt.Sprintf("%s_%s", pref, fieldName)
		}

		switch fieldType.String() {
		case "string", "common.RunMode":
			if val := os.Getenv(fieldName); val != "" {
				f.SetString(val)
			}
		case "bool":
			if val := os.Getenv(fieldName); val != "" {
				b, _ := strconv.ParseBool(val)
				f.SetBool(b)
			}
		case "int":
			if val := os.Getenv(fieldName); val != "" {
				i, _ := strconv.ParseInt(val, 10, 32)
				f.SetInt(i)
			}
		case "time.Duration":
			if val := os.Getenv(fieldName); val != "" {
				i, _ := strconv.ParseInt(val, 10, 32)
				f.SetInt(i)
			}
		default:
			log.Fatalf("unknown field type %s\n", fieldType.String())
		}
	}
}
