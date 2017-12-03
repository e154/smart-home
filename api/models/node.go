package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego"
	"net/rpc"
	"net"
	"sync"
	"github.com/e154/smart-home/api/common"
)

type Node struct {
	Id    		int64  		`orm:"pk;auto;column(id)" json:"id"`
	Name 		string 		`orm:"size(254)" json:"name" valid:"MaxSize(254);Required"`
	Ip		string		`orm:"size(128)" json:"ip" valid:"IP;Required"`			// Must be a valid IPv4 address
	Port        	int 		`orm:"size(11)" json:"port" valid:"Range(1, 65535);Required"`
	Status      	string 		`orm:"size(254)" json:"status"`
	Description 	string 		`orm:"type(longtext)" json:"description"`
	Created_at  	time.Time	`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Update_at   	time.Time	`orm:"auto_now;type(datetime);column(update_at)" json:"update_at"`
	rpcClient   	*rpc.Client	`orm:"-" json:"-"`
	netConn     	net.Conn	`orm:"-" json:"-"`
	Errors      	int64		`orm:"-" json:"-"`
	connStatus  	string		`orm:"-" json:"-"`
	mutex       	sync.Mutex
}

func (m *Node) TableName() string {
	return beego.AppConfig.String("db_nodes")
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

	var count int64
	count, err = o.QueryTable(&Device{}).Filter("node_id", id).Count()
	if err != nil {
		return
	}

	if count > 0 {
		err = errors.New("node: Not delete with child devices!")
		return
	}

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

func GetNodesCount() (total int64, err error) {
	o := orm.NewOrm()
	total, err = o.QueryTable(&Node{}).Count()
	return
}

func GetAllEnabledNodes() (nodes []*Node, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(&Node{}).Filter("status", "enabled").All(&nodes)
	return
}

func (n *Node) RpcDial() (*rpc.Client, error) {
	var err error
	if _ , err = n.TcpDial(); err != nil {return nil, err}
	if n.rpcClient == nil { n.rpcClient = rpc.NewClient(n.netConn) }
	return n.rpcClient, err
}

func (n *Node) TcpDial() (net.Conn, error) {
	var err error
	if n.netConn == nil {
		n.netConn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", n.Ip, n.Port), time.Second * 2)
		if err != nil {
			return nil, err
		}
	}
	//defer n.netConn.Close()
	return n.netConn, err
}

func (n *Node) TcpClose() {
	if n.netConn == nil {
		return
	}
	n.netConn.Close()
	n.netConn = nil
	n.rpcClient = nil
}

func (n *Node) GetVersion() (version string, err error) {
	if n.rpcClient == nil {
		err = errors.New("rpc.client is nil")
		return
	}
	err = n.rpcClient.Call("Node.Version", "", &version)
	return
}

func (n *Node) GetConn() net.Conn {
	return n.netConn
}

func (n *Node) GetConnectStatus() string {
	if n.Status == "disabled" {
		n.connStatus = "disabled"
	}

	return n.connStatus
}

func (n *Node) SetConnectStatus(st string) {
	n.connStatus = st
}

func (n *Node) ModbusSend(device *Device, return_result bool, command []byte, reply interface{}) error {

	if n.rpcClient == nil {
		return errors.New("rpc.client is nil")
	}

	if n.netConn == nil {
		n.Errors++
		return errors.New("node not connected")
	}

	request := &common.Request{
		Baud: device.Baud,
		Device: device.Tty,
		Timeout: device.Timeout,
		StopBits: device.StopBite,
		Sleep: device.Sleep,
		Result: return_result,
		Command: command,
	}

	if err := n.rpcClient.Call("Modbus.Send", request, reply); err != nil {
		n.Errors++
		return err
	}

	return nil
}

func (n *Node) Send(protocol string, device *Device, return_result bool, command []byte) (result common.Result) {
	switch protocol {
	case "modbus":
		if err := n.ModbusSend(device, return_result, command, &result); err != nil {
			result.Error = err.Error()
		}
	}
	return
}

func (n *Node) RpcCall(method string, request interface{}, reply interface{}) error {

	if n.rpcClient == nil {
		return errors.New("rpc.client is nil")
	}

	if n.netConn == nil {
		n.Errors++
		return errors.New("node not connected")
	}

	if err := n.rpcClient.Call(method, request, reply); err != nil {
		n.Errors++
		return err
	}

	return nil
}