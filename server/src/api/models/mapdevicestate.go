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

type MapDeviceState struct {
	Id           	int64  			`orm:"pk;auto;column(id)" json:"id"`
	Style      	string			`orm:"" json:"style"`
	MapDevice      	*MapDevice		`orm:"rel(fk)" json:"map_device"`
	DeviceState    	*DeviceState		`orm:"rel(fk);null" json:"device_state"`
	Image		*Image			`orm:"rel(fk);null" json:"image"`
	Created_at   	time.Time		`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at    	time.Time		`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
}

func (m *MapDeviceState) TableName() string {
	return beego.AppConfig.String("db_map_device_states")
}

func init() {
	orm.RegisterModel(new(MapDeviceState))
}

// AddMapDeviceState insert a new MapDeviceState into database and returns
// last inserted Id on success.
func AddMapDeviceState(m *MapDeviceState) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetMapDeviceStateById retrieves MapDeviceState by Id. Returns error if
// Id doesn't exist
func GetMapDeviceStateById(id int64) (v *MapDeviceState, err error) {
	o := orm.NewOrm()
	v = &MapDeviceState{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllMapDeviceState retrieves all MapDeviceState matches certain condition. Returns empty list if
// no records exist
func GetAllMapDeviceState(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(MapDeviceState))
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

	var l []MapDeviceState
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

// UpdateMapDeviceState updates MapDeviceState by Id and returns error if
// the record to be updated doesn't exist
func UpdateMapDeviceStateById(m *MapDeviceState) (err error) {
	o := orm.NewOrm()
	v := MapDeviceState{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMapDeviceState deletes MapDeviceState by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMapDeviceState(id int64) (err error) {
	o := orm.NewOrm()
	v := MapDeviceState{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&MapDeviceState{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// AddMultipleMapDeviceState
// Use a prepared statement to increase inserting speed with multiple inserts.
func AddMultipleMapDeviceState(states []*MapDeviceState) (ids []int64, errs []error) {

	o := orm.NewOrm()
	qs := o.QueryTable(&MapDeviceState{})
	i, _ := qs.PrepareInsert()
	for _, state := range states {
		id, err := i.Insert(state)
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