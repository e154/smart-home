package endpoint

import (
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"encoding/hex"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/e154/smart-home/system/access_list"
)

const (
	AdminId = 1
)

type AuthEndpoint struct {
	*CommonEndpoint
}

func NewAuthEndpoint(common *CommonEndpoint) *AuthEndpoint {
	return &AuthEndpoint{
		CommonEndpoint: common,
	}
}

func (a *AuthEndpoint) SignIn(email, password string, ip string) (user *m.User, accessToken string, err error) {

	if user, err = a.adaptors.User.GetByEmail(email); err != nil {
		err = errors.New("user not found")
		return
	} else if user.EncryptedPassword != common.Pwdhash(password) {
		err = errors.New("password not valid")
		return
	} else if user.Status == "blocked" && user.Id != AdminId {
		err = errors.New("account is blocked")
		return
	}

	if err = a.adaptors.User.SignIn(user, ip); err != nil {
		return
	}

	if _, err = a.adaptors.User.NewToken(user); err != nil {
		return
	}

	// ger hmac key
	var variable *m.Variable
	if variable, err = a.adaptors.Variable.GetByName("hmacKey"); err != nil {
		variable = &m.Variable{
			Name:  "hmacKey",
			Value: common.ComputeHmac256(),
		}
		if err = a.adaptors.Variable.Add(variable); err != nil {
			log.Error(err.Error())
		}
	}

	var hmacKey []byte
	hmacKey, err = hex.DecodeString(variable.Value)
	if err != nil {
		return
	}

	//key := common.GetKey("hmacKey")
	data := map[string]interface{}{
		"auth": user.AuthenticationToken,
		"nbf":  time.Now().Unix(),
	}

	if accessToken, err = GetHmacToken(data, hmacKey); err != nil {
		return
	}

	log.Infof("Successful login, user: %s", user.Email)

	return
}

func (a *AuthEndpoint) SignOut(user *m.User) (err error) {
	err = a.adaptors.User.ClearToken(user)
	return
}

func (a *AuthEndpoint) Recovery() {}

func (a *AuthEndpoint) Reset() {}

func (a *AuthEndpoint) AccessList(user *m.User, accessListService *access_list.AccessListService) (accessList *access_list.AccessList, err error) {
	accessList = accessListService.List
	return
}

func GetHmacToken(data map[string]interface{}, key []byte) (tokenString string, err error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(data))

	// Sign and get the complete encoded token as a string using the secret
	if tokenString, err = token.SignedString(key); err != nil {
		return
	}

	return
}
