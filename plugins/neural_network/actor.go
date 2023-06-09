package neural_network

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/plugins/notify"
	"github.com/e154/smart-home/plugins/webpush"
	"github.com/e154/smart-home/system/scripts"
	deep "github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
	"math"

	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
)

type Actor struct {
	entity_manager.BaseActor
	eventBus      bus.Bus
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	actionPool    chan events.EventCallAction
	train2        bool
	n2            *deep.Neural
}

func NewActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	eventBus bus.Bus) *Actor {

	actor := &Actor{
		BaseActor:     entity_manager.NewBaseActor(entity, scriptService, adaptors),
		adaptors:      adaptors,
		scriptService: scriptService,
		eventBus:      eventBus,
		actionPool:    make(chan events.EventCallAction, 10),
		train2:        true,
	}

	actor.Manager = entityManager

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			// bind
			a.ScriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
			_, _ = a.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
			_, _ = a.ScriptEngine.Do()
		}
	}

	if actor.ScriptEngine != nil {
		_, _ = actor.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
		actor.ScriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
	}

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return actor
}

func (e *Actor) destroy() {

}

func (e *Actor) Spawn() entity_manager.PluginActor {
	e.Update()
	return e
}

func (e *Actor) Update() {

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
	switch action.Name {
	case "ENABLE":
		e.train2 = true
	case "DISABLE":
		e.train2 = false
	case "TRAIN1":
		e.Train1()
	case "TRAIN2":
		e.Train2()
	case "CHECK2":
		if e.n2 != nil {
			fmt.Println([]float64{1.0, 1.0, 20.0}, "=>", e.n2.Predict([]float64{1.0, 1.0, 20.0})[0]) // 2
			fmt.Println([]float64{1.0, 1.0, 1.0}, "=>", e.n2.Predict([]float64{1.0, 1.0, 1.0})[0])   // 0
			fmt.Println([]float64{0.5, 0.0, 1.0}, "=>", e.n2.Predict([]float64{0.5, 0.0, 1.0})[0])   // 1
			fmt.Println([]float64{0, 1.0, 20.0}, "=>", e.n2.Predict([]float64{0, 1.0, 20.0})[0])   // 4
		}
	default:
		fmt.Sprintf("unknown comand: %s", action.Name)
	}
}

func (e *Actor) Start() {
}

func (e *Actor) Stop() {
}

