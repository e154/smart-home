package rbac

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
	"fmt"
	c "github.com/e154/smart-home/api/cache"
	"github.com/e154/smart-home/api/models"
	"net/url"
	"errors"
	"github.com/e154/smart-home/lib/common"
	"github.com/dgrijalva/jwt-go"
	"regexp"
	"github.com/e154/smart-home/api/log"
)

var (
	cache c.Cache
	err error
	auth_urls []string = []string{"/api/v1/signin", "/api/v1/signout", "/api/v1/reset", "/api/v1/recovery"}
)

func AccessFilter() {
	beego.InsertFilter("/*", beego.BeforeRouter, func(ctx *context.Context) {

		//TODO uncomment
		// for dev
		if beego.BConfig.RunMode == "dev" {
			//return
		}

		requestURI := strings.ToLower(ctx.Request.RequestURI)
		method := strings.ToLower(ctx.Request.Method)

		for _, auth_url := range auth_urls {
			if auth_url == requestURI {
				return
			}
		}


		// get access_token
		var access_token string
		if access_token, err = getToken(ctx); err != nil || access_token == "" {
			ctx.ResponseWriter.WriteHeader(401)
			ctx.Output.JSON(&map[string]interface{}{"status": "error", "message": "Unauthorized access"}, true, false)
			return
		}

		// get access list
		var access_list models.AccessList
		var user *models.User
		if user, access_list, err = getAccessList(access_token); err != nil {
			ctx.ResponseWriter.WriteHeader(401)
			ctx.Output.JSON(&map[string]interface{}{"status": "error", "message": "Unauthorized access"}, true, false)
			return
		}

		// если id == 1 значит суперадмин
		if user.Id == 1 {
			return
		}

		if ret := accessDecision(requestURI, method, access_list); ret {
			return
		}

		if beego.BConfig.RunMode == "dev" {
			beego.Warning(fmt.Sprintf("access denied: %s url: %s:%s", user.Role.Name, requestURI, method))
		}

		ctx.ResponseWriter.WriteHeader(401)
		ctx.Output.JSON(&map[string]interface{}{"status": "error", "message": "Unauthorized access"}, true, false)
	})
}

func accessDecision(params, method string, access_list models.AccessList) bool {

	//
	for _, levels := range access_list {
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

// access_token
func getToken(ctx *context.Context) (access_token string, err error) {

	if access_token = ctx.Input.Header("access_token"); access_token != "" {
		return
	}

	var u *url.URL
	if u, err = url.Parse(ctx.Input.URI()); err != nil {
		return
	}

	var m url.Values
	if m, err = url.ParseQuery(u.RawQuery); err != nil {
		return
	}

	if len(m["access_token"]) == 0 || m["access_token"][0] == "" {
		err = errors.New("")
		return
	}

	access_token = m["access_token"][0]

	return
}

// получить лист доступа
func getAccessList(token string) (user *models.User, access_list models.AccessList, err error) {

	// cache init
	cacheKey := cache.GetKey(token)
	if bc := cache.IsExist(cacheKey); bc {
		if r := cache.Get(cacheKey); r != nil {
			// load data from cache
			data := r.(cacheData)
			user = data.user
			access_list = data.access_list
			err = nil

			return
		}
	}

	// load user info
	key := common.GetKey("hmacKey")
	var claims jwt.MapClaims
	if claims, err = common.ParseHmacToken(token, key); err != nil {
		log.Warnf("rbac: %s", err.Error())
		return
	}

	var ok bool
	if token, ok = claims["auth"].(string); !ok {
		log.Warnf("rbac: no auth var in token")
		return
	}

	if user, err = models.UserGetByAuthenticationToken(token); err != nil {
		log.Warnf("rbac: %s", err.Error())
		return
	}

	if err = user.LoadRelated(); err != nil {
		return
	}

	access_list = user.Role.GetFullAccessList()

	// save info to cache
	cache.Put("rbac", cacheKey, cacheData{
		user: user,
		access_list: access_list,
	})

	return
}

func UpdateCache() {
	cache.ClearAll()
}

func Initialize() {}

func init() {
	var err error
	cache = c.Cache{}
	cache.Name = "rbac"
	if cache.Cachetime, err = beego.AppConfig.Int64("cachetime"); err != nil {
		cache.Cachetime = 360000
	}
}