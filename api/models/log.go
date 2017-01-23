package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type Log struct {
	Id    		int64  		`orm:"pk;auto;column(id)" json:"id"`
	Body 		string 		`orm:"" json:"body"`
	Level		string		`orm:"" json:"level"`
	Created_at	time.Time	`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
}

func (m *Log) TableName() string {
	return beego.AppConfig.String("db_logs")
}

func init() {
	orm.RegisterModel(new(Log))
}

// AddLog insert a new Log into database and returns
// last inserted Id on success.
func AddLog(m *Log) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetLogById retrieves Log by Id. Returns error if
// Id doesn't exist
func GetLogById(id int64) (v *Log, err error) {
	o := orm.NewOrm()
	v = &Log{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLog retrieves all Log matches certain condition. Returns empty list if
// no records exist
func GetAllLog(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (logs []Log, meta *map[string]int64, err error) {

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").From(new(Log).TableName())
	var start_date, end_date, levels string
	var ok bool

	if start_date, ok = query["start_date"]; ok {
		qb.Where(fmt.Sprintf("Date(logs.created_at) >= '%s'", start_date))
	}

	if end_date, ok = query["end_date"]; ok {
		if start_date != "" {
			qb.And(fmt.Sprintf("Date(logs.created_at) <= '%s'", end_date))
		} else {
			qb.Where(fmt.Sprintf("Date(logs.created_at) <= '%s'", end_date))
		}
	}

	if levels, ok = query["levels"]; ok {
		if start_date != "" || end_date != "" {
			qb.And(fmt.Sprintf("logs.level IN (%s)", levels))
		} else {
			qb.Where(fmt.Sprintf("logs.level IN (%s)", levels))
		}

	}

	o := orm.NewOrm()
	objects_count, _ := o.Raw(qb.String()).QueryRows(&logs)
	if len(sortby) > 0 && len(order) > 0 {
		qb.OrderBy(fmt.Sprintf("logs.%s",sortby[0]))
		if order[0] == "desc" {
			qb.Desc()
		} else {
			qb.Asc()
		}
	}

	qb.Limit(int(limit)).Offset(int(offset))
	o.Raw(qb.String()).QueryRows(&logs)
	meta = &map[string]int64{
		"objects_count": objects_count,
		"limit":         limit,
		"offset":        offset,
	}

	return
}

// UpdateLog updates Log by Id and returns error if
// the record to be updated doesn't exist
func UpdateLogById(m *Log) (err error) {
	o := orm.NewOrm()
	v := Log{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLog deletes Log by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLog(id int64) (err error) {
	o := orm.NewOrm()
	v := Log{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Log{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
