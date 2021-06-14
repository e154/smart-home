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

package endpoint

import (
	"bufio"
	"errors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"github.com/jinzhu/copier"
	"mime/multipart"
)

// ImageEndpoint ...
type ImageEndpoint struct {
	*CommonEndpoint
}

// NewImageEndpoint ...
func NewImageEndpoint(common *CommonEndpoint) *ImageEndpoint {
	return &ImageEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (i *ImageEndpoint) Add(params *m.Image) (image *m.Image, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = i.adaptors.Image.Add(params); err != nil {
		return
	}

	image, err = i.adaptors.Image.GetById(id)

	return
}

// GetById ...
func (i *ImageEndpoint) GetById(id int64) (image *m.Image, err error) {

	image, err = i.adaptors.Image.GetById(id)

	return
}

// Update ...
func (i *ImageEndpoint) Update(params *m.Image) (result *m.Image, errs []*validation.Error, err error) {

	var image *m.Image
	if image, err = i.adaptors.Image.GetById(params.Id); err != nil {
		return
	}

	if err = copier.Copy(&image, &params); err != nil {
		return
	}

	_, errs = image.Valid()
	if len(errs) > 0 {
		return
	}

	if err = i.adaptors.Image.Update(image); err != nil {
		return
	}

	image, err = i.adaptors.Image.GetById(params.Id)

	return
}

// Delete ...
func (i *ImageEndpoint) Delete(imageId int64) (err error) {

	if imageId == 0 {
		err = errors.New("image id is null")
		return
	}

	var image *m.Image
	if image, err = i.adaptors.Image.GetById(imageId); err != nil {
		return
	}

	err = i.adaptors.Image.Delete(image.Id)

	return
}

// Upload ...
func (i *ImageEndpoint) Upload(files map[string][]*multipart.FileHeader) (fileList []*m.Image, errs []error) {

	fileList = make([]*m.Image, 0)
	errs = make([]error, 0)

	for _, fileHeader := range files {

		file, err := fileHeader[0].Open()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		reader := bufio.NewReader(file)
		if err = i.adaptors.Image.UploadImage(reader, fileHeader[0].Filename); err != nil {
			errs = append(errs, err)
		}

		file.Close()
	}

	return
}

// GetList ...
func (i *ImageEndpoint) GetList(limit, offset int64, order, sortBy string) (items []*m.Image, total int64, err error) {

	items, total, err = i.adaptors.Image.List(limit, offset, order, sortBy)

	return
}
