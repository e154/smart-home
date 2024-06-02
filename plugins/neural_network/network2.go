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
	"math/rand"
	"time"

	"github.com/e154/bus"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	"github.com/julioguillermo/staticneurogenetic"
)

const fileName = "oxo.bin"

type Network2 struct {
	train2     bool
	eventBus   bus.Bus
	game       *Game
	entityId   common.EntityId
	agents     *staticneurogenetic.SNG
	moves      int
	individual int
}

func NewNetwork2(eventBus bus.Bus) (net *Network2) {
	net = &Network2{
		eventBus: eventBus,
		train2:   true,
		game:     NewGame(),
		entityId: "sensor.ticTacToe",
	}
	net.Start()
	return net
}

const (
	PopulationSize  = 100
	PopSizeChRandom = false
)

func (e *Network2) Start() {
	_ = e.eventBus.Subscribe("system/entities/"+e.entityId.Name(), e.eventHandler)

	rand.Seed(int64(time.Now().Nanosecond()))
	var err error
	e.agents, err = staticneurogenetic.LoadFromBin(fileName)
	if err != nil {
		e.agents = staticneurogenetic.NewSNG(
			[]int{9, 50, 30, 9},
			staticneurogenetic.Relu,          //Activation function for the neural network
			100,                              //PopulationSize (number of individual to work with)
			10,                               //Survivors (number of individual that will not change in next generation and to use as parents)
			0.1,                              //MutRate (probability to mutate a new individual)
			0.1,                              //MutSize (how big the mutation will be)
			staticneurogenetic.MultiMutation, //MutType
			staticneurogenetic.DivPointCross, //CrossType
		)
		_ = e.agents.SetPopulationSize(PopulationSize, PopSizeChRandom)
	}

	//e.agents.ResetFitness()
	e.SelectPopulation()

}

func (e *Network2) Stop() {
	_ = e.eventBus.Unsubscribe("system/entities/"+e.entityId.String(), e.eventHandler)
	if err := e.agents.SaveAsBin(fileName); err != nil {
		log.Error(err.Error())
	}
	e.agents = nil
}

func (e *Network2) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventCallEntityAction:
	case events.EventStateChanged:
		if v.EntityId != e.entityId {
			return
		}
		board := [3][3]GameState{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		}
		//debug.Println(msg)
		var exit = false
		for key, attr := range v.NewState.Attributes {
			switch key {
			case "PLAYER":
				if attr.String() != "X" {
					exit = true
				}
			case "R0C0":
				board[0][0] = getState(attr.String())
			case "R0C1":
				board[0][1] = getState(attr.String())
			case "R0C2":
				board[0][2] = getState(attr.String())
			case "R1C0":
				board[1][0] = getState(attr.String())
			case "R1C1":
				board[1][1] = getState(attr.String())
			case "R1C2":
				board[1][2] = getState(attr.String())
			case "R2C0":
				board[2][0] = getState(attr.String())
			case "R2C1":
				board[2][1] = getState(attr.String())
			case "R2C2":
				board[2][2] = getState(attr.String())
			case "WINNER":
				if attr.String() != "" {
					exit = true
				}
				if attr.String() == "O" {
					fmt.Println("ARGHHH!!!")
					e.agents.AddFitness(e.individual, 0.1)
				} else {
					e.agents.AddFitness(e.individual, -0.1)
				}
			}
		}
		if exit {
			return
		}
		e.game.UpdateBoard(board)
		e.MakeMove()
	}
}

func (e *Network2) sendMoveCommand(row, col int) {

}

func (e *Network2) MakeMove() {

	fmt.Println("MakeMove")

	e.moves++

	board := e.game.getBoardState()
	_, pos := e.agents.MaxOutput(0, board)
	if err := e.game.MakeMove(e.game.NumToCell(pos)); err != nil {

		fmt.Println(err.Error(), pos, board)
		if e.game.isGameOver() {
			return
		}
		e.agents.AddFitness(0, -1)
		row, col := e.game.getRandomMove()
		e.sendMoveCommand(row, col)
		return
	}
	row, col := e.game.NumToCell(pos)
	e.sendMoveCommand(row, col)
}

