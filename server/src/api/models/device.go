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

type Device struct {
	Id   		int64  		`orm:"pk;auto;column(id)" json:"id"`
	Device		*Device		`orm:"rel(fk);null" json:"device"`
	Node		*Node		`orm:"rel(fk);null" json:"node"`
	Address   	*int  		`orm:"" json:"address"`
	Baud		int		`orm:"size(11)" json:"baud"`
	Sleep		int		`orm:"size(32)" json:"sleep"`
	Description 	string 		`orm:"size(254)" json:"description" valid:"MaxSize(254)"`
	Name 		string 		`orm:"size(254)" json:"name" valid:"MaxSize(254);Required"`
	Status	 	string 		`orm:"size(254)" json:"status" valid:"MaxSize(254)"`
	StopBite   	int64  		`orm:"size(11)" json:"stop_bite"`
	Timeout   	time.Duration 	`orm:"" json:"timeout"`
	Tty 		string 		`orm:"size(254)" json:"tty" valid:"MaxSize(254)"`
	States		[]*DeviceState	`orm:"reverse(many)" json:"states"`
	Actions		[]*DeviceAction	`orm:"reverse(many)" json:"actions"`
	Created_at	time.Time	`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at	time.Time	`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
	IsGroup		bool		`orm:"-" json:"is_group"`
}

func (m *Device) TableName() string {
	return beego.AppConfig.String("db_devices")
}

func init() {
	orm.RegisterModel(new(Device))
}

// AddDevice insert a new Device into database and returns
// last inserted Id on success.
func AddDevice(m *Device) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetDeviceById retrieves Device by Id. Returns error if
// Id doesn't exist
func GetDeviceById(id int64) (v *Device, err error) {
	o := orm.NewOrm()
	v = &Device{Id: id}
	if err = o.Read(v); err != nil {
		return
	}

	if v.Device != nil {
		_, err = o.LoadRelated(v, "Device")
	}
	if v.Node != nil {
		_, err = o.LoadRelated(v, "Node")
	}

	return
}

// GetAllDevice retrieves all Device matches certain condition. Returns empty list if
// no records exist
func GetAllDevice(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Device))
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

	var l []Device
	qs = qs.RelatedSel("Device", "Node").OrderBy(sortFields...)
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

// UpdateDevice updates Device by Id and returns error if
// the record to be updated doesn't exist
func UpdateDeviceById(m *Device) (err error) {
	o := orm.NewOrm()
	v := Device{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDevice deletes Device by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDevice(id int64) (err error) {
	o := orm.NewOrm()
	v := Device{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Device{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetParentDeviceByChildId(id int64) (parent *Device, err error) {
	o := orm.NewOrm()
	err = o.Raw(fmt.Sprintf(`SELECT p2.*
		FROM %s as p1
		LEFT JOIn %[1]s as p2 on p2.id = p1.device_id
		WHERE p1.id = %d
		`, parent.TableName(), id)).QueryRow(&parent)

	parent.Id = id

	if parent.Device != nil {
		_, err = o.LoadRelated(parent, "Device")
	}

	if parent.Node != nil {
		_, err = o.LoadRelated(parent, "Node")
	}

	return
}

func (m *Device) GetChilds() (childs []*Device, all int64, err error) {
	o := orm.NewOrm()
	all, err = o.QueryTable(m).RelatedSel().Filter("device_id", m.Id).All(&childs)
	return
}

func (m *Device) GetNode() (*Node, error) {
	return GetNodeById(m.Node.Id)
}