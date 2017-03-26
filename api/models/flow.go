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

type Flow struct {
	Id          	int64  				`orm:"pk;auto;column(id)" json:"id"`
	Name        	string				`orm:"" json:"name"`
	Description 	string				`orm:"" json:"description"`
	Status      	string				`orm:"" json:"status"`
	Workflow    	*Workflow			`orm:"rel(fk)" json:"workflow"`
	Scenario    	*WorkflowScenario		`orm:"rel(fk);column(workflow_scenario_id)" json:"scenario"`
	Created_at  	time.Time			`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at   	time.Time			`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
	Connections 	[]*Connection			`orm:"-" json:"connections"`
	FlowElements	[]*FlowElement			`orm:"-" json:"flow_elements"`
	Workers     	[]*Worker			`orm:"-" json:"workers"`
}

func (m *Flow) TableName() string {
	return beego.AppConfig.String("db_flows")
}

func init() {
	orm.RegisterModel(new(Flow))
}

// AddFlow insert a new Flow into database and returns
// last inserted Id on success.
func AddFlow(m *Flow) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetFlowById retrieves Flow by Id. Returns error if
// Id doesn't exist
func GetFlowById(id int64) (v *Flow, err error) {
	o := orm.NewOrm()
	v = &Flow{Id: id}
	if err = o.Read(v); err != nil {
		return
	}

	if v.Workflow != nil {
		_, err = o.LoadRelated(v, "Workflow")
	}

	return
}

// GetFlowById retrieves Flow by Id. Returns error if
// Id doesn't exist
func GetFullFlowById(id int64) (v *Flow, err error) {
	o := orm.NewOrm()
	v = &Flow{Id: id}
	if err = o.Read(v); err != nil {
		return
	}

	err = FlowGetRelatedDate(v)

	return
}

// GetAllFlow retrieves all Flow matches certain condition. Returns empty list if
// no records exist
func GetAllFlow(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Flow))
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

	var l []Flow
	qs = qs.RelatedSel("Workflow").OrderBy(sortFields...)
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

// UpdateFlow updates Flow by Id and returns error if
// the record to be updated doesn't exist
func UpdateFlowById(m *Flow) (err error) {
	o := orm.NewOrm()
	v := Flow{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteFlow deletes Flow by Id and returns error if
// the record to be deleted doesn't exist
func DeleteFlow(id int64) (err error) {
	o := orm.NewOrm()
	v := Flow{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Flow{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAllEnabledFlows() (fs []*Flow, err error) {
	o := orm.NewOrm()
	fs = []*Flow{}
	_, err = o.QueryTable(&Flow{}).Filter("status", "enabled").RelatedSel("Workflow").All(&fs)

	return
}

func GetEnabledFlowById(id int64) (flow *Flow, err error) {
	o := orm.NewOrm()
	flow = new(Flow)
	err = o.QueryTable(&Flow{}).Filter("id", id).Filter("status", "enabled").One(flow)
	FlowGetRelatedDate(flow)
	return
}

func GetAllEnabledFlowsByWf(wf *Workflow) (flows []*Flow, err error) {
	o := orm.NewOrm()
	flows = []*Flow{}
	_, err = o.QueryTable(&Flow{}).Filter("status", "enabled").Filter("workflow_id", wf.Id).All(&flows)
	return
}

func FlowGetRelatedDate(flow *Flow) (err error) {

	if flow.FlowElements, err = GetFlowElementsByFlow(flow); err != nil {
		return
	}

	if flow.Connections, err = GetConnectionsByFlow(flow); err != nil {
		return
	}

	for _, conn := range flow.Connections {
		for _, element := range flow.FlowElements {
			if element.Uuid == conn.ElementFrom {
				conn.FlowElementFrom = element
			} else if element.Uuid == conn.ElementTo {
				conn.FlowElementTo = element
			}
		}
	}

	return
}

func (f *Flow) AddConnection(connection *Connection) {
	f.Connections = append(f.Connections, connection)
}

func (f *Flow) GetAllEnabledWorkers() (workers []*Worker, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(&Worker{}).Filter("flow_id", f.Id).Filter( "status", "enabled").RelatedSel().All(&workers)
	return
}

func (f *Flow) GetWorkers() (workers []*Worker, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(&Worker{}).Filter("flow_id", f.Id).All(&workers)
	return
}

func (f *Flow) GetScenario() (int64, error) {
	if f.Scenario != nil {
		o := orm.NewOrm()
		return  o.LoadRelated(f, "Scenario")
	}

	return 0, nil
}