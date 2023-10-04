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

package endpoint

import (
	"bufio"
	"context"
	"mime/multipart"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
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
func (i *ImageEndpoint) Add(ctx context.Context, params *m.Image) (image *m.Image, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = i.validation.Valid(params); !ok {
		return
	}

	var id int64
	if id, err = i.adaptors.Image.Add(ctx, params); err != nil {
		return
	}

	image, err = i.adaptors.Image.GetById(ctx, id)

	return
}

// GetById ...
func (i *ImageEndpoint) GetById(ctx context.Context, imageId int64) (image *m.Image, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = i.validation.ValidVar(imageId, "id", "required,numeric"); !ok {
		return
	}

	image, err = i.adaptors.Image.GetById(ctx, imageId)

	return
}

// Update ...
func (i *ImageEndpoint) Update(ctx context.Context, params *m.Image) (result *m.Image, errs validator.ValidationErrorsTranslations, err error) {

	var image *m.Image
	if image, err = i.adaptors.Image.GetById(ctx, params.Id); err != nil {
		return
	}

	if err = copier.Copy(&image, &params); err != nil {
		return
	}

	var ok bool
	if ok, errs = i.validation.Valid(params); !ok {
		return
	}

	if err = i.adaptors.Image.Update(ctx, image); err != nil {
		return
	}

	result, err = i.adaptors.Image.GetById(ctx, params.Id)

	return
}

// Delete ...
func (i *ImageEndpoint) Delete(ctx context.Context, imageId int64) (errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = i.validation.ValidVar(imageId, "id", "required,numeric"); !ok {
		return
	}

	var image *m.Image
	if image, err = i.adaptors.Image.GetById(ctx, imageId); err != nil {
		return
	}

	err = i.adaptors.Image.Delete(ctx, image.Id)

	return
}

// Upload ...
func (i *ImageEndpoint) Upload(ctx context.Context, files map[string][]*multipart.FileHeader) (fileList []*m.Image, errs []error) {

	fileList = make([]*m.Image, 0)
	errs = make([]error, 0)

	for _, fileHeader := range files {

		file, err := fileHeader[0].Open()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		reader := bufio.NewReader(file)
		var newImage *m.Image
		if newImage, err = i.adaptors.Image.UploadImage(ctx, reader, fileHeader[0].Filename); err != nil {
			errs = append(errs, err)
		} else {
			fileList = append(fileList, newImage)
		}

		file.Close()
	}

	return
}

// GetList ...
func (i *ImageEndpoint) GetList(ctx context.Context, pagination common.PageParams) (items []*m.Image, total int64, err error) {

	items, total, err = i.adaptors.Image.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)

	return
}

// GetListByDate ...
func (i *ImageEndpoint) GetListByDate(ctx context.Context, filter string) (images []*m.Image, err error) {

	images, err = i.adaptors.Image.GetAllByDate(ctx, filter)

	return
}

// GetFilterList ...
func (i *ImageEndpoint) GetFilterList(ctx context.Context) (filterList []*m.ImageFilterList, err error) {

	filterList, err = i.adaptors.Image.GetFilterList(ctx)

	return
}
