package dto

import (
	"github.com/e154/smart-home/api/stub/api"
	m "github.com/e154/smart-home/models"
)

func GetStatistic(statistic []*m.Statistic) (result *api.Statistics) {
	result = &api.Statistics{
		Items: make([]*api.Statistic, 0, len(statistic)),
	}
	for _, item := range statistic {
		result.Items = append(result.Items, &api.Statistic{
			Name:        item.Name,
			Description: item.Description,
			Value:       item.Value,
			Diff:        item.Diff,
		})
	}
	return
}
