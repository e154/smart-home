package use_case

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common"
	"errors"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/jinzhu/copier"
	"encoding/hex"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/e154/smart-home/system/access_list"
)

const (
	AdminId = 1
)

func SignIn(email, password string, adaptors *adaptors.Adaptors, ip string) (currentUser *models.CurrentUserModel, accessToken string, err error) {

	var user *m.User

	if user, err = adaptors.User.GetByEmail(email); err != nil {
		err = errors.New("user not found")
		return
	} else if user.EncryptedPassword != common.Pwdhash(password) {
		err = errors.New("password not valid")
		return
	} else if user.Status == "blocked" && user.Id != AdminId {
		err = errors.New("account is blocked")
		return
	}

	adaptors.User.SignIn(user, ip)

	if _, err = adaptors.User.NewToken(user); err != nil {
		return
	}

	currentUser = &models.CurrentUserModel{}
	copier.Copy(&currentUser, &user)

	// meta
	currentUser.Meta = make([]*models.UserByIdModelMeta, 0)
	for _, meta := range user.Meta {
		m := &models.UserByIdModelMeta{}
		copier.Copy(&m, &meta)
		currentUser.Meta = append(currentUser.Meta, m)
	}

	// history
	currentUser.History = make([]*models.UserHistory, 0)
	for _, story := range user.History {
		s := &models.UserHistory{}
		copier.Copy(&s, &story)
		currentUser.History = append(currentUser.History, s)
	}

	// role
	currentUser.Role = &models.RoleModel{}
	copier.Copy(&currentUser.Role, &user.Role)

	// image
	if user.Image != nil {
		currentUser.Image = &models.Image{}
		copier.Copy(&currentUser.Image, &user.Image)
	}

	// ger hmac key
	var variable *m.Variable
	if variable, err = adaptors.Variable.GetByName("hmacKey"); err != nil {
		variable = &m.Variable{
			Name:  "hmacKey",
			Value: common.ComputeHmac256(),
		}
		if err = adaptors.Variable.Add(variable); err != nil {
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

func SignOut(user *m.User, adaptors *adaptors.Adaptors) (err error) {

	err = adaptors.User.ClearToken(user)

	return
}

func Recovery() {}

func Reset() {}

func AccessList(user *m.User, accessListService *access_list.AccessListService) (accessList *access_list.AccessList, err error) {
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
