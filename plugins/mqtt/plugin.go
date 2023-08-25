package mqtt

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.mqtt")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	actorsLock *sync.Mutex
	actors     map[string]*Actor
	mqttServ   mqtt.MqttServ
	mqttClient mqtt.MqttCli
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin:     supervisor.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]*Actor),
	}
}

// Load ...
func (p *plugin) Load(service supervisor.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.mqttServ = service.MqttServ()

	p.mqttClient = p.mqttServ.NewClient("plugins.mqtt")
	if err := p.EventBus.Subscribe(bus.TopicEntities, p.eventHandler); err != nil {
		log.Error(err.Error())
	}

	_ = p.mqttServ.Authenticator().Register(p.Authenticator)

	return nil
}

// Unload ...
func (p plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.mqttServ.RemoveClient("plugins.mqtt")
	_ = p.EventBus.Unsubscribe(bus.TopicEntities, p.eventHandler)

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	// remove actors
	for entityId, actor := range p.actors {
		actor.destroy()
		delete(p.actors, entityId)
	}

	_ = p.mqttServ.Authenticator().Unregister(p.Authenticator)

	return
}

// Name ...
func (p plugin) Name() string {
	return Name
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) error {
	return p.addOrUpdateEntity(entity, entity.Attributes.Serialize())
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {
	return p.removeEntity(entityId.Name())
}

func (p *plugin) addOrUpdateEntity(entity *m.Entity, attributes m.AttributeValue) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	name := entity.Id.Name()
	if _, ok := p.actors[name]; ok {
		return
	}

	if actor, ok := p.actors[name]; ok {
		// update
		_ = actor.SetState(supervisor.EntityStateParams{
			AttributeValues: attributes,
		})
		return
	}

	var actor *Actor
	if actor, err = NewActor(entity, attributes,
		p.Adaptors, p.ScriptService, p.Supervisor, p.EventBus, p.mqttClient); err != nil {
		return
	}
	p.actors[name] = actor
	p.Supervisor.Spawn(p.actors[name].Spawn)

	return
}

func (p *plugin) removeEntity(name string) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[name]; !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("failed remove '%s", name))
		return
	}

	p.actors[name].destroy()

	delete(p.actors, name)

	return
}

func (p *plugin) topic(bridgeId string) string {
	return fmt.Sprintf("%s/#", bridgeId)
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventCallAction:
		actor, ok := p.actors[v.EntityId.Name()]
		if !ok {
			return
		}
		actor.addAction(v)

	default:
		//fmt.Printf("new event: %v\n", reflect.TypeOf(v).String())
	}
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginInstallable
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Version ...
func (p *plugin) Version() string {
	return "0.0.1"
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Triggers:           false,
		Actors:             true,
		ActorCustomAttrs:   true,
		ActorCustomActions: true,
		ActorCustomStates:  true,
		ActorCustomSetts:   true,
		ActorSetts:         NewSettings(),
		Setts:              nil,
	}
}

// Authenticator ...
func (p *plugin) Authenticator(login, password string) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	for _, actor := range p.actors {
		attrs := actor.Settings()

		if _login, ok := attrs[AttrMqttLogin]; !ok || _login.String() != login {
			continue
		}

		if _password, ok := attrs[AttrMqttPass]; !ok || _password.String() != password {
			continue
		}

		err = nil
		return

		// todo add encripted password
		//if ok := common.CheckPasswordHash(password, settings[AttrNodePass].String()); ok {
		//	return
		//}
	}

	err = apperr.ErrBadLoginOrPassword

	return
}
