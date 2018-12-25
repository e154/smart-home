package controllers

import (
	"github.com/gin-gonic/gin"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"github.com/e154/smart-home/api/server/v1/models"
	m "github.com/e154/smart-home/models"
)

type ControllerAuth struct {
	*ControllerCommon
}

func NewControllerAuth(common *ControllerCommon) *ControllerAuth {
	return &ControllerAuth{ControllerCommon: common}
}

// Auth godoc
// @tags auth
// @Summary Add new auth
// @Description
// @Produce json
// @Accept  json
// @Success 200 {object} models.AuthSignInResponse
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /signin [post]
// @Security BasicAuth
func (c ControllerAuth) SignIn(ctx *gin.Context) {

	email, password, ok := ctx.Request.BasicAuth()
	if !ok {
		NewError(403, "bad request").Send(ctx)
		return
	}

	currentUser, accessToken, err := SignIn(email, password, c.adaptors, ctx.ClientIP())
	if err != nil {
		code := 500
		switch err.Error() {
		case "user not found":
			code = 401
		case "password not valid":
			code = 403
		}

		NewError(code, err.Error()).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(&models.AuthSignInResponse{
		AccessToken: accessToken,
		CurrentUser: currentUser,
	}).Send(ctx)
}

// Auth godoc
// @tags auth
// @Summary Sign out
// @Description
// @Produce json
// @Accept  json
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /signout [post]
func (c ControllerAuth) SignOut(ctx *gin.Context) {

	u, ok := ctx.Get("currentUser")
	if !ok {
		NewError(403, "bad request").Send(ctx)
	}

	user, ok := u.(*m.User)
	if !ok {
		NewError(403, "bad request").Send(ctx)
	}

	if err := SignOut(user, c.adaptors); err != nil {
		NewError(500, err.Error()).Send(ctx)
		return
	}

	NewSuccess().Send(ctx)
}

// Auth godoc
// @tags auth
// @Summary Recovery access
// @Description
// @Produce json
// @Accept  json
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /recovery [post]
func (c ControllerAuth) Recovery(ctx *gin.Context) {

}

// Auth godoc
// @tags auth
// @Summary Reset password
// @Description
// @Produce json
// @Accept  json
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /reset [post]
func (c ControllerAuth) Reset(ctx *gin.Context) {

}

// Auth godoc
// @tags auth
// @Summary Get user access list
// @Description
// @Produce json
// @Accept  json
// @Success 200 {object} models.AccessList
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Security ApiKeyAuth
// @Router /access_list [get]
func (c ControllerAuth) AccessList(ctx *gin.Context) {

	u, ok := ctx.Get("currentUser")
	if !ok {
		NewError(403, "bad request").Send(ctx)
	}

	user, ok := u.(*m.User)
	if !ok {
		NewError(403, "bad request").Send(ctx)
	}

	accessList, err := AccessList(user, c.accessList)
	if err != nil {
		NewError(500, err.Error()).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(&map[string]interface{}{
		"access_list": accessList,
	}).Send(ctx)
}
