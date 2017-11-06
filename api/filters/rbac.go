package filters

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
	"fmt"
	"github.com/e154/smart-home/api/common"
	"github.com/e154/smart-home/api/models"
	"net/url"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"regexp"
	"github.com/e154/smart-home/api/log"
	"github.com/e154/smart-home/api/variable"
	"encoding/hex"
)

var (
	cache common.Cache
	err error
	auth_urls []string = []string{"/api/v1/signin", "/api/v1/signout", "/api/v1/reset", "/api/v1/recovery"}
)

func AccessFilter() {
	beego.InsertFilter("/*", beego.BeforeRouter, func(ctx *context.Context) {

		//if beego.BConfig.RunMode == "dev" {
		//	return
		//}

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
	//cacheKey := cache.GetKey(token)
	//if bc := cache.IsExist(cacheKey); bc {
	//	if r := cache.Get(cacheKey); r != nil {
	//		// load data from cache
	//		data := r.(cacheData)
	//		user = data.user
	//		access_list = data.access_list
	//		err = nil
	//
	//		return
	//	}
	//}

	// ger hmac key
	key, ok := variable.Get("hmacKey")
	if !ok {
		key = common.ComputeHmac256()
		if err = variable.Set("hmacKey", key); err != nil {
			log.Error(err.Error())
		}
	}

	hmacKey, err := hex.DecodeString(key)
	if err != nil {
		log.Error(err.Error())
	}

	// load user info
	var claims jwt.MapClaims
	if claims, err = ParseHmacToken(token, hmacKey); err != nil {
		log.Warnf("rbac: %s", err.Error())
		return
	}

	if token, ok = claims["auth"].(string); !ok {
		log.Warnf("rbac: no auth var in token")
		return
	}

	if user, err = models.UserGetByAuthenticationToken(token); err != nil {
		return
	}

	if err = user.LoadRelated(); err != nil {
		return
	}

	access_list = user.Role.GetFullAccessList()

	// save info to cache
	//cache.Put("rbac", cacheKey, cacheData{
	//	user: user,
	//	access_list: access_list,
	//})

	return
}

func UpdateCache() {
	cache.ClearAll()
}


func ParseHmacToken(tokenString string, key []byte) (jwt.MapClaims, error) {

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
		return key, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func Initialize() {}

func init() {
	var err error
	cache = common.Cache{}
	cache.Name = "rbac"
	if cache.Cachetime, err = beego.AppConfig.Int64("cachetime"); err != nil {
		cache.Cachetime = 360000
	}
}