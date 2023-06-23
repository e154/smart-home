package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// EntityStorage ...
type EntityStorage struct{}

// NewEntityStorageDto ...
func NewEntityStorageDto() EntityStorage {
	return EntityStorage{}
}

func (_ EntityStorage) List(list []*m.EntityStorage, total uint64, pagination common.PageParams, entity *m.Entity) (result *api.GetEntityStorageResult) {

	var items = make([]*api.EntityStorage, 0, len(list))

	for _, item := range list {
		entity.Attributes.Deserialize(item.Attributes)
		items = append(items, &api.EntityStorage{
			Id:         item.Id,
			EntityId:   string(item.EntityId),
			State:      item.State,
			Attributes: AttributeToApi(entity.Attributes),
			CreatedAt:  timestamppb.New(item.CreatedAt),
		})
	}

	return &api.GetEntityStorageResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}
