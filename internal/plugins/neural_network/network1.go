// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package neural_network

import (
	"fmt"
	"math"

	"github.com/e154/smart-home/internal/plugins/notify"
	"github.com/e154/smart-home/internal/plugins/notify/common"
	"github.com/e154/smart-home/internal/plugins/webpush"

	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"

	"github.com/e154/bus"
)

type Network1 struct {
	train2   bool
	n2       *deep.Neural
	eventBus bus.Bus
}

func NewNetwork1(eventBus bus.Bus) (net *Network1) {
	net = &Network1{
		eventBus: eventBus,
		train2:   true,
	}
	return net
}

// Train222 ...
func (e *Network1) Train222() {
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

	e.eventBus.Publish(notify.TopicNotify, common.Message{
		Type: webpush.Name,
		Attributes: map[string]interface{}{
			webpush.AttrUserIDS: "14",
			webpush.AttrTitle:   "neural network",
			webpush.AttrBody:    "all completed",
		},
	})
}

func (e *Network1) Train1() {
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

	for _, i := range data {
		num := n.Predict(i.Input)[0]
		fmt.Println(i.Input, "=>", num, modRound(num))
	}
}

func (e *Network1) Train2() {
	var data = training.Examples{
		{[]float64{0, 0}, []float64{0}},
		{[]float64{0, 1.0}, []float64{1.0}},
		{[]float64{1.0, 0}, []float64{1.0}},
		{[]float64{1.0, 1.0}, []float64{0}},
	}

	n := deep.NewNeural(&deep.Config{
		/* Input dimensionality */
		Inputs: 2,
		/* Two hidden layers consisting of two neurons each, and a single output */
		Layout: []int{20, 20, 1},
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

	training, heldout := data.Split(1)
	trainer.Train(n, training, heldout, 1000) // training, validation, iterations

	for _, i := range data {
		num := n.Predict(i.Input)[0]
		fmt.Println(i.Input, "=>", num, modRound(num))
	}
}

var modRound = func(x float64) float64 {
	if x < 0 {
		x *= -1.0
	}
	return math.Round(x)
}
