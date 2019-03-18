package use_case

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
	m "github.com/e154/smart-home/models"
	"encoding/json"
	"time"
	"strings"
)

func AddLog(newLog *models.NewLog, adaptors *adaptors.Adaptors) (result *models.Log, errs []*validation.Error, err error) {

	log := &m.Log{}
	if err = common.Copy(&log, &newLog); err != nil {
		return
	}
	_, errs = log.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = adaptors.Log.Add(log); err != nil {
		return
	}

	if log, err = adaptors.Log.GetById(id); err != nil {
		return
	}

	result = &models.Log{}
	err = common.Copy(&result, &log)

	return
}
func GetLogById(logId int64, adaptors *adaptors.Adaptors) (logDto *models.Log, err error) {

	var log *m.Log
	if log, err = adaptors.Log.GetById(logId); err != nil {
		return
	}

	logDto = &models.Log{}
	err = common.Copy(&logDto, &log, common.JsonEngine)

	return
}

func GetLogList(limit, offset int64, order, sortBy, query string, adaptors *adaptors.Adaptors) (listDto []*models.Log, total int64, err error) {

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

	var list []*m.Log
	if list, total, err = adaptors.Log.List(limit, offset, order, sortBy, queryObj); err != nil {
		return
	}

	listDto = make([]*models.Log, 0)
	err = common.Copy(&listDto, &list)

	return
}

func SearchLog(query string, limit, offset int, adaptors *adaptors.Adaptors) (listDto []*models.Log, total int64, err error) {
	var list []*m.Log
	if list, total, err = adaptors.Log.Search(query, limit, offset); err != nil {
		return
	}

	listDto = make([]*models.Log, 0)
	err = common.Copy(&listDto, &list)
	return
}

func DeleteLogById(logId int64, adaptors *adaptors.Adaptors) (err error) {

	if _, err = adaptors.Log.GetById(logId); err != nil {
		return
	}

	err = adaptors.Log.Delete(logId)

	return
}
