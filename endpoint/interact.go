package endpoint

import (
	"context"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/event_bus/events"
	"github.com/go-playground/validator/v10"
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
func (d InteractEndpoint) EntityCallAction(ctx context.Context, entityId string, action string, args map[string]interface{}) (errs validator.ValidationErrorsTranslations, err error) {

	id := common.EntityId(entityId)
	_, err = d.adaptors.Entity.GetById(id)
	if err != nil {
		return
	}

	d.eventBus.Publish(event_bus.TopicEntities, events.EventCallAction{
		PluginName: id.PluginName(),
		EntityId:   id,
		ActionName: action,
		Args:       args,
	})

	return
}
