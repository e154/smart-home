package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type Connection struct {
	Uuid   		string  	`orm:"pk" json:"uuid"`
	Name		string		`orm:"" json:"name"`
	ElementFrom	string		`orm:"column(element_from);type(string)" json:"element_from"`
	ElementTo	string		`orm:"column(element_to);type(string)" json:"element_to"`
	PointFrom	int64		`orm:"column(point_from)" json:"point_from"`
	PointTo		int64		`orm:"column(point_to)" json:"point_to"`
	FlowId		int64		`orm:"column(flow_id)" json:"flow_id"`
	GraphSettings	string		`orm:"column(graph_settings)" json:"graph_settings"`
	Created_at	time.Time	`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at	time.Time	`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
	FlowElementFrom	*FlowElement	`orm:"-" json:"flow_elemet_from"`
	FlowElementTo	*FlowElement	`orm:"-" json:"flow_elemet_to"`
	Flow		*Flow		`orm:"-" json:"flow"`
}

func (m *Connection) TableName() string {
	return beego.AppConfig.String("db_connections")
}

func init() {
	orm.RegisterModel(new(Connection))
}

// AddConnection insert a new Connection into database and returns
// last inserted Id on success.
func AddConnection(m *Connection) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddFlowElement insert a new FlowElement into database and returns
// last inserted Id on success.
func AddOrUpdateConnection(m *Connection) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.InsertOrUpdate(m, "uuid")
	return
}

// GetConnectionById retrieves Connection by Id. Returns error if
// Id doesn't exist
func GetConnectionById(id string) (v *Connection, err error) {
	o := orm.NewOrm()
	v = &Connection{Uuid: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllConnection retrieves all Connection matches certain condition. Returns empty list if
// no records exist
func GetAllConnection(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Connection))
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

	var l []Connection
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

// UpdateConnection updates Connection by Id and returns error if
// the record to be updated doesn't exist
func UpdateConnectionById(m *Connection) (err error) {
	o := orm.NewOrm()
	v := Connection{Uuid: m.Uuid}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteConnection deletes Connection by Id and returns error if
// the record to be deleted doesn't exist
func DeleteConnection(uuid string) (err error) {
	o := orm.NewOrm()
	v := Connection{Uuid: uuid}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Connection{Uuid: uuid}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}


func GetConnectionsByFlow(flow *Flow) (connections []*Connection, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(&Connection{}).Filter("flow_id", flow.Id).All(&connections)

	return
}
