package mqtt

import (
	"sync"

	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	message          *Message
	mqttMessageQueue chan *Message
	actionPool       chan events.EventCallEntityAction
	mqttClient       mqtt.MqttCli
	newMsgMu         *sync.Mutex
	stateMu          *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service,
	mqttClient mqtt.MqttCli) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor:        supervisor.NewBaseActor(entity, service),
		message:          NewMessage(),
		mqttMessageQueue: make(chan *Message, 10),
		actionPool:       make(chan events.EventCallEntityAction, 10),
		mqttClient:       mqttClient,
		newMsgMu:         &sync.Mutex{},
		stateMu:          &sync.Mutex{},
	}

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			a.ScriptEngine.PushStruct("Actor", supervisor.NewScriptBind(actor))
			_, _ = a.ScriptEngine.Do()
		}
	}

	if actor.ScriptEngine != nil {
		actor.ScriptEngine.PushStruct("message", actor.message)
		actor.ScriptEngine.PushStruct("Actor", supervisor.NewScriptBind(actor))
	}

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

func (e *Actor) Destroy() {
	if e.Setts != nil && e.Setts[AttrSubscribeTopic] != nil {
		e.mqttClient.Unsubscribe(e.Setts[AttrSubscribeTopic].String())
	}
}

// Spawn ...
func (e *Actor) Spawn() {

	if e.Setts != nil && e.Setts[AttrSubscribeTopic] != nil {
		_ = e.mqttClient.Subscribe(e.Setts[AttrSubscribeTopic].String(), e.mqttOnPublish)
	}

	return
}

// SetState ...
func (e *Actor) SetState(params supervisor.EntityStateParams) error {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()

	oldState := e.GetEventState()
	now := e.Now(oldState)

	if params.NewState != nil {
		if state, ok := e.States[*params.NewState]; ok {
			e.State = &state
		}
	}

	e.AttrMu.Lock()
	changed, err := e.Attrs.Deserialize(params.AttributeValues)
	if !changed {
		if err != nil {
			log.Warn(err.Error())
		}

		if oldState.LastUpdated != nil {
			delta := now.Sub(*oldState.LastUpdated).Milliseconds()
			//fmt.Println("delta", delta)
			if delta < 200 {
				e.AttrMu.Unlock()
				return nil
			}
		}
	}
	e.AttrMu.Unlock()

	e.Service.EventBus().Publish("system/entities/"+e.Id.String(), events.EventStateChanged{
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(),
		StorageSave: params.StorageSave,
	})

	return nil
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

func (e *Actor) addAction(event events.EventCallEntityAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallEntityAction) {
	action, ok := e.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if action.ScriptEngine == nil {
		return
	}
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId, action.Name, msg.Args); err != nil {
		log.Error(err.Error())
	}
}
