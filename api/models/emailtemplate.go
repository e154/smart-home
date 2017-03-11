package models

import (
	"encoding/json"
	"unicode/utf8"
	"regexp"
	"sort"
	"strings"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego/orm"
	"errors"
	"reflect"
)

type Field	struct {
	Name		string		`json:"name"`
	Value		string		`json:"value"`
}

type EmailTemplate struct {
	Items		[]string		`json:"items"`
	Title		string			`json:"title"`
	Fields		[]*Field		`json:"fields"`
}

type EmailRender struct {
	Subject		string
	Body		string
}

// GetAllEmailTemplates retrieves all EmailItem matches certain condition. Returns empty list if
// no records exist
func GetAllEmailTemplate(query map[string]string, fields []string, sortby []string, order []string,
offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(EmailItem))
	// query k=v
	query["type"] = "template"
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []EmailItem
	qs = qs.OrderBy(sortFields...)
	objects_count, err := qs.Count()
	if err != nil {
		return
	}
	if _, err := qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		meta := &map[string]int64{
			"objects_count": objects_count,
			"limit": limit,
			"offset": offset,
		}
		return ml, meta, nil
	}
	return nil, nil, err
}

func EmailTemplateGetByName(name string) (template *EmailItem, err error) {

	o := orm.NewOrm()
	template = new(EmailItem)
	err = o.QueryTable(template).Filter("type", "template").Filter("name", name).One(template)
	return
}

func EmailTemplatePreview(template *EmailTemplate) (buf string, err error) {

	_, items, err := EmailItemGetList()
	if err != nil {
		return
	}

	result := Items{}
	for _, item := range template.Items {
		getItemParents(items, &result, item)
	}

	sort.Sort(result)

	// замена [xxxx:block] на реальные блоки
	for key, item := range result {
		if item.Status != "active" {
			continue
		}

		if key == 0 {
			buf = item.Content
		} else {
			buf = strings.Replace(buf, fmt.Sprintf("[%s:block]", item.Name), item.Content, -1)
		}
	}

	// поиск маркера [xxx:content] и замена на контейнер с возможностью редактирования
	reg := regexp.MustCompile(`(\[{1}[a-z]{2,64}\:content\]{1})`)
	reg2 := regexp.MustCompile(`(\[{1})([a-z]{2,64})(\:)([content]+)([\]]{1})`)
	markers := reg.FindAllString(buf, -1)
	for _, m := range markers {
		marker := reg2.FindStringSubmatch(m)[2]

		var f string = m
		for _, field := range template.Fields {
			if field.Name == marker {
				if utf8.RuneCountInString(field.Value) > 0 {
					f = field.Value
				}
			}
		}

		buf = strings.Replace(buf, m, fmt.Sprintf("<div class=\"edit_inline\" data-id=\"%s\">%s</div>", marker, f), -1)
	}

	// скрыть не использованные маркеры [xxxx:block] блоков
	reg = regexp.MustCompile(`(\[{1}[a-z]{2,64}\:block\]{1})`)
	blocks := reg.FindAllString(buf, -1)
	for _, block := range blocks {
		buf = strings.Replace(buf, block, "", -1)
	}

	return
}

func EmailTemplateAddNew(body []byte) (id int64, b bool, valid validation.Validation, err error) {

	tpl := new(EmailItem)
	if err = json.Unmarshal(body, tpl); err != nil {
		return
	}

	tpl.Parent = ""
	tpl.Status = "active"
	tpl.Type = "template"

	b, err = valid.Valid(tpl)
	if err != nil || !b {
		return
	}

	o := orm.NewOrm()
	id, err = o.Insert(tpl)

	return
}

func EmailTemplateUpdate(body []byte, name string) (b bool, valid validation.Validation, err error) {

	o := orm.NewOrm()
	tpl := new(EmailItem)
	if err = o.QueryTable(tpl).Filter("name", name).Filter("type", "template").One(tpl); err != nil {
		return
	}

	if err = json.Unmarshal(body, tpl); err != nil {
		return
	}

	tpl.Parent = ""
	tpl.Status = "active"
	tpl.Type = "template"

	b, err = valid.Valid(tpl)
	if err != nil || !b {
		return
	}

	_, err = o.Update(tpl)

	return
}

func EmailTemplateDelete(name string) (err error) {

	o := orm.NewOrm()
	tpl := new(EmailItem)
	_, err = o.QueryTable(tpl).Filter("name", name).Filter("type", "template").Delete()
	return
}

func EmailTemplateRender(name string, params map[string]interface {}) (r *EmailRender, err error) {

	var item *EmailItem
	var template *EmailTemplate
	var items Items

	if item, err = EmailItemGet(name); err != nil {
		return
	}

	if template, err = item.GetTemplate(); err != nil {
		return
	}

	if _, items, err = EmailItemGetList(); err != nil {
		return
	}

	result := Items{}
	for _, item := range template.Items {
		getItemParents(items, &result, item)
	}

	sort.Sort(result)
	var buf string

	// замена [xxxx:block] на реальные блоки
	for key, item := range result {
		if item.Status != "active" {
			continue
		}

		if key == 0 {
			buf = item.Content
		} else {
			buf = strings.Replace(buf, fmt.Sprintf("[%s:block]", item.Name), item.Content, -1)
		}
	}

	// поиск маркера [xxx:content] и замена на контейнер с возможностью редактирования
	reg := regexp.MustCompile(`(\[{1}[a-z]{2,64}\:content\]{1})`)
	reg2 := regexp.MustCompile(`(\[{1})([a-z]{2,64})(\:)([content]+)([\]]{1})`)
	markers := reg.FindAllString(buf, -1)
	var f string
	for _, m := range markers {
		marker := reg2.FindStringSubmatch(m)[2]

		f = m
		for _, field := range template.Fields {
			if field.Name == marker {
				if utf8.RuneCountInString(field.Value) > 0 {
					f = field.Value
				}
			}
		}

		buf = strings.Replace(buf, m, f, -1)
	}

	// скрыть не использованные маркеры [xxxx:block] блоков
	reg = regexp.MustCompile(`(\[{1}[a-z]{2,64}\:block\]{1})`)
	blocks := reg.FindAllString(buf, -1)
	for _, block := range blocks {
		buf = strings.Replace(buf, block, "", -1)
	}

	// заполнение формы
	title := template.Title
	for k, v := range params {
		buf = strings.Replace(buf, fmt.Sprintf("[%s]", k), v.(string), -1)
		title = strings.Replace(title, fmt.Sprintf("[%s]", k), v.(string), -1)
	}

	r = new(EmailRender)
	r.Subject = title
	r.Body = buf

	return
}
