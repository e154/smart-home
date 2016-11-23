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

type DeviceAction struct {
	Id   		int64  		`orm:"pk;auto;column(id)" json:"id"`
	StartAddr   	int64  		`orm:"column(start_addr)" json:"start_addr"`
	ColCells	int64		`orm:"column(col_cells)" json:"col_cells" valid:"Required"`
	Device		*Device 	`orm:"rel(fk)" json:"device"`
	Function   	int64  		`orm:"size(11)" json:"function"`
	Command 	string 		`orm:"" json:"command" valid:"Required"`
	Name 		string 		`orm:"size(254)" json:"name" valid:"MaxSize(254);Required"`
	Direction	string 		`orm:"size(254)" json:"direction" valid:"MaxSize(254);Required"`
	Description 	string 		`orm:"" json:"description"`
	ResultType 	string 		`orm:"" json:"result_type"`
	Created_at	time.Time	`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at	time.Time	`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
}

func (m *DeviceAction) TableName() string {
	return beego.AppConfig.String("db_device_actions")
}

func init() {
	orm.RegisterModel(new(DeviceAction))
}

// AddDeviceAction insert a new DeviceAction into database and returns
// last inserted Id on success.
func AddDeviceAction(m *DeviceAction) (id int64, err error) {
	o := orm.NewOrm()

	device := &Device{}
	if device, err = GetDeviceById(m.Device.Id); err != nil {
		return
	}

	if device.Device != nil {
		m.Device = device.Device
	}

	id, err = o.Insert(m)
	return
}

// GetDeviceActionById retrieves DeviceAction by Id. Returns error if
// Id doesn't exist
func GetDeviceActionById(id int64) (v *DeviceAction, err error) {
	o := orm.NewOrm()
	v = &DeviceAction{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllDeviceAction retrieves all DeviceAction matches certain condition. Returns empty list if
// no records exist
func GetAllDeviceAction(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(DeviceAction))
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

	var l []DeviceAction
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

// UpdateDeviceAction updates DeviceAction by Id and returns error if
// the record to be updated doesn't exist
func UpdateDeviceActionById(m *DeviceAction) (err error) {
	o := orm.NewOrm()
	v := DeviceAction{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDeviceAction deletes DeviceAction by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDeviceAction(id int64) (err error) {
	o := orm.NewOrm()
	v := DeviceAction{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&DeviceAction{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetDeviceActionsByDeviceId(ids []int64) (actions []*DeviceAction, err error) {

	o := orm.NewOrm()

	actions = []*DeviceAction{}
	_, err = o.QueryTable(&DeviceAction{}).Filter("device_id__in", ids).RelatedSel("Device").All(&actions)
	return
}