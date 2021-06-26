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

package dto

import (
	"github.com/e154/smart-home/api/stub/api"
	m "github.com/e154/smart-home/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Image struct{}

func NewImageDto() Image {
	return Image{}
}

func (i Image) ToImage(image *m.Image) (result *api.Image) {
	result = &api.Image{
		Id:        int32(image.Id),
		Thumb:     image.Thumb,
		Image:     image.Image,
		MimeType:  image.MimeType,
		Title:     image.Title,
		Size:      int32(image.Size),
		Name:      image.Name,
		CreatedAt: timestamppb.New(image.CreatedAt),
	}
	return
}

func (i Image) FromNewImageRequest(req *api.NewImageRequest) (image *m.Image) {
	image = &m.Image{
		Thumb:    req.Thumb,
		Image:    req.Image,
		MimeType: req.MimeType,
		Title:    req.Title,
		Name:     req.Name,
	}
	return
}

func (i Image) FromUpdateImageRequest(req *api.UpdateImageRequest) (image *m.Image) {
	image = &m.Image{
		Id:       int64(req.Id),
		Thumb:    req.Thumb,
		Image:    req.Image,
		MimeType: req.MimeType,
		Title:    req.Title,
		Name:     req.Name,
		Size:     int64(req.Size),
	}
	return
}

func (i Image) ToImageListResult(items []*m.Image, total, limit, offset uint32) (result *api.GetImageListResult) {

	result = &api.GetImageListResult{
		Items: make([]*api.Image, 0, len(items)),
		Meta: &api.GetImageListResult_Meta{
			Limit:        limit,
			ObjectsCount: total,
			Offset:       offset,
		},
	}

	for _, item := range items {
		result.Items = append(result.Items, &api.Image{
			Id:        int32(item.Id),
			Thumb:     item.Thumb,
			Url:       item.Url,
			Image:     item.Image,
			MimeType:  item.MimeType,
			Title:     item.Title,
			Size:      int32(item.Size),
			Name:      item.Name,
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}

	return
}
