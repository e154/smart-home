package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"encoding/json"
	"time"
	"strings"
)

type LogCommand struct {
	*CommonCommand
}

func NewLogCommand(common *CommonCommand) *LogCommand {
	return &LogCommand{
		CommonCommand: common,
	}
}

func (l *LogCommand) Add(log *m.Log) (result *m.Log, errs []*validation.Error, err error) {

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

func (l *LogCommand) GetById(id int64) (log *m.Log, err error) {

	log, err = l.adaptors.Log.GetById(id)

	return
}

func (l *LogCommand) GetList(limit, offset int64, order, sortBy, query string) (list []*m.Log, total int64, err error) {

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

func (l *LogCommand) Search(query string, limit, offset int) (list []*m.Log, total int64, err error) {

	list, total, err = l.adaptors.Log.Search(query, limit, offset)

	return
}

func (l *LogCommand) Delete(logId int64) (err error) {

	if _, err = l.adaptors.Log.GetById(logId); err != nil {
		return
	}

	err = l.adaptors.Log.Delete(logId)

	return
}
