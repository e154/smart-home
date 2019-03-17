package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerImage struct {
	*ControllerCommon
}

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

	var newImage models.NewImage
	if err := ctx.ShouldBindJSON(&newImage); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	_, image, errs, err := AddImage(newImage, c.adaptors)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(image).Send(ctx)
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

	image, err := GetImageById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(image).Send(ctx)
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
//     $ref: '#/responses/Success'
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

	n := &models.UpdateImage{}
	if err := ctx.ShouldBindJSON(&n); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	n.Id = int64(aid)

	_, errs, err := UpdateImage(n, c.adaptors)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
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
	items, total, err := GetImageList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
	return
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

	if err := DeleteImageById(int64(aid), c.adaptors); err != nil {
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

	images, errs := UploadImages(form.File, c.adaptors)

	resp := NewSuccess()
	resp.SetData(&map[string]interface{}{
		"images": images,
		"errors": errs,
	})
	resp.Send(ctx)
}
