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

package common

import (
	"crypto/md5"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"go/build"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

const DefaultPageSize int64 = 15

//create md5 string
func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

const (
	// Alphanum ...
	Alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// Alpha ...
	Alpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// Number ...
	Number = "0123456789"
)

// RandStr ...
func RandStr(strSize int, dictionary string) string {

	var bytes = make([]byte, strSize)
	crypto_rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

// RandInt ...
func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

// RandomString ...
func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 129))
	}
	return string(bytes)
}

func TestMode() bool {
	return os.Getenv("TEST_MODE") == "true"
}

func Dir() string {
	dir, _ := os.Getwd()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	project := path.Join(gopath, "src", "")
	dir = strings.Replace(dir, project, "+", -1)
	dir = strings.Replace(dir, "+/", "", -1)
	return dir
}
