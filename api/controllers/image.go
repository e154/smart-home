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

package controllers

import (
	"context"
	"net/http"

	"github.com/e154/smart-home/api/stub/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ControllerImage ...
type ControllerImage struct {
	*ControllerCommon
}

// NewControllerImage ...
func NewControllerImage(common *ControllerCommon) ControllerImage {
	return ControllerImage{
		ControllerCommon: common,
	}
}

// AddImage ...
func (c ControllerImage) AddImage(ctx context.Context, req *api.NewImageRequest) (*api.Image, error) {

	image, errs, err := c.endpoint.Image.Add(ctx, c.dto.Image.FromNewImageRequest(req))
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Image.ToImage(image), nil
}

// GetImageById ...
func (c ControllerImage) GetImageById(ctx context.Context, req *api.GetImageRequest) (*api.Image, error) {

	image, errs, err := c.endpoint.Image.GetById(ctx, int64(req.Id))
	if err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Image.ToImage(image), nil
}

// UpdateImageById ...
func (c ControllerImage) UpdateImageById(ctx context.Context, req *api.UpdateImageRequest) (*api.Image, error) {

	image, errs, err := c.endpoint.Image.Update(ctx, c.dto.Image.FromUpdateImageRequest(req))
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.Image.ToImage(image), nil
}

// GetImageList ...
func (c ControllerImage) GetImageList(ctx context.Context, req *api.GetImageListRequest) (*api.GetImageListResult, error) {

	pagination := c.Pagination(req.Limit, req.Offset, req.Order, req.SortBy)
	items, total, err := c.endpoint.Image.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.Image.ToImageListResult(items, uint64(total), pagination), nil
}

// DeleteImageById ...
func (c ControllerImage) DeleteImageById(ctx context.Context, req *api.DeleteImageRequest) (*emptypb.Empty, error) {

	if errs, err := c.endpoint.Image.Delete(ctx, int64(req.Id)); err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return &emptypb.Empty{}, nil
}

// UploadImage ...
func (c ControllerImage) UploadImage(ctx context.Context, req *api.UploadImageRequest) (*api.Image, error) {

	return nil, nil
}

// MuxUploadImage ...
func (c ControllerImage) MuxUploadImage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseMultipartForm(8 << 20); err != nil {
			log.Error(err.Error())
		}

		form := r.MultipartForm
		if len(form.File) == 0 {
			c.writeErr(403, "bad request", w)
			return
		}

		images, errs := c.endpoint.Image.Upload(r.Context(), form.File)

		var resultImages = make([]interface{}, 0)

		for _, img := range images {
			resultImages = append(resultImages, &map[string]int64{
				"id": img.Id,
			})
		}

		c.writeJson(w, &map[string]interface{}{
			"images": resultImages,
			"errors": errs,
		})
	}
}
