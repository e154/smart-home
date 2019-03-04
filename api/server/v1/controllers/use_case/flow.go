package use_case

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

func GetFlowList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (items []*m.Flow, total int64, err error) {

	items, total, err = adaptors.Flow.List(limit, offset, order, sortBy)

	return
}