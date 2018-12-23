package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
	m "github.com/e154/smart-home/models"
)

type ControllerUser struct {
	*ControllerCommon
}

func NewControllerUser(common *ControllerCommon) *ControllerUser {
	return &ControllerUser{ControllerCommon: common}
}

// User godoc
// @tags user
// @Summary Add new user
// @Description
// @Produce json
// @Accept  json
// @Param user body models.NewUser true "user params"
// @Success 200 {object} models.NewObjectSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /user [post]
func (c ControllerUser) AddUser(ctx *gin.Context) {

	var params models.NewUser
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	var currentUser *m.User
	if user, ok := ctx.Get("currentUser"); ok {
		currentUser = user.(*m.User)
	}

	_, id, errs, err := AddUser(params, c.adaptors, currentUser)
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
	resp.Item("id", id).Send(ctx)
}

// User godoc
// @tags user
// @Summary Show user
// @Description Get user by id
// @Produce json
// @Accept  json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserByIdModel
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /user/{id} [Get]
func (c ControllerUser) GetUserById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	user, err := GetUserById(int64(aid), c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("user", user).Send(ctx)
}

// User godoc
// @tags user
// @Summary Update user
// @Description Update user by id
// @Produce json
// @Accept  json
// @Param  id path int true "User ID"
// @Param  user body models.UpdateUser true "Update user"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /user/{id} [Put]
//func (c ControllerUser) UpdateUser(ctx *gin.Context) {
//
//	aid, err := strconv.Atoi(ctx.Param("id"))
//	if err != nil {
//		log.Error(err.Error())
//		NewError(400, err).Send(ctx)
//		return
//	}
//
//	n := &m.User{}
//	if err := ctx.ShouldBindJSON(&n); err != nil {
//		log.Error(err.Error())
//		NewError(400, err).Send(ctx)
//		return
//	}
//
//	n.Id = int64(aid)
//
//	_, errs, err := UpdateUser(n, c.adaptors, c.core)
//	if len(errs) > 0 {
//		err400 := NewError(400)
//		err400.ValidationToErrors(errs).Send(ctx)
//		return
//	}
//
//	if err != nil {
//		NewError(500, err).Send(ctx)
//		return
//	}
//
//	resp := NewSuccess()
//	resp.Send(ctx)
//}

// User godoc
// @tags user
// @Summary User list
// @Description Get user list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(id)
// @Success 200 {array} models.UserShotModel
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /users [Get]
func (c ControllerUser) GetUserList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetUserList(limit, offset, order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, int(total), items).Send(ctx)
	return
}

// User godoc
// @tags user
// @Summary Delete user
// @Description Delete user by id
// @Produce json
// @Accept  json
// @Param  id path int true "User ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /user/{id} [Delete]
func (c ControllerUser) DeleteUserById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := DeleteUserById(int64(aid), c.adaptors); err != nil {
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
