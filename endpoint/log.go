package endpoint

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"encoding/json"
	"time"
	"strings"
)

type LogEndpoint struct {
	*CommonEndpoint
}

func NewLogEndpoint(common *CommonEndpoint) *LogEndpoint {
	return &LogEndpoint{
		CommonEndpoint: common,
	}
}

func (l *LogEndpoint) Add(log *m.Log) (result *m.Log, errs []*validation.Error, err error) {

	_, errs = log.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = l.adaptors.Log.Add(log); err != nil {
		return
	}

	log, err = l.adaptors.Log.GetById(id)

	return
}

func (l *LogEndpoint) GetById(id int64) (log *m.Log, err error) {

	log, err = l.adaptors.Log.GetById(id)

	return
}

func (l *LogEndpoint) GetList(limit, offset int64, order, sortBy, query string) (list []*m.Log, total int64, err error) {

	var queryObj *m.LogQuery
	if query != "" {
		queryObj = &m.LogQuery{}
		d := make(map[string]string, 0)
		if err = json.Unmarshal([]byte(query), &d); err != nil {
			return
		}

		if startDate, ok := d["start_date"]; ok {
			date, _ := time.Parse("2006-01-02", startDate)
			queryObj.StartDate = &date
		}
		if endDate, ok := d["end_date"]; ok {
			date, _ := time.Parse("2006-01-02", endDate)
			queryObj.EndDate = &date
		}
		if levels, ok := d["levels"]; ok {
			queryObj.Levels = strings.Split(strings.Replace(levels, "'", "", -1), ",")
		}
	}

	list, total, err = l.adaptors.Log.List(limit, offset, order, sortBy, queryObj)

	return
}

func (l *LogEndpoint) Search(query string, limit, offset int) (list []*m.Log, total int64, err error) {

	list, total, err = l.adaptors.Log.Search(query, limit, offset)

	return
}

func (l *LogEndpoint) Delete(logId int64) (err error) {

	if _, err = l.adaptors.Log.GetById(logId); err != nil {
		return
	}

	err = l.adaptors.Log.Delete(logId)

	return
}
