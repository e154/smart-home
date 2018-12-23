package controllers

import (
	"github.com/gin-gonic/gin"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"github.com/e154/smart-home/api/server/v1/models"
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
		NewError(403, "bar request").Send(ctx)
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