// Train2 ...
// 0 Атаковать
// 1 Красться
// 2 Убегать
// 3 Ничего не делать
// 4 Сдох
func (e *Actor) Train2() {
	var data = training.Examples{
		{[]float64{0.5, 1, 1}, []float64{0}},
		{[]float64{0.9, 1, 2}, []float64{0}},
		{[]float64{0.8, 0, 1}, []float64{0}},
		{[]float64{0.3, 1, 1}, []float64{1}},
		{[]float64{0.6, 1, 2}, []float64{1}},
		{[]float64{0.4, 0, 1}, []float64{1}},
		{[]float64{0.9, 1, 7}, []float64{2}},
		{[]float64{0.6, 1, 4}, []float64{2}},
		{[]float64{0.1, 0, 1}, []float64{2}},
		{[]float64{0.6, 1, 0}, []float64{3}},
		{[]float64{1, 0, 0}, []float64{3}},
		{[]float64{0.5, 0, 0}, []float64{3}},
		{[]float64{0.0, 0.0, 0.0}, []float64{4}},
		{[]float64{0.0, 1.0, 0.0}, []float64{4}},
		{[]float64{0.0, 0.0, 1.0}, []float64{4}},
		{[]float64{0.0, 0.0, 10.0}, []float64{4}},
		{[]float64{0.0, 1.0, 10.0}, []float64{4}},
	}

	var iterations int

LOOP:
	if iterations >= 100 {
		fmt.Println("BREAK!!!")
		return
	}
	iterations++
	fmt.Println("Iteration", iterations)

	//if e.n2 == nil {
	e.n2 = deep.NewNeural(&deep.Config{
		/* Input dimensionality */
		Inputs: 3,
		/* Two hidden layers consisting of two neurons each, and a single output */
		Layout: []int{3, 3, 1},
		/* Activation functions: Sigmoid, Tanh, ReLU, Linear */
		Activation: deep.ActivationSigmoid,
		/* Determines output layer activation & loss function:
		ModeRegression: linear outputs with MSE loss
		ModeMultiClass: softmax output with Cross Entropy loss
		ModeMultiLabel: sigmoid output with Cross Entropy loss
		ModeBinary: sigmoid output with binary CE loss */
		Mode: deep.ModeRegression,
		/* Weight initializers: {deep.NewNormal(μ, σ), deep.NewUniform(μ, σ)} */
		Weight: deep.NewNormal(1, 0),
		/* Apply bias */
		Bias: true,
	})
	//}

	if e.train2 {
		// params: learning rate, momentum, alpha decay, nesterov
		optimizer := training.NewSGD(0.05, 0.1, 1e-6, true)
		// params: optimizer, verbosity (print stats at every 50th iteration)
		trainer := training.NewTrainer(optimizer, 100)

		training, heldout := data.Split(0.9)
		trainer.Train(e.n2, training, heldout, 10000) // training, validation, iterations
	}

	fmt.Println("Train result:")

	modRound := func(x float64) float64 {
		if x < 0 {
			x *= -1.0
		}
		return math.Round(x)
	}

	for _, i := range data {
		num := e.n2.Predict(i.Input)[0]
		fmt.Println(i.Input, "=>", num, modRound(num))
	}

	fmt.Println("new predict:")

	t := training.Examples{
		{[]float64{0, 1, 18}, []float64{4}},      // 4
		{[]float64{0.0, 0.0, 1.0}, []float64{4}}, // 4
		{[]float64{0.2, 0, 0}, []float64{3}},     // 3
		{[]float64{0.2, 1, 0}, []float64{3}},     // 3
		{[]float64{0.9, 1, 8}, []float64{2}},     // 2
		{[]float64{0.9, 1, 18}, []float64{2}},    // 2
	}

	for _, i := range t {
		num := e.n2.Predict(i.Input)[0]
		fmt.Println(i.Input, "=>", num, modRound(num))
		if modRound(num) != i.Response[0] {
			goto LOOP
		}
	}

	e.eventBus.Publish(notify.TopicNotify, notify.Message{
		Type: webpush.Name,
		Attributes: map[string]interface{}{
			webpush.AttrUserIDS: "14",
			webpush.AttrTitle:   "neural network",
			webpush.AttrBody:    "all completed",
		},
	})
}

func (e *Actor) Train1() {
	var data = training.Examples{
		{[]float64{2.7810836, 2.550537003}, []float64{0}},
		{[]float64{1.465489372, 2.362125076}, []float64{0}},
		{[]float64{3.396561688, 4.400293529}, []float64{0}},
		{[]float64{1.38807019, 1.850220317}, []float64{0}},
		{[]float64{7.627531214, 2.759262235}, []float64{1}},
		{[]float64{5.332441248, 2.088626775}, []float64{1}},
		{[]float64{6.922596716, 1.77106367}, []float64{1}},
		{[]float64{8.675418651, -0.242068655}, []float64{1}},
	}

	n := deep.NewNeural(&deep.Config{
		/* Input dimensionality */
		Inputs: 2,
		/* Two hidden layers consisting of two neurons each, and a single output */
		Layout: []int{2, 2, 1},
		/* Activation functions: Sigmoid, Tanh, ReLU, Linear */
		Activation: deep.ActivationSigmoid,
		/* Determines output layer activation & loss function:
		ModeRegression: linear outputs with MSE loss
		ModeMultiClass: softmax output with Cross Entropy loss
		ModeMultiLabel: sigmoid output with Cross Entropy loss
		ModeBinary: sigmoid output with binary CE loss */
		Mode: deep.ModeBinary,
		/* Weight initializers: {deep.NewNormal(μ, σ), deep.NewUniform(μ, σ)} */
		Weight: deep.NewNormal(1.0, 0.0),
		/* Apply bias */
		Bias: true,
	})

	// params: learning rate, momentum, alpha decay, nesterov
	optimizer := training.NewSGD(0.05, 0.1, 1e-6, true)
	// params: optimizer, verbosity (print stats at every 50th iteration)
	trainer := training.NewTrainer(optimizer, 50)

	training, heldout := data.Split(0.5)
	trainer.Train(n, training, heldout, 1000) // training, validation, iterations

	fmt.Println(data[0].Input, "=>", n.Predict(data[0].Input))
	fmt.Println(data[5].Input, "=>", n.Predict(data[5].Input))
}
