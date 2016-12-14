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

	//TODO refactor!!!
	q := `SELECT *
FROM logs
`

	start_date, ok := query["start_date"];
	if  ok {
		q += `
		WHERE Date(logs.created_at) >= '`+start_date+`'`
	}

	end_date, ok := query["end_date"];
	if  ok {
		if start_date != "" {
			q += `
			AND`
		} else {
			q += `
			WHERE`
		}

		q += `
		Date(logs.created_at) <= '`+end_date+`'`
	}

	if levels, ok := query["levels"]; ok {
		if start_date != "" || end_date != "" {
			q += `
			AND`
		} else {
			q += `
			WHERE`
		}

		q += `
		logs.level IN (`+levels+`)`
	}

	o := orm.NewOrm()
	objects_count, _ := o.Raw(fmt.Sprint(q)).QueryRows(&logs)

	if len(sortby) > 0 && len(order) > 0 {

		q += `
		ORDER BY logs.`+ sortby[0] +` ` + order[0]
	}

	q += `
LIMIT %d
OFFSET %d`

	o.Raw(fmt.Sprintf(q, limit, offset)).QueryRows(&logs)
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
