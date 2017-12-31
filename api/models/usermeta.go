package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type UserMeta struct {
	Id		int64		`orm:"pk;auto" json:"id"`
	User		*User		`orm:"rel(fk)" json:"-"`
	Key		string		`orm:"size(255)" valid:"MaxSize(255)" json:"key"`
	Value		string		`orm:"size(255)" valid:"MaxSize(255)" json:"value"`
}

func (m *UserMeta) TableName() string {
	return beego.AppConfig.String("db_user_metas")
}

// AddUserMeta insert a new UserMeta into database and returns
// last inserted Id on success.
func AddUserMeta(m *UserMeta) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserMetaById retrieves UserMeta by Id. Returns error if
// Id doesn't exist
func GetUserMetaById(id int64) (v *UserMeta, err error) {
	o := orm.NewOrm()
	v = &UserMeta{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserMeta retrieves all UserMeta matches certain condition. Returns empty list if
// no records exist
func GetAllUserMeta(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserMeta))
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

	var l []UserMeta
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

// UpdateUserMeta updates UserMeta by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserMetaById(m *UserMeta) (err error) {
	o := orm.NewOrm()
	v := UserMeta{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserMeta deletes UserMeta by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserMeta(id int64) (err error) {
	o := orm.NewOrm()
	v := UserMeta{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserMeta{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func SetUserMeta(u *User, key, value string) (err error) {

	o := orm.NewOrm()

	meta := new(UserMeta)
	qs := o.QueryTable(meta).Filter("user_id", u.Id).Filter("key", key)

	if exist := qs.Exist(); exist {

		if _, err = qs.Update(orm.Params{"value": value}); err != nil {
			return
		}

	} else {
		meta.Key = key
		meta.Value = value
		meta.User = u

		if _, err = o.Insert(meta); err != nil {
			return
		}
	}

	return
}

func GetUserMeta(u *User, key string) string {

	o := orm.NewOrm()

	meta := UserMeta{}
	if err := o.QueryTable(&meta).Filter("user_id", u.Id).Filter("key", key).One(&meta); err != nil {
		meta.Key = key
		meta.Value = ""
		meta.User = u
		o.Insert(&meta)
	}

	return meta.Value
}