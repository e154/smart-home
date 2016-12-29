package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type MapText struct {
	Id           	int64  			`orm:"pk;auto;column(id)" json:"id"`
	Text           	string 			`orm:"" json:"text"`
	Style        	string			`orm:"" json:"style"`
}

func (m *MapText) TableName() string {
	return beego.AppConfig.String("db_map_text")
}

func init() {
	orm.RegisterModel(new(MapText))
}

// AddMapText insert a new MapText into database and returns
// last inserted Id on success.
func AddMapText(m *MapText) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetMapTextById retrieves MapText by Id. Returns error if
// Id doesn't exist
func GetMapTextById(id int64) (v *MapText, err error) {
	o := orm.NewOrm()
	v = &MapText{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllMapText retrieves all MapText matches certain condition. Returns empty list if
// no records exist
func GetAllMapText(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(MapText))
	// query k=v
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

	var l []MapText
	qs = qs.OrderBy(sortFields...)
	objects_count, err := qs.Count()
	if err != nil {
		return
	}
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
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
		meta = &map[string]int64{
			"objects_count": objects_count,
			"limit":         limit,
			"offset":        offset,
		}
		return ml, meta, nil
	}
	return nil, nil, err
}

// UpdateMapText updates MapText by Id and returns error if
// the record to be updated doesn't exist
func UpdateMapTextById(m *MapText) (err error) {
	o := orm.NewOrm()
	v := MapText{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMapText deletes MapText by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMapText(id int64) (err error) {
	o := orm.NewOrm()
	v := MapText{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&MapText{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
