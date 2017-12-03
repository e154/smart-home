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

type WorkflowScenario struct {
	Id    		int64 				`orm:"auto" json:"id"`
	Name 		string				`orm:"size(255)" json:"name"`
	Workflow	*Workflow			`orm:"rel(fk)" json:"workflow"`
	Scripts		[]*Script			`orm:"rel(m2m);rel_through(github.com/e154/smart-home/api/models.WorkflowScenarioScript)" json:"scripts"`
	SystemName  	string				`orm:"size(255)" json:"system_name"`
	Created_at	time.Time			`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at	time.Time			`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
}

func (m *WorkflowScenario) TableName() string {
	return beego.AppConfig.String("db_workflow_scenarios")
}

// AddScenario insert a new Scenario into database and returns
// last inserted Id on success.
func AddWorkflowScenario(m *WorkflowScenario) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetScenarioById retrieves Scenario by Id. Returns error if
// Id doesn't exist
func GetWorkflowScenarioById(id int64) (v *WorkflowScenario, err error) {
	o := orm.NewOrm()
	v = &WorkflowScenario{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllScenario retrieves all Scenario matches certain condition. Returns empty list if
// no records exist
func GetAllWorkflowScenario(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(WorkflowScenario))
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

	var l []WorkflowScenario
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

// UpdateScenario updates Scenario by Id and returns error if
// the record to be updated doesn't exist
func UpdateWorkflowScenarioById(m *WorkflowScenario) (err error) {
	o := orm.NewOrm()
	v := WorkflowScenario{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteScenario deletes Scenario by Id and returns error if
// the record to be deleted doesn't exist
func DeleteWorkflowScenario(id int64) (err error) {
	o := orm.NewOrm()
	v := WorkflowScenario{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&WorkflowScenario{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func (ws *WorkflowScenario) GetScripts() (int64, error) {

	o := orm.NewOrm()
	return  o.LoadRelated(ws, "Scripts", 3)

}

func (wf *WorkflowScenario) AddScripts(scripts []*Script) (num int64, err error) {
	if len(scripts) == 0 {
		return
	}

	o := orm.NewOrm()
	m2m := o.QueryM2M(wf, "Scripts")
	num, err = m2m.Add(scripts)
	if err == nil {
		fmt.Println("Added nums: ", num)
	}

	return
}

func (wf *WorkflowScenario) RemoveScripts(scripts []*Script) (num int64, err error) {
	if len(scripts) == 0 {
		return
	}

	o := orm.NewOrm()
	m2m := o.QueryM2M(wf, "Scripts")
	num, err = m2m.Remove(scripts)
	if err == nil {
		fmt.Println("Removed nums: ", num)
	}

	return
}

func (wf *WorkflowScenario) UpdateScripts(scripts []*Script) (num int64, err error) {
	var add, rem []*Script
	var exist bool

	for _, s1 := range wf.Scripts {
		exist = false
		for _, s2 := range scripts {
			if s1.Id == s2.Id {
				exist = true
				break
			}
		}
		if !exist {
			rem = append(rem, s1)
		}
	}

	for _, s1 := range scripts {
		exist = false
		for _, s2 := range wf.Scripts {
			if s1.Id == s2.Id {
				exist = true
				break
			}
		}
		if !exist {
			add = append(add, s1)
		}
	}

	if _, err = wf.RemoveScripts(rem); err != nil {
		return
	}

	if _, err = wf.AddScripts(add); err != nil {
		return
	}

	return
}