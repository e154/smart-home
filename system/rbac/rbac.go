package rbac

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/op/go-logging"
	"errors"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"fmt"
	"regexp"
)

var (
	log = logging.MustGetLogger("rbac")
)

type AccessFilter struct {
	adaptors          *adaptors.Adaptors
	accessListService *access_list.AccessListService
}

func NewAccessFilter(adaptors *adaptors.Adaptors,
	accessListService *access_list.AccessListService) *AccessFilter {
	return &AccessFilter{
		adaptors:          adaptors,
		accessListService: accessListService,
	}
}

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

	// get access list
	var accessList access_list.AccessList
	var user *m.User
	if user, accessList, err = f.getAccessList(accessToken); err != nil {
		ctx.AbortWithError(403, errors.New("unauthorized access"))
		return
	}

	ctx.Set("currentUser", user)

	// если id == 1 is admin
	if user.Id == 1 {
		return
	}

	if ret := f.accessDecision(requestURI, method, accessList); ret {
		return
	}

	log.Warningf(fmt.Sprintf("access denied: role(%s) [%s] url(%s)", user.Role.Name, method, requestURI))

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
	var variable *m.Variable
	if variable, err = f.adaptors.Variable.GetByName("hmacKey"); err != nil {
		variable = &m.Variable{
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
		//log.Warning(err.Error())
		return
	}

	var ok bool
	if token, ok = claims["auth"].(string); !ok {
		log.Warning("no auth var in token")
		return
	}

	if user, err = f.adaptors.User.GetByAuthenticationToken(token); err != nil {
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
