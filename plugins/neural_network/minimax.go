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

type Move struct {
	Row int
	Col int
}

// This function returns true if there are moves
// remaining on the board. It returns false if
// there are no moves left to play.
func isMovesLeft(board [3][3]GameState) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				return true
			}
		}
	}
	return false
}

// This is the evaluation function as discussed
// in the previous article
func evaluate(b [3][3]GameState, player GameState) int {
	// Checking for Rows for X or O victory.
	for row := 0; row < 3; row++ {
		if b[row][0] == b[row][1] && b[row][1] == b[row][2] {
			if b[row][0] == player {
				return 10
			} else if b[row][0] == player.opponent() {
				return -10
			}
		}
	}

	// Checking for Columns for X or O victory.
	for col := 0; col < 3; col++ {
		if b[0][col] == b[1][col] && b[1][col] == b[2][col] {
			if b[0][col] == player {
				return 10
			} else if b[0][col] == player.opponent() {
				return -10
			}
		}
	}

	// Checking for Diagonals for X or O victory.
	if b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		if b[0][0] == player {
			return 10
		} else if b[0][0] == player.opponent() {
			return -10
		}
	}

	if b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		if b[0][2] == player {
			return 10
		} else if b[0][2] == player.opponent() {
			return -10
		}
	}

	// Else, if none of them have won, then return 0
	return 0
}

// This is the minimax function. It
// considers all the possible ways
// the game can go and returns the
// value of the board
func minimax(board [3][3]GameState, depth int, isMax bool, player GameState) int {
	score := evaluate(board, player)

	// If Maximizer has won the game
	// return his/her evaluated score
	if score == 10 {
		return score
	}

	// If Minimizer has won the game
	// return his/her evaluated score
	if score == -10 {
		return score
	}

	// If there are no more moves and
	// no winner then it is a tie
	if !isMovesLeft(board) {
		return 0
	}

	// If this is the maximizer's move
	if isMax {
		best := -1000

		// Traverse all cells
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {

				// Check if cell is empty
				if board[i][j] == 0 {

					// Make the move
					board[i][j] = player

					// Call minimax recursively
					// and choose the maximum value
					best = max(best, minimax(board, depth+1, !isMax, player))

					// Undo the move
					board[i][j] = 0
				}
			}
		}
		return best
	}

	// If this is the minimizer's move
	best := 1000

	// Traverse all cells
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {

			// Check if cell is empty
			if board[i][j] == 0 {

				// Make the move
				board[i][j] = player.opponent()

				// Call minimax recursively and
				// choose the minimum value
				best = min(best, minimax(board, depth+1, !isMax, player))

				// Undo the move
				board[i][j] = 0
			}
		}
	}
	return best
}

// This will return the best possible
// move for the player
func findBestMove(board [3][3]GameState, player GameState) Move {
	bestVal := -1000
	bestMove := Move{-1, -1}

	// Traverse all cells, evaluate
	// minimax function for all empty
	// cells, and return the cell
	// with the optimal value.
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {

			// Check if cell is empty
			if board[i][j] == 0 {

				// Make the move
				board[i][j] = player

				// Compute the evaluation function
				// for this move.
				moveVal := minimax(board, 0, false, player)

				// Undo the move
				board[i][j] = 0

				// If the value of the current move
				// is more than the best value, then
				// update best
				if moveVal > bestVal {
					bestMove.Row = i
					bestMove.Col = j
					bestVal = moveVal
				}
			}
		}
	}

	//fmt.Printf("The value of the best Move is: %d\n", bestVal)

	return bestMove
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
