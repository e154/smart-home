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

type MessageDeliverie struct {
	Id    			int64  		`orm:"pk;auto;column(id)" json:"id"`
	Message    		*Message	`orm:"rel(fk)" json:"message"`
	State 			string 		`orm:"size(254)" json:"state"`
	Address			string		`orm:"" json:"address"`
	Error_system_code 	string 		`orm:"size(254)" json:"error_system_code"`
	Error_system_message 	string 		`orm:"size(254)" json:"error_system_message"`
	Created_at  		time.Time	`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at   		time.Time	`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
}

func (m *MessageDeliverie) TableName() string {
	return beego.AppConfig.String("db_message_deliveries")
}

// AddMessageDeliverie insert a new MessageDeliverie into database and returns
// last inserted Id on success.
func AddMessageDeliverie(m *MessageDeliverie) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetMessageDeliverieById retrieves MessageDeliverie by Id. Returns error if
// Id doesn't exist
func GetMessageDeliverieById(id int64) (v *MessageDeliverie, err error) {
	o := orm.NewOrm()
	v = &MessageDeliverie{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllMessageDeliverie retrieves all MessageDeliverie matches certain condition. Returns empty list if
// no records exist
func GetAllMessageDeliverie(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(MessageDeliverie))
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

	var l []MessageDeliverie
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

// UpdateMessageDeliverie updates MessageDeliverie by Id and returns error if
// the record to be updated doesn't exist
func UpdateMessageDeliverieById(m *MessageDeliverie) (err error) {
	o := orm.NewOrm()
	v := MessageDeliverie{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMessageDeliverie deletes MessageDeliverie by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMessageDeliverie(id int64) (err error) {
	o := orm.NewOrm()
	v := MessageDeliverie{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&MessageDeliverie{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAllMessageDeliveriesInProgress() (messages []*MessageDeliverie, count int64, err error) {

	o := orm.NewOrm()
	count, err = o.QueryTable(&MessageDeliverie{}).RelatedSel().Filter("state", "in_progress").All(&messages)

	return
}

func AddMessageDeliverieMultiple(mds []*MessageDeliverie) (count int64, _errors []error) {

	o := orm.NewOrm()
	qs := o.QueryTable(&MessageDeliverie{})
	i, _ := qs.PrepareInsert()

	for _, m := range mds {
		if _, err := i.Insert(m); err != nil {
			_errors = append(_errors, err)
		}
		count++
	}

	i.Close()

	return
}