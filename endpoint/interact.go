package endpoint

import (
	"context"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
)

// InteractEndpoint ...
type InteractEndpoint struct {
	*CommonEndpoint
}

// NewInteractEndpoint ...
func NewInteractEndpoint(common *CommonEndpoint) *InteractEndpoint {
	return &InteractEndpoint{
		CommonEndpoint: common,
	}
}

// EntityCallAction ...
func (d InteractEndpoint) EntityCallAction(ctx context.Context, entityId string, action string, args map[string]interface{}) (errs map[string]string, err error) {

	id := common.EntityId(entityId)
	_, err = d.adaptors.Entity.GetById(ctx, id)
	if err != nil {
		return
	}

	d.eventBus.Publish("system/entities/"+id.String(), events.EventCallEntityAction{
		PluginName: id.PluginName(),
		EntityId:   id,
		ActionName: action,
		Args:       args,
	})

	return
}
