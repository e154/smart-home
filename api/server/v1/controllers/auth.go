package controllers

import (
	"github.com/gin-gonic/gin"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"github.com/e154/smart-home/api/server/v1/models"
	m "github.com/e154/smart-home/models"
	"net/http"
)

type ControllerAuth struct {
	*ControllerCommon
}

func NewControllerAuth(common *ControllerCommon) *ControllerAuth {
	return &ControllerAuth{ControllerCommon: common}
}

// swagger:operation POST /signin authSignin
// ---
// summary: sign in
// description:
// security:
// - BasicAuth: []
// tags:
// - auth
// responses:
//   "200":
//	   $ref: '#/responses/AuthSignInResponse'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "user not found"
//   "403":
//     description: "password not valid"
//   "500":
//	   $ref: '#/responses/Error'
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

// swagger:operation POST /signout authSignout
// ---
// summary: sign out
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - auth
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "401":
//     description: "user not found"
//   "403":
//     $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
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

// swagger:operation POST /recovery authRecovery
// ---
// summary: recovery access
// description:
// tags:
// - auth
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//     $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAuth) Recovery(ctx *gin.Context) {
	ctx.String(http.StatusOK, "operation Recovery has not yet been implemented")
}

// swagger:operation POST /reset authReset
// ---
// summary: reset access
// description:
// tags:
// - auth
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//     $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerAuth) Reset(ctx *gin.Context) {
	ctx.String(http.StatusOK, "operation Reset has not yet been implemented")
}

// swagger:operation GET /access_list authGetAccessList
// ---
// summary: get user access list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - auth
// responses:
//   "200":
//	   $ref: '#/responses/AccessList'
//   "401":
//     description: "user not found"
//   "403":
//     $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
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
