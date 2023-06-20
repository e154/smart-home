package neural_network

import (
	"fmt"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/system/bus"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/plugins"
)

var (
	log = logger.MustGetLogger("plugins.neural_network")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	actorsLock *sync.Mutex
	actors     map[common.EntityId]*Actor
	quit       chan struct{}
	pause      time.Duration
}

func New() plugins.Plugable {
	return &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[common.EntityId]*Actor),
		pause:      240,
	}
}

func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.quit = make(chan struct{})

	_ = p.EventBus.Subscribe(bus.TopicEntities, p.eventHandler)

	return nil
}

func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	_ = p.EventBus.Unsubscribe(bus.TopicEntities, p.eventHandler)

	// remove actors
	for entityId, actor := range p.actors {
		actor.destroy()
		delete(p.actors, entityId)
	}

	return nil
}

func (p *plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(topic string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventStateChanged:
	case events.EventCallAction:
		actor, ok := p.actors[v.EntityId]
		if !ok {
			return
		}
		actor.addAction(v)
	}
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entity.Id]; ok {
		p.actors[entity.Id].Update()
		return
	}

	p.actors[entity.Id] = NewActor(entity, p.EntityManager, p.Adaptors, p.ScriptService, p.EventBus)
	p.EntityManager.Spawn(p.actors[entity.Id].Spawn)

	return
}

func (p *plugin) RemoveActor(entityId common.EntityId) error {
	return p.removeEntity(entityId)
}

func (p *plugin) removeEntity(entityId common.EntityId) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entityId]; !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("failed remove \"%s\"", entityId.Name()))
		return
	}

	delete(p.actors, entityId)

	return
}

func (p *plugin) updateForAll() {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	for _, actor := range p.actors {
		actor.Update()
	}
}

func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
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
		ActorActions: entity_manager.ToEntityActionShort(NewActions()),
		ActorStates:  entity_manager.ToEntityStateShort(NewStates()),
	}
}
