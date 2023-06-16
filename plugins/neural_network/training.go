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
