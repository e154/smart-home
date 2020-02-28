// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"io"
	"log"
	"os"
	"path/filepath"
	"unicode/utf8"
)

const (
	dataDir         = "data"
	fileStoragePath = "file_storage"
	staticPath      = "static"
	depth           = 3
)

func GetFileSize(name string) (int64, error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

func GetFullPath(name string) string {

	dir := filepath.Join(dataDir, fileStoragePath)

	for i := 0; i < depth; i++ {
		dir = filepath.Join(dir, name[i*3:(i+1)*3])
	}

	return dir
}

func GetLinkPath(name string) string {

	dir := "/upload"

	count := utf8.RuneCountInString(name)
	if count < 9 {
		return filepath.Join(dir, name)
	}

	for i := 0; i < depth; i++ {
		dir = filepath.Join(dir, name[i*3:(i+1)*3])
	}

	return filepath.Join(dir, name)
}

func StoragePath() string {
	return filepath.Join(dataDir, fileStoragePath)
}

func StaticPath() string {
	return filepath.Join(dataDir, staticPath)
}

func FileExist(path string) (exist bool) {

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
		} else {
			// other error
		}
		return
	}

	exist = true

	return
}

func CopyFile(f, t string) {

	from, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(t, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}
