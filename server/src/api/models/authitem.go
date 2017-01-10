package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type AuthItem struct {
	Id			int64		`orm:"pk;auto" json:"id"`
	Name			string		`orm:"size(255)" valid:"Required;MaxSize(255)" json:"name"`
	Method			string		`orm:"size(255)" valid:"MaxSize(255)" json:"method"`
	Description		string		`orm:"size(255)" valid:"MaxSize(255)" json:"description"`
	Active			int		`orm:"-" json:"active"`
	Role			*Role		`orm:"rel(fk);null" json:"role"`
	Package			string		`orm:"size(255)" json:"package"`
	Frontend		int		`orm:"size(1);default(0)" json:"frontend"`
}

func (m *AuthItem) TableName() string {
	return beego.AppConfig.String("db_auth_items")
}

func init() {
	orm.RegisterModel(new(AuthItem))
}

// AddAuthItem insert a new AuthItem into database and returns
// last inserted Id on success.
func AddAuthItem(m *AuthItem) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAuthItemById retrieves AuthItem by Id. Returns error if
// Id doesn't exist
func GetAuthItemById(id int64) (v *AuthItem, err error) {
	o := orm.NewOrm()
	v = &AuthItem{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAuthItem retrieves all AuthItem matches certain condition. Returns empty list if
// no records exist
func GetAllAuthItem(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AuthItem))
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

	var l []AuthItem
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

// UpdateAuthItem updates AuthItem by Id and returns error if
// the record to be updated doesn't exist
func UpdateAuthItemById(m *AuthItem) (err error) {
	o := orm.NewOrm()
	v := AuthItem{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAuthItem deletes AuthItem by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAuthItem(id int64) (err error) {
	o := orm.NewOrm()
	v := AuthItem{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AuthItem{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
