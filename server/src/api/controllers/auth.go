package controllers

import (
	"encoding/json"
	"../models"
	"../../lib/common"
	"time"
)

const ADMIN_ID = 1

// AuthController operations for Auth
type AuthController struct {
	CommonController
}

// URLMapping ...
func (c *AuthController) URLMapping() {
	c.Mapping("SignIn", c.SignIn)
	c.Mapping("SignOut", c.SignOut)
	c.Mapping("Recovery", c.Recovery)
	c.Mapping("Reset", c.Reset)
	c.Mapping("AccessList", c.AccessList)
}

// @Title SignIn
// @Description user account page
// @Param	body
// @Success 201 {object}
// @Failure 403
// @router /signin [post]
func (h *AuthController) SignIn() {

	input := map[string]string{}
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &input); err != nil {
		h.ErrHan(403, err.Error())
		return
	}

	var user *models.User
	var err error

	if user, err = models.UserGetByEmail(input["email"]); err != nil {
		h.ErrHan(401, "Пользователь не найден")
		return
	} else if user.EncryptedPassword != common.Pwdhash(input["password"]) {
		h.ErrHan(403, "Не верный пароль")
		return
	} else if user.Status == "blocked" && user.Id != ADMIN_ID {
		h.ErrHan(401, "Аккаунт заблокирован")
		return
	}

	user.SignIn(h.Ctx.Input.IP())
	user.LoadRelated()
	user.NewToken()

	//access_list := user.Role.GetFullAccessList()
	//fmt.Println(access_list)
	user.Role.GetAccessList()

	current_user := map[string]interface{}{
		"id": user.Id,
		"nickname": user.Nickname,
		"first_name": user.FirstName,
		"last_name": user.LastName,
		"email": user.Email,
		"history": user.History,
		"avatar": user.Avatar,
		"sign_in_count": user.SignInCount,
		"meta": user.Meta,
		"role": user.Role,
	}

	key := common.GetKey("hmacKey")
	data := map[string]interface{}{
		"auth": user.AuthenticationToken,
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	}

	var token string
	if token, err = common.GetHmacToken(data, key); err != nil {
		h.ErrHan(403, err.Error())
		return
	}

	h.Data["json"] = &map[string]interface{}{"token": token, "current_user": current_user}
	h.ServeJSON()
}

// @Title SignOut
// @Description user account page
// @Param	body
// @Success 201 {object}
// @Failure 403
// @router /signout [post]
func (h *AuthController) SignOut() {}

// @Title Recovery
// @Description user account page
// @Param	body
// @Success 201 {object}
// @Failure 403
// @router /recovery [post]
func (h *AuthController) Recovery() {}

// @Title Reset
// @Description user account page
// @Param	body
// @Success 201 {object}
// @Failure 403
// @router /reset [post]
func (h *AuthController) Reset() {}

func (c *AuthController) AccessList() {

	c.Data["json"] = &map[string]interface{}{"access_list": models.AccessConfigList}
	c.ServeJSON()
}