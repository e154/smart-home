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
	"github.com/e154/smart-home/api/stub/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
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
func (c ControllerImage) AddImage(_ context.Context, req *api.NewImageRequest) (*api.Image, error) {

	image, errs, err := c.endpoint.Image.Add(c.dto.Image.FromNewImageRequest(req))
	if len(errs) > 0 {
		return nil, c.prepareErrors(errs)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Image.ToImage(image), nil
}

// GetImageById ...
func (c ControllerImage) GetImageById(_ context.Context, req *api.GetImageRequest) (*api.Image, error) {

	image, err := c.endpoint.Image.GetById(int64(req.Id))
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Image.ToImage(image), nil
}

// UpdateImageById ...
func (c ControllerImage) UpdateImageById(_ context.Context, req *api.UpdateImageRequest) (*api.Image, error) {

	image, errs, err := c.endpoint.Image.Update(c.dto.Image.FromUpdateImageRequest(req))
	if len(errs) > 0 {
		return nil, c.prepareErrors(errs)
	}

	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Image.ToImage(image), nil
}

// GetImageList ...
func (c ControllerImage) GetImageList(_ context.Context, req *api.GetImageListRequest) (*api.GetImageListResult, error) {

	items, total, err := c.endpoint.Image.GetList(int64(req.Limit), int64(req.Offset), req.Order, req.SortBy)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return c.dto.Image.ToImageListResult(items, uint32(total), req.Limit, req.Offset), nil
}

// DeleteImageById ...
func (c ControllerImage) DeleteImageById(_ context.Context, req *api.DeleteImageRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.Image.Delete(int64(req.Id)); err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// UploadImage ...
func (c ControllerImage) UploadImage(_ context.Context, req *api.UploadImageRequest) (*api.Image, error) {

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

		images, errs := c.endpoint.Image.Upload(form.File)

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
