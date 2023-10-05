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
)

type GameState int

const (
	Empty GameState = 0
	X     GameState = 1
	O     GameState = -1
)

type Game struct {
	Board  [3][3]GameState
	Player GameState
}

func NewGame() *Game {
	game := &Game{}
	game.initBoard()
	game.Player = X
	return game
}

func (g *Game) initBoard() {
	g.Board = [3][3]GameState{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
}

func (g *Game) UpdateBoard(board [3][3]GameState) {
	g.Board = board
}

func (g *Game) NumToCell(num int) (row, col int) {
	row = num / 3
	col = num - row*3
	return
}

func (g *Game) MakeMove(row, col int) error {
	if row < 0 || row >= 3 || col < 0 || col >= 3 || g.Board[row][col] != Empty {
		return fmt.Errorf("Недопустимый ход")
	}

	g.Board[row][col] = g.Player
	g.Player = g.Player.opponent()
	return nil
}

func (g *Game) isGameOver() bool {
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == Empty {
				return false
			}
		}
	}
	return true
}

func (g *Game) checkWinner() GameState {
	// Проверка горизонтальных и вертикальных линий
	for i := 0; i < 3; i++ {
		if g.Board[i][0] != Empty && g.Board[i][0] == g.Board[i][1] && g.Board[i][0] == g.Board[i][2] {
			return g.Board[i][0]
		}
		if g.Board[0][i] != Empty && g.Board[0][i] == g.Board[1][i] && g.Board[0][i] == g.Board[2][i] {
			return g.Board[0][i]
		}
	}

	// Проверка диагональных линий
	if g.Board[0][0] != Empty && g.Board[0][0] == g.Board[1][1] && g.Board[0][0] == g.Board[2][2] {
		return g.Board[0][0]
	}
	if g.Board[0][2] != Empty && g.Board[0][2] == g.Board[1][1] && g.Board[0][2] == g.Board[2][0] {
		return g.Board[0][2]
	}

	return Empty
}

func (g *Game) PrintBoard() {
	for _, row := range g.Board {
		for _, cell := range row {
			switch cell {
			case X:
				fmt.Print("X ")
			case O:
				fmt.Print("O ")
			case Empty:
				fmt.Print("- ")
			}
		}
		fmt.Println()
	}
}

func (g *Game) getRandomMove() (int, int) {
	bestMove := findBestMove(g.Board, g.Player)
	return bestMove.Row, bestMove.Col
}

func (g *Game) getBoardState() []float64 {
	var boardState []float64
	for _, row := range g.Board {
		for _, cell := range row {
			switch cell {
			case X:
				boardState = append(boardState, 1.0)
			case O:
				boardState = append(boardState, -1.0)
			case Empty:
				boardState = append(boardState, 0.0)
			}
		}
	}
	return boardState
}

func (g *Game) getMoveOutput(row, col int) []float64 {
	return []float64{float64(row), float64(col)}
}

func (state GameState) opponent() GameState {
	if state == X {
		return O
	}
	return X
}
