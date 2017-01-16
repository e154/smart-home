package controllers

import (
	"encoding/json"
	"../models"
	"../../lib/common"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"fmt"
	"time"
)

var hmacSampleSecret []byte
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

	token := h.getHmacToken(user)

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

func (h *AuthController) getHmacToken(user *models.User) (tokenString string){

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"auth": user.AuthenticationToken,
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	var err error

	// Sign and get the complete encoded token as a string using the secret
	if tokenString, err = token.SignedString(hmacSampleSecret); err != nil {
		h.ErrHan(401, err.Error())
		return
	}

	return
}

func (h *AuthController) parseHmacToken(tokenString string) (jwt.MapClaims, error) {

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (c *AuthController) AccessList() {

	c.Data["json"] = &map[string]interface{}{"access_list": models.AccessConfigList}
	c.ServeJSON()
}

func init() {
	// Load sample key data
	if keyData, e := ioutil.ReadFile("keys/hmacKey"); e == nil {
		hmacSampleSecret = keyData
	} else {
		panic(e)
	}
}
