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

package controllers

import (
	"github.com/e154/smart-home/internal/api/stub"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/common"
	"github.com/labstack/echo/v4"
)

// ControllerImage ...
type ControllerImage struct {
	*ControllerCommon
}

// NewControllerImage ...
func NewControllerImage(common *ControllerCommon) *ControllerImage {
	return &ControllerImage{
		ControllerCommon: common,
	}
}

// AddImage ...
func (c ControllerImage) ImageServiceAddImage(ctx echo.Context, _ stub.ImageServiceAddImageParams) error {

	obj := &stub.ApiNewImageRequest{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	image, err := c.endpoint.Image.Add(ctx.Request().Context(), c.dto.Image.FromNewImageRequest(obj))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP201(ctx, ResponseWithObj(ctx, c.dto.Image.ToImage(image)))
}

// GetImageById ...
func (c ControllerImage) ImageServiceGetImageById(ctx echo.Context, id int64) error {

	image, err := c.endpoint.Image.GetById(ctx.Request().Context(), id)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Image.ToImage(image)))
}

// UpdateImageById ...
func (c ControllerImage) ImageServiceUpdateImageById(ctx echo.Context, id int64, _ stub.ImageServiceUpdateImageByIdParams) error {

	obj := &stub.ImageServiceUpdateImageByIdJSONBody{}
	if err := c.Body(ctx, obj); err != nil {
		return c.ERROR(ctx, err)
	}

	image, err := c.endpoint.Image.Update(ctx.Request().Context(), c.dto.Image.FromUpdateImageRequest(obj, id))
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, c.dto.Image.ToImage(image)))
}

// GetImageList ...
func (c ControllerImage) ImageServiceGetImageList(ctx echo.Context, params stub.ImageServiceGetImageListParams) error {

	pagination := c.Pagination(params.Page, params.Limit, params.Sort)
	items, total, err := c.endpoint.Image.GetList(ctx.Request().Context(), pagination)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithList(ctx, c.dto.Image.ToImageListResult(items), total, pagination))
}

// DeleteImageById ...
func (c ControllerImage) ImageServiceDeleteImageById(ctx echo.Context, id int64) error {

	if err := c.endpoint.Image.Delete(ctx.Request().Context(), id); err != nil {
		return c.ERROR(ctx, err)
	}

	return c.HTTP200(ctx, ResponseWithObj(ctx, struct{}{}))
}

// MuxUploadImage ...
func (c ControllerImage) ImageServiceUploadImage(ctx echo.Context, _ stub.ImageServiceUploadImageParams) error {

	r := ctx.Request()

	if err := r.ParseMultipartForm(8 << 20); err != nil {
		log.Error(err.Error())
	}

	form := r.MultipartForm
	if len(form.File) == 0 {
		return c.ERROR(ctx, apperr.ErrInvalidRequest)
	}

	images, errs := c.endpoint.Image.Upload(r.Context(), form.File)

	var resultImages = make([]interface{}, 0)

	for _, img := range images {
		resultImages = append(resultImages, map[string]int64{
			"id": img.Id,
		})
	}

	return c.HTTP200(ctx, map[string]interface{}{
		"images": resultImages,
		"errors": errs,
	})
}

// GetImageListByDate ...
func (c ControllerImage) ImageServiceGetImageListByDate(ctx echo.Context, request stub.ImageServiceGetImageListByDateParams) error {
	images, err := c.endpoint.Image.GetListByDate(ctx.Request().Context(), common.StringValue(request.Filter))
	if err != nil {
		return c.ERROR(ctx, err)
	}
	return c.HTTP200(ctx, c.dto.Image.ToImageList(images))
}

// GetImageFilterList ...
func (c ControllerImage) ImageServiceGetImageFilterList(ctx echo.Context) error {
	filters, err := c.endpoint.Image.GetFilterList(ctx.Request().Context())
	if err != nil {
		return c.ERROR(ctx, err)
	}
	return c.HTTP200(ctx, c.dto.Image.ToFilterList(filters))
}