func (e *Network2) Train2() {

	// Генерация обучающих данных
	//trainingData := GenerateTrainingData()
	//debug.Println(len(trainingData))
	//debug.Println(trainingData)

	fmt.Println("Start")

	e.agents.ResetFitness()

	for k := 0; k < 50; k++ {

		generation := e.agents.GetGeneration()
		fmt.Println("generation:", generation)

		for index := range e.agents.Population {

			for i := 0; i < 100; i++ { //for2

				game := NewGame()
				moves := 0.0

				// X
				game.MakeMove(game.getRandomMove())

				for !game.isGameOver() { //for1
					moves++
					//game.PrintBoard()

					_, pos := e.agents.MaxOutput(index, game.getBoardState())
					// O
					if err := game.MakeMove(game.NumToCell(pos)); err != nil {
						//fmt.Println("error1", err.Error())
						e.agents.AddFitness(index, -1/moves)
						// O
						game.MakeMove(game.getRandomMove())
						// X
						game.MakeMove(game.getRandomMove())
						if game.isGameOver() {
							if winner := game.checkWinner(); winner == X {
								e.agents.AddFitness(index, -1/moves)
							}
						}
					} else {

						if game.isGameOver() {
							if winner := game.checkWinner(); winner == X {
								e.agents.AddFitness(index, -1/moves)
							} else if winner == O {
								e.agents.AddFitness(index, 3/moves)
							}
						} else {
							//e.agents.AddFitness(index, 3/moves)
							// X
							game.MakeMove(game.getRandomMove())
							if game.isGameOver() {
								e.agents.AddFitness(index, -1/moves)
							}
						}

					}

				} // \for1

			} // \for2
		}

		e.agents.NextGeneration()

	} // \300

	if err := e.agents.SaveAsBin(fileName); err != nil {
		log.Error(err.Error())
	}

	e.SelectPopulation()

	fmt.Println("END")
}

func (e *Network2) SelectPopulation() {

	e.agents.Sort()

	e.individual = e.agents.GetLastBestIndex()

	fmt.Println("total populations:")
	for index, population := range e.agents.Population {
		fmt.Printf("population: %d, %f\n", index, population.Fitness)
	}

	fmt.Printf("selected population: %f, num: %d\n", e.agents.Population[e.individual].Fitness, e.individual)
}

var (
	inputs = [][]float64{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	}
	targets = []float64{
		1,
		0,
		0,
		1,
	}
)

func (e *Network2) Train1() {

	agents := staticneurogenetic.NewSNG(
		[]int{2, 3, 1},                   //Neural network's layers [input, hiddens..., output]
		staticneurogenetic.Sigmoid,       //Activation function for the neural network
		300,                              //PopulationSize (number of individual to work with)
		10,                               //Survivors (number of individual that will not change in next generation and to use as parents)
		0.1,                              //MutRate (probability to mutate a new individual)
		0.1,                              //MutSize (how big the mutation will be)
		staticneurogenetic.OneMutation,   //MutType
		staticneurogenetic.DivPointCross, //CrossType
	)

	for i := 0; i < 300; i++ {
		agents.ResetFitness() //Set all fitness to 0, for use AddFitness
		evalAll(agents)
		agents.NextGeneration() //Evolve each neural networks
	}

	//fmt.Println("-----")
	//for i := 0; i < 5; i++ {
	//	value := agents.Output(i, []float64{0, 0})
	//	fmt.Println(value)
	//	val, indexMax := agents.MaxOutput(i, []float64{0, 0})
	//	fmt.Println(val, indexMax)
	//	val, indexMin := agents.MinOutput(i, []float64{0, 0})
	//	fmt.Println(val, indexMin)
	//	fmt.Println("-----")
	//}

}

// Eval an individual
func eval(agents *staticneurogenetic.SNG, individual int) {

	for i, input := range inputs {
		// Get individual output ([]float64)
		output := agents.Output(individual, input)
		// Calculate how wrong is the output
		dif := math.Abs(targets[i] - output[0])
		// Added to the fitness
		fmt.Println(1 - dif)
		agents.AddFitness(individual, 1-dif)
	}
}

// Eval each individual
func evalAll(agents *staticneurogenetic.SNG) {
	for i := range agents.Population {
		eval(agents, i)
	}
}

func getState(val string) (result GameState) {
	switch val {
	case "X":
		result = 1
	case "O":
		result = -1
	}
	return
}
