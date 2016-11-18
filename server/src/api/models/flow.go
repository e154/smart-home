package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"encoding/json"
)

type Flow struct {
	Id   		int64  		`orm:"pk;auto;column(id)" json:"id"`
	Name		string		`orm:"" json:"name"`
	Status		string		`orm:"" json:"status"`
	WorkflowId	int64		`orm:"column(workflow_id)" json:"workflow_id"`
	Created_at	time.Time	`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at	time.Time	`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
	Connections	[]*Connection	`orm:"-" json:"connections"`
	FlowElements	[]*FlowElement	`orm:"-" json:"flow_elements"`
	Cursor		[]*FlowElement	`orm:"-" json:"-"`
	wf		*Workflow	`orm:"-" json:"-"`
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
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetFlowById retrieves Flow by Id. Returns error if
// Id doesn't exist
func GetFullFlowById(id int64) (v *Flow, err error) {
	o := orm.NewOrm()
	v = &Flow{Id: id}
	if err = o.Read(v); err == nil {
		err = FlowGetRelatedDate(v)
		return v, nil
	}

	return nil, err
}

// GetFlowById retrieves Flow by Id. Returns error if
// Id doesn't exist
func GetRedactorFlowById(id int64) (r *RedactorFlow, err error) {
	o := orm.NewOrm()
	v := &Flow{Id: id}
	if err = o.Read(v); err != nil {
		return nil, err
	}

	if r, err = v.ExportToRedactor(); err == nil {
		return r, nil
	}

	return nil, err
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
	_, err = o.QueryTable(&Flow{}).Filter("status", "enabled").All(&fs)

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

	for _, flow := range flows {
		FlowGetRelatedDate(flow)
	}

	return
}

func FlowGetRelatedDate(flow *Flow) (err error) {

	//_, err = o.LoadRelated(f, "Connections")
	//_, err = o.LoadRelated(f, "FlowElements")

	if flow.FlowElements, err = GetFlowElementsByFlow(flow); err != nil {
		return
	}

	if flow.Connections, err = GetConnectionsByFlow(flow); err != nil {
		return
	}

	for _, conn := range flow.Connections {
		for _, element := range flow.FlowElements {
			if element.Id == conn.ElementFrom {
				conn.FlowElementFrom = element
			} else if element.Id == conn.ElementTo {
				conn.FlowElementTo = element
			}
		}
	}

	for _, element := range flow.FlowElements {
		element.Flow = flow
		switch element.PrototypeType  {
		case "MessageHandler":
			element.Prototype = &MessageHandler{}
			break
		case "MessageEmitter":
			element.Prototype = &MessageEmitter{}
			break
		case "Task":
			element.Prototype = &Task{}
			break
		}
	}

	return
}

//---------------------------------------------------
// Workflow
//---------------------------------------------------
func (f *Flow) AddConnection(connection *Connection) {
	f.Connections = append(f.Connections, connection)
}

func (f *Flow) AddElement(flow_element *FlowElement) {
	flow_element.Flow = f
	f.FlowElements = append(f.FlowElements, flow_element)
}

func (f *Flow) NewMessage(message *Message) (err error) {

	var exist bool
	for _, element := range f.FlowElements {
		if element.Prototype == nil {
			continue
		}

		if element.PrototypeType != "MessageHandler" {
			continue
		}

		element.Run(message)

		exist = true
	}

	if !exist {
		err = errors.New("Message handler not found")
	}

	return
}

func (f *Flow) ExportToRedactor() (flow *RedactorFlow, err error) {

	flow = new(RedactorFlow)
	flow.Name = f.Name
	flow.Description = ""
	flow.Objects = make([]*RedactorObject, 0)
	flow.Connectors = make([]*RedactorConnector, 0)

	var flowElements []*FlowElement
	if flowElements, err = GetFlowElementsByFlow(f); err != nil {
		return
	}

	for _, el := range flowElements {
		object := &RedactorObject{
			Id: el.Id,
			Title: el.Name,
		}

		switch el.PrototypeType {
		case "MessageHandler":
			object.Type.Name = "event"
			object.Type.Start = map[int64]interface{}{0: &map[int64]interface{}{0: true}}
		case "MessageEmitter":
			object.Type.Name = "event"
			object.Type.End = map[string]interface{}{"simply": &map[string]interface{}{"top_level": true}}
		case "Task":
			object.Type.Name = "task"
		default:

		}

		gst := new(RedactorGrapSettings)
		if err = json.Unmarshal([]byte(el.GraphSettings), &gst); err != nil {
			return
		}

		object.Position = gst.Position

		flow.Objects = append(flow.Objects, object)
	}

	var connections []*Connection
	if connections, err = GetConnectionsByFlow(f); err != nil {
		return
	}

	for _, con := range connections {
		connector := &RedactorConnector{
			Id: con.Id,
			flow_type: "default",
			title: con.Name,
		}
		connector.Start.Object = con.ElementFrom
		connector.Start.Point = con.ElementFrom

		connector.End.Object = con.ElementTo
		connector.End.Point = con.PointTo

		flow.Connectors = append(flow.Connectors, connector)
	}

	return
}