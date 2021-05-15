package alexa

import (
	"fmt"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/message_queue"
	"sync"
)

var _ triggers.ITrigger = (*Trigger)(nil)

const (
	TriggerName         = "alexa"
	TriggerFunctionName = "automationTriggerAlexa"
	queueSize           = 10 //todo update
)

type Trigger struct {
	eventBus     event_bus.EventBus
	msgQueue     message_queue.MessageQueue
	functionName string
	name         string
}

func NewTrigger(eventBus event_bus.EventBus) (tr triggers.ITrigger) {
	return &Trigger{
		eventBus:     eventBus,
		msgQueue:     message_queue.New(queueSize),
		functionName: TriggerFunctionName,
		name:         TriggerName,
	}
}

func (t Trigger) Name() string {
	return t.name
}

func (t Trigger) AsyncAttach(wg *sync.WaitGroup) {

	if err := t.eventBus.Subscribe(TopicPluginAlexa, t.eventHandler); err != nil {
		log.Error(err.Error())
	}

	wg.Done()
}

func (t *Trigger) eventHandler(topic string, msg interface{}) {
	switch v := msg.(type) {
	case EventAlexaAction:
		t.msgQueue.Publish(fmt.Sprintf("skill_%d", v.SkillId), v)
	}
}

func (t Trigger) Subscribe(options triggers.Subscriber) error {
	if options.Payload == nil {
		return fmt.Errorf("trigger '%s' subscribe to empty topic", t.name)
	}
	log.Infof("trigger '%s' subscribe topic '%s'", t.name, t.topic(options.Payload))
	return t.msgQueue.Subscribe(t.topic(options.Payload), options.Handler)
}

func (t Trigger) Unsubscribe(options triggers.Subscriber) error {
	if options.Payload == nil {
		return fmt.Errorf("trigger '%s' unsubscribe from empty topic", t.name)
	}
	log.Infof("trigger '%s' unsubscribe topic '%s'", t.name, t.topic(options.Payload))
	return t.msgQueue.Unsubscribe(t.topic(options.Payload), options.Handler)
}

func (t Trigger) FunctionName() string {
	return t.functionName
}

func (t Trigger) topic(n interface{}) string {
	return fmt.Sprintf("skill_%v", n)
}
