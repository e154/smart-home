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

// Image godoc
// @tags image
// @Summary Add new image
// @Description
// @Produce json
// @Accept  json
// @Param image body models.NewImage true "image params"
// @Success 200 {object} models.Image
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /image [post]
// @Security ApiKeyAuth
func (c ControllerImage) AddImage(ctx *gin.Context) {

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

// Image godoc
// @tags image
// @Summary Show image
// @Description Get image by id
// @Produce json
// @Accept  json
// @Param id path int true "Image ID"
// @Success 200 {object} models.Image
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /image/{id} [Get]
// @Security ApiKeyAuth
func (c ControllerImage) GetImageById(ctx *gin.Context) {

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

// Image godoc
// @tags image
// @Summary Update image
// @Description Update image by id
// @Produce json
// @Accept  json
// @Param  id path int true "Image ID"
// @Param  image body models.UpdateImage true "Update image"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /image/{id} [Put]
// @Security ApiKeyAuth
func (c ControllerImage) UpdateImage(ctx *gin.Context) {

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

// Image godoc
// @tags image
// @Summary Image list
// @Description Get image list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(id)
// @Success 200 {object} models.ImageListModel
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /images [Get]
// @Security ApiKeyAuth
func (c ControllerImage) GetImageList(ctx *gin.Context) {

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

// Image godoc
// @tags image
// @Summary Delete image
// @Description Delete image by id
// @Produce json
// @Accept  json
// @Param  id path int true "Image ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /image/{id} [Delete]
// @Security ApiKeyAuth
func (c ControllerImage) DeleteImageById(ctx *gin.Context) {

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
