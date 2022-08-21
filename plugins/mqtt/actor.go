package mqtt

import (
	"fmt"
	"sync"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	eventBus         bus.Bus
	adaptors         *adaptors.Adaptors
	scriptService    scripts.ScriptService
	message          *Message
	mqttMessageQueue chan *Message
	actionPool       chan events.EventCallAction
	mqttClient       mqtt.MqttCli
	newMsgMu         *sync.Mutex
	stateMu          *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	params map[string]interface{},
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	entityManager entity_manager.EntityManager,
	eventBus bus.Bus,
	mqttClient mqtt.MqttCli) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor:        entity_manager.NewBaseActor(entity, scriptService, adaptors),
		eventBus:         eventBus,
		adaptors:         adaptors,
		scriptService:    scriptService,
		message:          NewMessage(),
		mqttMessageQueue: make(chan *Message, 10),
		actionPool:       make(chan events.EventCallAction, 10),
		mqttClient:       mqttClient,
		newMsgMu:         &sync.Mutex{},
		stateMu:          &sync.Mutex{},
	}

	actor.Manager = entityManager
	_, _ = actor.Attrs.Deserialize(params)

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			_, _ = a.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
			a.ScriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
			_, _ = a.ScriptEngine.Do()
		}
	}

	if actor.ScriptEngine != nil {
		// message
		actor.ScriptEngine.PushStruct("message", actor.message)

		// binds
		_, _ = actor.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
		actor.ScriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
	}

	actor.Manager = entityManager
	actor.Setts = entity.Settings

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	// mqtt worker
	go func() {
		for message := range actor.mqttMessageQueue {
			actor.mqttNewMessage(message)
		}
	}()

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return
}

func (e *Actor) destroy() {
	if e.Setts != nil && e.Setts[AttrSubscribeTopic] != nil {
		e.mqttClient.Unsubscribe(e.Setts[AttrSubscribeTopic].String())
	}
}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {

	if e.Setts != nil && e.Setts[AttrSubscribeTopic] != nil {
		_ = e.mqttClient.Subscribe(e.Setts[AttrSubscribeTopic].String(), e.mqttOnPublish)
	}

	return e
}

// SetState ...
func (e *Actor) SetState(params entity_manager.EntityStateParams) error {
	if !e.setState(params) {
		return nil
	}

	message := NewMessage()
	message.NewState = params

	e.mqttMessageQueue <- message

	return nil
}

func (e *Actor) setState(params entity_manager.EntityStateParams) (changed bool) {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()

	oldState := e.GetEventState(e)
	now := e.Now(oldState)

	if params.NewState != nil {
		if state, ok := e.States[*params.NewState]; ok {
			e.State = &state
		}
	}

	e.AttrMu.Lock()
	var err error
	if changed, err = e.Attrs.Deserialize(params.AttributeValues); !changed {
		if err != nil {
			log.Warn(err.Error())
		}

		if oldState.LastUpdated != nil {
			delta := now.Sub(*oldState.LastUpdated).Milliseconds()
			//fmt.Println("delta", delta)
			if delta < 200 {
				e.AttrMu.Unlock()
				return
			}
		}
	}
	e.AttrMu.Unlock()

	e.eventBus.Publish(bus.TopicEntities, events.EventStateChanged{
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
		StorageSave: params.StorageSave,
	})

	return
}

func (e *Actor) mqttOnPublish(client mqtt.MqttCli, msg mqtt.Message) {
	message := NewMessage()
	message.Payload = string(msg.Payload)
	message.Topic = msg.Topic
	message.Qos = msg.Qos
	message.Duplicate = msg.Dup

	e.mqttMessageQueue <- message
}

func (e *Actor) mqttNewMessage(message *Message) {

	e.newMsgMu.Lock()
	defer e.newMsgMu.Unlock()

	e.message.Update(message)
	if e.ScriptEngine == nil {
		return
	}
	if _, err := e.ScriptEngine.AssertFunction(FuncMqttEvent); err != nil {
		log.Error(err.Error())
		return
	}
}

func (e *Actor) addAction(event events.EventCallAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallAction) {
	action, ok := e.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if action.ScriptEngine == nil {
		return
	}
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId, action.Name); err != nil {
		log.Error(err.Error())
	}
}
