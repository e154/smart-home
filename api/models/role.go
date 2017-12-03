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

type Role struct {
	Name		string			`orm:"pk;unique;size(255);index" valid:"Required;MaxSize(255)" json:"name"`
	Description	string			`orm:"size(255)" json:"description"`
	Parent		*Role			`orm:"rel(fk);column(parent);null" json:"parent"`
	Children	[]*Role			`orm:"reverse(many)" json:"children"`
	Permissions	[]*Permission		`orm:"reverse(many)" json:"-"`
	AccessList	map[string][]string	`orm:"-" json:"access_list"`
	Created_at	time.Time		`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at	time.Time		`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
}

func (m *Role) TableName() string {
	return beego.AppConfig.String("db_roles")
}

// AddRole insert a new Role into database and returns
// last inserted Name on success.
func AddRole(m *Role) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRoleById retrieves Role by Name. Returns error if
// Name doesn't exist
func GetRoleByName(name string) (v *Role, err error) {
	o := orm.NewOrm()
	v = &Role{Name: name}
	if err = o.Read(v, "Name"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRole retrieves all Role matches certain condition. Returns empty list if
// no records exist
func GetAllRole(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Role))
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

	var l []Role
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

// UpdateRole updates Role by Name and returns error if
// the record to be updated doesn't exist
func UpdateRoleByName(m *Role) (err error) {
	o := orm.NewOrm()
	v := Role{Name: m.Name}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRole deletes Role by Name and returns error if
// the record to be deleted doesn't exist
func DeleteRole(name string) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(&Role{}).Filter("name", name).Delete()
	return
}

func (r *Role) LoadRelated() (err error) {
	o := orm.NewOrm()

	_, err = o.LoadRelated(r, "Parent", 3)
	_, err = o.LoadRelated(r, "Children", 3)

	r.GetAccessList()

	return
}

func (r *Role) GetFullAccessList() (access_list AccessList) {

	var item AccessItem
	var levels AccessLevels
	var ok bool
	roleList := []*Role{}
	roles := []*Role{}
	access_list = NewAccessList()

	// получим полный список ролей
	o := orm.NewOrm()
	if _, err := o.QueryTable(r).All(&roles); err != nil {
		return
	}

	getRoleParentList(&roles, r.Name, &roleList)

	//
	for _, role := range roleList {
		if _, err := o.LoadRelated(role, "Permissions"); err != nil {
			break
		}

		if role.Permissions == nil {
			continue
		}

		for _, perm := range role.Permissions {
			if levels, ok = AccessConfigList[perm.PackageName]; !ok {
				continue
			}

			if item, ok = levels[perm.LevelName]; !ok {
				continue
			}


			if access_list[perm.PackageName] == nil {
				access_list[perm.PackageName] = NewAccessLevels()
			}

			//fmt.Println(perm.PackageName, item)
			item.RoleName = role.Name
			access_list[perm.PackageName][perm.LevelName] = item
		}
	}

	//fmt.Println(permissions)

	return
}

func (r *Role) GetAccessList() {
	access_list := r.GetFullAccessList()

	r.AccessList = make(map[string][]string)
	for pack_name, pack := range access_list {
		levels := []string{}
		for level_name, _ := range pack {
			levels = append(levels, level_name)
		}

		r.AccessList[pack_name] = levels
	}
}

func (r *Role) UpdateAccessList(access_list map[string]map[string]bool) (err error) {


	o := orm.NewOrm()
	add_perms := []*Permission{}
	del_perms := []string{}
	for pack_name, pack := range access_list {
		for level_name, dir := range pack {
			if dir {
				add_perms = append(add_perms, &Permission{
					Role: r,
					PackageName: pack_name,
					LevelName: level_name,
				})
			} else {
				del_perms = append(del_perms, level_name)
			}

			if len(del_perms) > 0 {
				if _, err = o.QueryTable(&Permission{}).Filter("package_name", pack_name).Filter("level_name__in", del_perms).Delete(); err != nil {
					return
				}
				del_perms = []string{}
			}
		}

	}

	if len(add_perms) == 0 {
		return
	}

	for _, perm := range add_perms {
		if _, _, err = o.ReadOrCreate(perm, "PackageName", "LevelName", "Role"); err != nil {
			fmt.Println(err.Error())
		}
	}

	return
}

// Поиск всех родителей по ребенку
func getRoleParentList(items *[]*Role, name string, list *[]*Role) {

	PARENT:
	for _, item := range *items {
		if item.Name == name {
			*list = append(*list, item)
			name = ""
			if item.Parent != nil {
				name = item.Parent.Name
			}

			goto PARENT
			break
		}

	}
}