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

// Code generaTed by fileb0x at "2017-11-26 17:57:18.489607614 +0600 +06 m=+4.245926076" from config file "b0x.yml" DO NOT EDIT.

package swaggerFiles

import (
	"log"
	"os"
)

// FileSwaggerUICSSMap is "/swagger-ui.css.map"
var FileSwaggerUICSSMap = []byte("\x7b\x22\x76\x65\x72\x73\x69\x6f\x6e\x22\x3a\x33\x2c\x22\x73\x6f\x75\x72\x63\x65\x73\x22\x3a\x5b\x5d\x2c\x22\x6e\x61\x6d\x65\x73\x22\x3a\x5b\x5d\x2c\x22\x6d\x61\x70\x70\x69\x6e\x67\x73\x22\x3a\x22\x22\x2c\x22\x66\x69\x6c\x65\x22\x3a\x22\x73\x77\x61\x67\x67\x65\x72\x2d\x75\x69\x2e\x63\x73\x73\x22\x2c\x22\x73\x6f\x75\x72\x63\x65\x52\x6f\x6f\x74\x22\x3a\x22\x22\x7d")

func init() {

	f, err := FS.OpenFile(CTX, "/swagger-ui.css.map", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(FileSwaggerUICSSMap)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}
