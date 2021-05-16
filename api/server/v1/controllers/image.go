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
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

// ControllerImage ...
type ControllerImage struct {
	*ControllerCommon
}

// NewControllerImage ...
func NewControllerImage(common *ControllerCommon) *ControllerImage {
	return &ControllerImage{ControllerCommon: common}
}

// swagger:operation POST /image imageAdd
// ---
// parameters:
// - description: image params
//   in: body
//   name: image
//   required: true
//   schema:
//     $ref: '#/definitions/NewImage'
//     type: object
// summary: add new image
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - image
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Image'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerImage) Add(ctx *gin.Context) {

	var params models.NewImage
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	image := &m.Image{}
	common.Copy(&image, &params)

	image, errs, err := c.endpoint.Image.Add(image)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Image{}
	common.Copy(&result, &image)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /image/{id} imageGetById
// ---
// parameters:
// - description: Image ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get image by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - image
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Image'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerImage) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	image, err := c.endpoint.Image.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.Image{}
	common.Copy(&result, &image)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /image/{id} imageUpdateById
// ---
// parameters:
// - description: Image ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update image params
//   in: body
//   name: image
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateImage'
//     type: object
// summary: update image by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - image
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Image'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerImage) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params := &models.UpdateImage{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	params.Id = int64(aid)

	image := &m.Image{}
	common.Copy(&image, &params)

	image, errs, err := c.endpoint.Image.Update(image)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := &models.Image{}
	common.Copy(&result, &image)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /images imageList
// ---
// summary: get image list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - image
// parameters:
// - default: 10
//   description: limit
//   in: query
//   name: limit
//   required: true
//   type: integer
// - default: 0
//   description: offset
//   in: query
//   name: offset
//   required: true
//   type: integer
// - default: DESC
//   description: order
//   in: query
//   name: order
//   type: string
// - default: id
//   description: sort_by
//   in: query
//   name: sort_by
//   type: string
// responses:
//   "200":
//	   $ref: '#/responses/ImageList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerImage) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := c.endpoint.Image.GetList(int64(limit), int64(offset), order, sortBy)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.Image, 0)
	common.Copy(&result, items)

	resp := NewSuccess()
	resp.Page(limit, offset, total, result).Send(ctx)
}

// swagger:operation DELETE /image/{id} imageDeleteById
// ---
// parameters:
// - description: Image ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete image by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - image
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerImage) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.Image.Delete(int64(aid)); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation POST /image/upload imageUpload
// ---
// consumes:
//   - multipart/form-data
// parameters:
//   - in: formData
//     name: file
//     type: array
//     required: true
//     description: "image file"
//     items:
//       type: file
// summary: upload image files
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - image
// responses:
//   "200":
//	   $ref: '#/responses/NewObjectSuccess'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c *ControllerImage) Upload(ctx *gin.Context) {

	form, _ := ctx.MultipartForm()

	if len(form.File) == 0 {
		NewError(403, "http: no such file").Send(ctx)
		return
	}

	images, errs := c.endpoint.Image.Upload(form.File)

	resultImages := make([]*models.Image, 0)
	common.Copy(&resultImages, images)

	resp := NewSuccess()
	resp.SetData(&map[string]interface{}{
		"images": resultImages,
		"errors": errs,
	})
	resp.Send(ctx)
}
