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
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/pkg/apperr"
	m "github.com/e154/smart-home/pkg/models"

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
func (i *ImageEndpoint) Add(ctx context.Context, params *m.Image) (image *m.Image, err error) {

	if ok, errs := i.validation.Valid(params); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	var id int64
	if id, err = i.adaptors.Image.Add(ctx, params); err != nil {
		return
	}

	if image, err = i.adaptors.Image.GetById(ctx, id); err != nil {
		return
	}

	log.Infof("added new image id:(%d)", image.Id)

	return
}

// GetById ...
func (i *ImageEndpoint) GetById(ctx context.Context, imageId int64) (image *m.Image, err error) {

	if ok, errs := i.validation.ValidVar(imageId, "id", "required,numeric"); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	image, err = i.adaptors.Image.GetById(ctx, imageId)

	return
}

// Update ...
func (i *ImageEndpoint) Update(ctx context.Context, params *m.Image) (result *m.Image, err error) {

	var image *m.Image
	if image, err = i.adaptors.Image.GetById(ctx, params.Id); err != nil {
		return
	}

	if err = copier.Copy(&image, &params); err != nil {
		return
	}

	if ok, errs := i.validation.Valid(params); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = i.adaptors.Image.Update(ctx, image); err != nil {
		return
	}

	if result, err = i.adaptors.Image.GetById(ctx, params.Id); err != nil {
		return
	}

	log.Infof("updated image id:(%d)", result.Id)

	return
}

// Delete ...
func (i *ImageEndpoint) Delete(ctx context.Context, imageId int64) (err error) {

	if ok, errs := i.validation.ValidVar(imageId, "id", "required,numeric"); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	var image *m.Image
	if image, err = i.adaptors.Image.GetById(ctx, imageId); err != nil {
		return
	}

	if err = i.adaptors.Image.Delete(ctx, image.Id); err != nil {
		return
	}

	log.Infof("image id:(%d) was deleted", imageId)

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
		newImage, err = UploadImage(reader, fileHeader[0].Filename)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		newImage.Id, err = i.adaptors.Image.Add(ctx, newImage)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		fileList = append(fileList, newImage)

		file.Close()
	}

	log.Infof("uploaded %d images", len(fileList))

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

func UploadImage(reader *bufio.Reader, fileName string) (newFile *m.Image, err error) {

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
	newFile = &m.Image{
		Size:     size,
		MimeType: contentType,
		Image:    newname,
		Name:     fileName,
	}

	return
}
