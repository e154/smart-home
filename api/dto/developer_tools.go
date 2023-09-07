package dto

import (
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/system/bus"
)

type DeveloperTools struct {
}

func NewDeveloperToolsDto() DeveloperTools {
	return DeveloperTools{}
}

func (DeveloperTools) GetEventBusState(state bus.Stats, total int64) (result *api.EventBusStateListResult) {
	result = &api.EventBusStateListResult{
		Items: make([]*api.BusStateItem, 0, len(state)),
		Meta: &api.Meta{
			Limit: uint64(total),
			Page:  1,
			Total: uint64(total),
		},
	}
	for _, item := range state {
		result.Items = append(result.Items, &api.BusStateItem{
			Topic:       item.Topic,
			Subscribers: int32(item.Subscribers),
		})
	}
	return
}
