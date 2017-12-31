package controllers

import  (
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"reflect"
	"html/template"
	"encoding/json"
	"github.com/e154/smart-home/api/log"
	"github.com/e154/smart-home/api/models"
	"github.com/pkg/errors"
	"github.com/e154/smart-home/api/variable"
	"github.com/e154/smart-home/api/common"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
)

type Request struct {
	Query  map[string]string                `json:"query"`
	Fields []string                        `json:"fields"`
	Sortby []string                        `json:"sortby"`
	Order  []string                        `json:"order"`
	Offset int64                        `json:"offset"`
	Limit  int64                        `json:"limit"`
}

type CommonController struct {
	beego.Controller
}

func (c *CommonController) ErrHan(code int, message string) {

	switch code {
	case 401, 403:
		log.Warn("error:", message)
	}

	c.Ctx.ResponseWriter.WriteHeader(code)
	c.Data["json"] = &map[string]interface{}{"status":"error", "message": message}
	c.ServeJSON()
}

func (c *CommonController) pagination() (query map[string]string, fields []string, sortby []string, order []string,
offset int64, limit int64) {

	link, _ := url.ParseRequestURI(c.Ctx.Request.URL.String())
	q := link.Query()

	query = map[string]string{}
	fields = []string{}
	sortby = []string{}
	order = []string{}

	values, _ := url.ParseQuery(q.Encode())
	if val, ok := values["query"]; ok {
		for _, v := range val {
			json.Unmarshal([]byte(v), &query)
		}
	}

	if val, ok := q["fields"]; ok {
		for _, v := range val {
			fields = append(fields, v)
		}
	}

	if val, ok := q["sortby"]; ok {
		for _, v := range val {
			sortby = append(sortby, v)
		}
	}

	if val, ok := q["order"]; ok {
		for _, v := range val {
			order = append(order, v)
		}
	}

	if val, ok := q["offset"]; ok {
		offset, _ = strconv.ParseInt(val[0], 10, 0)
	}

	if val, ok := q["limit"]; ok {
		limit, _ = strconv.ParseInt(val[0], 10, 0)
	}

	return
}

func (c *CommonController) GetTemplate() string {

	templatetype := beego.AppConfig.String("template_type")
	if templatetype == "" {
		templatetype = "default"
	}
	return templatetype
}

func (c *CommonController) GetCurrentUser() (user *models.User, err error) {

	token := c.Ctx.Input.Header("access_token")
	if token == "" {
		err = errors.New("access_token is empty")
		return
	}

	key, ok := variable.Get("hmacKey")
	if !ok {
		key = common.ComputeHmac256()
		if err = variable.Set("hmacKey", key); err != nil {
			return
		}
	}

	hmacKey, err := hex.DecodeString(key)
	if err != nil {
		return
	}

	// load user info
	var claims jwt.MapClaims
	if claims, err = common.ParseHmacToken(token, hmacKey); err != nil {
		//log.Warnf("rbac: %s", err.Error())
		return
	}

	if token, ok = claims["auth"].(string); !ok {
		//log.Warnf("rbac: no auth var in token")
		return
	}

	if user, err = models.UserGetByAuthenticationToken(token); err != nil {
		return
	}

	return
}


func (c *CommonController) Prepare() {


}

func init() {
	beego.AddFuncMap("safeHtml", func(s string) template.HTML {return template.HTML(s)})
	beego.AddFuncMap("safeCss", func(s string) template.CSS {return template.CSS(s)})
	beego.AddFuncMap("safeUrl", func(s string) template.URL {return template.URL(s)})
	beego.AddFuncMap("safeJs", func(s string) template.JS {return template.JS(s)})
	beego.AddFuncMap("attr", func(s string) template.HTMLAttr {return template.HTMLAttr(s)})
	beego.AddFuncMap("last", func(i int, s interface{}) bool {return i == reflect.ValueOf(s).Len() - 1})
	beego.AddFuncMap("len", func(s interface{}) int {return reflect.ValueOf(s).Len()})
}