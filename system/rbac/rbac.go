// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package rbac

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
	"strings"
)

var (
	log = common.MustGetLogger("rbac")
)

// AccessFilter ...
type AccessFilter struct {
	adaptors          *adaptors.Adaptors
	accessListService *access_list.AccessListService
}

// NewAccessFilter ...
func NewAccessFilter(adaptors *adaptors.Adaptors,
	accessListService *access_list.AccessListService) *AccessFilter {
	return &AccessFilter{
		adaptors:          adaptors,
		accessListService: accessListService,
	}
}

// Auth ...
func (f *AccessFilter) Auth(ctx *gin.Context) {

	requestURI := ctx.Request.RequestURI
	method := strings.ToLower(ctx.Request.Method)

	var err error

	// get access_token
	var accessToken string
	if accessToken, err = f.getToken(ctx); err != nil || accessToken == "" {
		ctx.AbortWithError(401, errors.New("unauthorized access"))
		return
	}

	if len(strings.Split(accessToken, ".")) != 3 {
		ctx.AbortWithError(401, errors.New("access token invalid"))
		return
	}

	// get access list
	var accessList access_list.AccessList
	var user *m.User
	if user, accessList, err = f.getAccessList(accessToken); err != nil {
		ctx.AbortWithError(403, errors.New("unauthorized access"))
		return
	}

	ctx.Set("currentUser", user)

	// если id == 1 is admin
	if user.Id == 1 || user.Role.Name == "admin" {
		return
	}

	if ret := f.accessDecision(requestURI, method, accessList); ret {
		return
	}

	log.Warnf(fmt.Sprintf("access denied: role(%s) [%s] url(%s)", user.Role.Name, method, requestURI))

	ctx.AbortWithError(403, errors.New("unauthorized access"))
}

// access_token
func (f *AccessFilter) getToken(ctx *gin.Context) (accessToken string, err error) {

	if accessToken = ctx.Request.Header.Get("access_token"); accessToken != "" {
		return
	}

	if accessToken = ctx.Request.Header.Get("Authorization"); accessToken != "" {
		return
	}

	if accessToken = ctx.Request.URL.Query().Get("access_token"); accessToken != "" {
		return
	}

	return
}

// получить лист доступа
func (f *AccessFilter) getAccessList(token string) (user *m.User, accessList access_list.AccessList, err error) {

	//TODO cache start

	// ger hmac key
	var variable m.Variable
	if variable, err = f.adaptors.Variable.GetByName("hmacKey"); err != nil {
		variable = m.Variable{
			Name:  "hmacKey",
			Value: common.ComputeHmac256(),
		}
		if err = f.adaptors.Variable.Add(variable); err != nil {
			log.Error(err.Error())
		}
	}

	hmacKey, err := hex.DecodeString(variable.Value)
	if err != nil {
		log.Error(err.Error())
	}

	// load user info
	var claims jwt.MapClaims
	if claims, err = common.ParseHmacToken(token, hmacKey); err != nil {
		//log.Warn(err.Error())
		return
	}

	//var ok bool
	//if token, ok = claims["auth"].(string); !ok {
	//	log.Warn("no auth var in token")
	//	return
	//}
	//
	//if user, err = f.adaptors.User.GetByAuthenticationToken(token); err != nil {
	//	return
	//}

	id, ok := claims["userId"]
	if !ok {
		err = fmt.Errorf("no userId var in token")
		return
	}

	var userId int64
	if userId, err = strconv.ParseInt(fmt.Sprintf("%v", id), 10, 0); err != nil {
		return
	}

	if user, err = f.adaptors.User.GetById(userId); err != nil {
		return
	}

	if accessList, err = f.accessListService.GetFullAccessList(user.Role); err != nil {
		return
	}

	//TODO cache end

	return
}

func (f *AccessFilter) accessDecision(params, method string, accessList access_list.AccessList) bool {

	for _, levels := range accessList {
		for _, item := range levels {
			for _, action := range item.Actions {
				if item.Method != method {
					continue
				}

				if ok, _ := regexp.MatchString(action, params); ok {
					return true
				}
			}
		}
	}

	return false
}
