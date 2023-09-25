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

type TrainingData struct {
	Input  []float64
	Output []int
}

func GenerateTrainingData() []TrainingData {
	var trainingData []TrainingData

	// Генерируем случайные игровые ситуации и сохраняем оптимальные ходы компьютера
	for i := 0; i < 10; i++ {
		game := NewGame()

		for !game.isGameOver() {
			row, col := game.getRandomMove()
			trainingData = append(trainingData, TrainingData{
				Input:  game.getBoardState(),
				Output: []int{row, col},
			})

			game.MakeMove(row, col)
		}
	}

	return trainingData
}
