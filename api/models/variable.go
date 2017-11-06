package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type Variable struct {
	Name			string		`orm:"pk;size(128);unique;column(name)" json:"name"`
	Value			string		`orm:"column(value)" json:"value"`
	Autoload		string		`orm:"default(yes);column(autoload)" json:"autoload"`
}

func (m *Variable) TableName() string {
	return beego.AppConfig.String("db_variables")
}

func init() {
	orm.RegisterModel(new(Variable))
}

// GetVariableByName retrieves Variable by Name. Returns error if
// Name doesn't exist
func GetVariableByName(name string) (v *Variable, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(&Variable{}).Filter("name", name).One(v)
	return
}

// GetAllVariable retrieves all Variable matches certain condition. Returns empty list if
// no records exist
func GetAllVariable() (variables []*Variable, err error) {
	o := orm.NewOrm()
	variables = []*Variable{}
	_, err = o.QueryTable(&Variable{}).Filter("autoload", "yes").All(&variables)

	return
}

// InsertOrUpdateVariable updates Variable by Name and returns error if
// the record to be updated doesn't exist
func InsertOrUpdateVariableByName(m *Variable) (err error) {
	o := orm.NewOrm()
	//v := Variable{Name: m.Name}
	// ascertain name exists in the database
	//if err = o.Read(&v, "name"); err == nil {
	var num int64
	if num, err = o.InsertOrUpdate(m, "name"); err == nil {
		fmt.Println("Number of records updated in database:", num)
	}
	//}
	return
}

// DeleteVariable deletes Variable by Name and returns error if
// the record to be deleted doesn't exist
func DeleteVariable(name string) (err error) {
	o := orm.NewOrm()
	v := Variable{Name: name}
	// ascertain name exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Variable{Name: name}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
