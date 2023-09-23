package neural_network

import (
	"context"
	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.neural_network")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
}

func New() supervisor.Pluggable {
	return &plugin{
		Plugin: supervisor.NewPlugin(),
	}
}

func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	_ = p.Service.EventBus().Subscribe("system/entities/+", p.eventHandler)

	return nil
}

func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	_ = p.Service.EventBus().Unsubscribe("system/entities/+", p.eventHandler)

	return nil
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor supervisor.PluginActor, err error) {
	actor = NewActor(entity, p.Service)
	return
}

func (p *plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(topic string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventStateChanged:
	case events.EventCallEntityAction:
		value, ok := p.Actors.Load(v.EntityId)
		if !ok {
			return
		}
		actor := value.(*Actor)
		actor.addAction(v)
	}
}

func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginBuiltIn
}

func (p *plugin) Depends() []string {
	return nil
}

func (p *plugin) Version() string {
	return "0.0.1"
}

func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Actors:       true,
		ActorAttrs:   NewAttr(),
		ActorSetts:   NewSettings(),
		ActorActions: supervisor.ToEntityActionShort(NewActions()),
		ActorStates:  supervisor.ToEntityStateShort(NewStates()),
	}
}
