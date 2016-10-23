package controllers

import  (
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
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