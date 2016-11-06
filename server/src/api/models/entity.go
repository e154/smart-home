package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"time"
	"github.com/astaxie/beego"
)

type Entity struct {
	Id   		int64  		`orm:"pk;auto;column(id)" json:"id"`
	EntityId   	int64  		`orm:"size(11)" json:"entity_id"`
	NodeId   	int64  		`orm:"size(11)" json:"node_id"`
	StopBite   	int64  		`orm:"size(11)" json:"stop_bite"`
	Address   	int  		`orm:"" json:"address"`
	Timeout   	time.Duration 	`orm:"" json:"timeout"`
	Status	 	string 		`orm:"size(254)" json:"status" valid:"MaxSize(254)"`
	Name 		string 		`orm:"size(254)" json:"name" valid:"MaxSize(254);Required"`
	Device 		string 		`orm:"size(254)" json:"device" valid:"MaxSize(254);Required"`
	Description 	string 		`orm:"size(254)" json:"description" valid:"MaxSize(254)"`
	Baud		int		`orm:"" json:"baud" valid:"Range(1, 115200);Required"`
	Created_at	time.Time	`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at	time.Time	`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
}

func (m *Entity) TableName() string {
	return beego.AppConfig.String("db_nodes")
}

func init() {
	orm.RegisterModel(new(Entity))
}

// AddEntity insert a new Entity into database and returns
// last inserted Id on success.
func AddEntity(m *Entity) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetEntityById retrieves Entity by Id. Returns error if
// Id doesn't exist
func GetEntityById(id int64) (v *Entity, err error) {
	o := orm.NewOrm()
	v = &Entity{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllEntity retrieves all Entity matches certain condition. Returns empty list if
// no records exist
func GetAllEntity(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Entity))
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

	var l []Entity
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

// UpdateEntity updates Entity by Id and returns error if
// the record to be updated doesn't exist
func UpdateEntityById(m *Entity) (err error) {
	o := orm.NewOrm()
	v := Entity{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteEntity deletes Entity by Id and returns error if
// the record to be deleted doesn't exist
func DeleteEntity(id int64) (err error) {
	o := orm.NewOrm()
	v := Entity{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Entity{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
