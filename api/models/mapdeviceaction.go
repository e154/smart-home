package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"time"
)

type MapDeviceAction struct {
	Id           	int64  			`orm:"pk;auto;column(id)" json:"id"`
	Type           	string  		`orm:"" json:"type"`
	MapDevice      	*MapDevice		`orm:"rel(fk)" json:"map_device"`
	DeviceAction    *DeviceAction		`orm:"rel(fk)" json:"device_action"`
	Image		*Image			`orm:"rel(fk);null" json:"image"`
	Created_at   	time.Time		`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at    	time.Time		`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
}

func (m *MapDeviceAction) TableName() string {
	return beego.AppConfig.String("db_map_device_actions")
}

// AddMapDeviceAction insert a new MapDeviceAction into database and returns
// last inserted Id on success.
func AddMapDeviceAction(m *MapDeviceAction) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetMapDeviceActionById retrieves MapDeviceAction by Id. Returns error if
// Id doesn't exist
func GetMapDeviceActionById(id int64) (v *MapDeviceAction, err error) {
	o := orm.NewOrm()
	v = &MapDeviceAction{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllMapDeviceAction retrieves all MapDeviceAction matches certain condition. Returns empty list if
// no records exist
func GetAllMapDeviceAction(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(MapDeviceAction))
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

	var l []MapDeviceAction
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

// UpdateMapDeviceAction updates MapDeviceAction by Id and returns error if
// the record to be updated doesn't exist
func UpdateMapDeviceActionById(m *MapDeviceAction) (err error) {
	o := orm.NewOrm()
	v := MapDeviceAction{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMapDeviceAction deletes MapDeviceAction by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMapDeviceAction(id int64) (err error) {
	o := orm.NewOrm()
	v := MapDeviceAction{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&MapDeviceAction{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// AddMultipleMapDeviceAction
// Use a prepared statement to increase inserting speed with multiple inserts.
func AddMultipleMapDeviceAction(actions []*MapDeviceAction) (ids []int64, errs []error) {

	o := orm.NewOrm()
	qs := o.QueryTable(&MapDeviceAction{})
	i, _ := qs.PrepareInsert()
	for _, action := range actions {
		id, err := i.Insert(action)
		if err != nil {
			errs = append(errs, err)
		} else {
			ids = append(ids, id)
		}
	}
	// PREPARE INSERT INTO user (`name`, ...) VALUES (?, ...)
	// EXECUTE INSERT INTO user (`name`, ...) VALUES ("slene", ...)
	// EXECUTE ...
	// ...
	i.Close() // Don't forget to close the statement

	return
}