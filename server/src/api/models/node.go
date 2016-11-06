package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"time"
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego"
)

type Node struct {
	Id    		int64  		`orm:"pk;auto;column(id)" json:"id"`
	Name 		string 		`orm:"size(254)" json:"name" valid:"MaxSize(254);Required"`
	Ip		string		`orm:"size(128)" json:"ip" valid:"IP;Required"`			// Must be a valid IPv4 address
	Port 		int 		`orm:"size(11)" json:"port" valid:"Range(1, 65535);Required"`
	Status	 	string 		`orm:"size(254)" json:"status"`
	Description 	string 		`orm:"type(longtext)" json:"description"`
	Created_at	time.Time	`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at	time.Time	`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
}

func (m *Node) TableName() string {
	return beego.AppConfig.String("db_node")
}

func init() {
	orm.RegisterModel(new(Node))
}

// AddNode insert a new Node into database and returns
// last inserted Id on success.
func AddNode(m *Node) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetNodeById retrieves Node by Id. Returns error if
// Id doesn't exist
func GetNodeById(id int64) (v *Node, err error) {
	o := orm.NewOrm()
	v = &Node{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllNode retrieves all Node matches certain condition. Returns empty list if
// no records exist
func GetAllNode(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Node))
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

	var l []Node
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
			"limit": limit,
			"offset": offset,
		}
		return ml, meta, nil
	}
	return nil, nil, err
}

// UpdateNode updates Node by Id and returns error if
// the record to be updated doesn't exist
func UpdateNodeById(m *Node) (err error) {
	o := orm.NewOrm()
	v := Node{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteNode deletes Node by Id and returns error if
// the record to be deleted doesn't exist
func DeleteNode(id int64) (err error) {
	o := orm.NewOrm()
	v := Node{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Node{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func (n *Node) Valid(v *validation.Validation)  {

	o := orm.NewOrm()
	nn := Node{Ip: n.Ip, Port: n.Port}
	o.Read(&nn, "Ip", "Port")

	if nn.Id != 0 {
		v.SetError("ip", "Not unique")
		v.SetError("port", "Not unique")
		return
	}

	return
}