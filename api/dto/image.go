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

package dto

import (
	stub "github.com/e154/smart-home/api/stub"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// Image ...
type Image struct{}

// NewImageDto ...
func NewImageDto() Image {
	return Image{}
}

// ToImage ...
func (i Image) ToImage(image *m.Image) (result *stub.ApiImage) {
	result = &stub.ApiImage{
		Id:        image.Id,
		Thumb:     image.Thumb,
		Image:     image.Image,
		MimeType:  image.MimeType,
		Title:     image.Title,
		Size:      image.Size,
		Name:      image.Name,
		Url:       image.Url,
		CreatedAt: image.CreatedAt,
	}
	return
}

// ToImageShort ...
func (i Image) ToImageShort(image *m.Image) (result *stub.ApiImage) {
	result = &stub.ApiImage{
		Id:   image.Id,
		Name: image.Name,
		Url:  image.Url,
	}
	return
}

// FromNewImageRequest ...
func (i Image) FromNewImageRequest(req *stub.ApiNewImageRequest) (image *m.Image) {
	image = &m.Image{
		Thumb:    req.Thumb,
		Image:    req.Image,
		MimeType: req.MimeType,
		Title:    req.Title,
		Name:     req.Name,
	}
	return
}

// FromUpdateImageRequest ...
func (i Image) FromUpdateImageRequest(req *stub.ImageServiceUpdateImageByIdJSONBody, id int64) (image *m.Image) {
	image = &m.Image{
		Id:       id,
		Thumb:    req.Thumb,
		Image:    req.Image,
		MimeType: req.MimeType,
		Title:    req.Title,
		Name:     req.Name,
		Size:     req.Size,
	}
	return
}

// ToImageListResult ...
func (i Image) ToImageListResult(images []*m.Image) []*stub.ApiImage {

	var items = make([]*stub.ApiImage, 0, len(images))
	for _, item := range images {
		items = append(items, &stub.ApiImage{
			Id:        item.Id,
			Thumb:     item.Thumb,
			Url:       item.Url,
			Image:     item.Image,
			MimeType:  item.MimeType,
			Title:     item.Title,
			Size:      item.Size,
			Name:      item.Name,
			CreatedAt: item.CreatedAt,
		})
	}

	return items
}

// ToImageList ...
func (i Image) ToImageList(items []*m.Image) (result *stub.ApiGetImageListByDateResult) {

	result = &stub.ApiGetImageListByDateResult{
		Items: make([]stub.ApiImage, 0, len(items)),
	}

	for _, item := range items {
		result.Items = append(result.Items, stub.ApiImage{
			Id:        item.Id,
			Thumb:     item.Thumb,
			Url:       item.Url,
			Image:     item.Image,
			MimeType:  item.MimeType,
			Title:     item.Title,
			Size:      item.Size,
			Name:      item.Name,
			CreatedAt: item.CreatedAt,
		})
	}

	return
}

// ToFilterList ...
func (i Image) ToFilterList(items []*m.ImageFilterList) (result stub.ApiGetImageFilterListResult) {

	result = stub.ApiGetImageFilterListResult{
		Items: make([]stub.GetImageFilterListResultfilter, 0, len(items)),
	}

	for _, item := range items {
		result.Items = append(result.Items, stub.GetImageFilterListResultfilter{
			Date:  item.Date,
			Count: int32(item.Count),
		})
	}

	return
}

func ImportImage(from *stub.ApiImage) (*int64, *m.Image) {
	if from == nil {
		return nil, nil
	}
	return common.Int64(from.Id), &m.Image{
		Id:       from.Id,
		Thumb:    from.Thumb,
		Image:    from.Image,
		MimeType: from.MimeType,
		Title:    from.Title,
		Name:     from.Name,
		Size:     from.Size,
	}
}
