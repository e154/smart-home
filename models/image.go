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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/e154/smart-home/common"
)

// Image ...
type Image struct {
	Id        int64     `json:"id"`
	Thumb     string    `json:"thumb"`
	Url       string    `json:"url"`
	Image     string    `json:"image"`
	MimeType  string    `json:"mime_type"`
	Title     string    `json:"title"`
	Size      int64     `json:"size"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ImageFilterList ...
type ImageFilterList struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// UploadImage ...
func UploadImage(ctx context.Context, reader *bufio.Reader, fileName string) (newFile *Image, err error) {

	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 128)

	var count int
	for {
		if count, err = reader.Read(part); err != nil {
			break
		}
		buffer.Write(part[:count])
	}
	if err != io.EOF {
		return
	}

	contentType := http.DetectContentType(buffer.Bytes())
	//log.Infof("Content-type from buffer, %s", contentType)

	//------
	// rename & save
	name := common.Strtomd5(common.RandomString(10))
	ext := strings.ToLower(path.Ext(fileName))
	newname := fmt.Sprintf("%s%s", name, ext)

	//create destination file making sure the path is writeable.
	dir := common.GetFullPath(name)
	_ = os.MkdirAll(dir, os.ModePerm)
	var dst *os.File
	if dst, err = os.Create(filepath.Join(dir, newname)); err != nil {
		return
	}

	defer dst.Close()

	//copy the uploaded file to the destination file
	if _, err = io.Copy(dst, buffer); err != nil {
		return
	}

	size, _ := common.GetFileSize(filepath.Join(dir, newname))
	newFile = &Image{
		Size:     size,
		MimeType: contentType,
		Image:    newname,
		Name:     fileName,
	}

	return
}
