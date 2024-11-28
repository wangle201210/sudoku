package game

import (
	"github.com/wangle201210/sudoku/internal/types"
)

// IsValidMove checks if placing a number at the specified position is valid
func IsValidMove(g *types.Game, board [][]int, row, col, num int) bool {
	// Check row
	for x := 0; x < g.Size; x++ {
		if board[row][x] == num {
			return false
		}
	}

	// Check column
	for x := 0; x < g.Size; x++ {
		if board[x][col] == num {
			return false
		}
	}

	// Check 3x3 box
	startRow := row - (row % 3)
	startCol := col - (col % 3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[startRow+i][startCol+j] == num {
				return false
			}
		}
	}

	return true
}

// HasUniqueSolution checks if the Sudoku puzzle has exactly one solution
func HasUniqueSolution(g *types.Game) bool {
	// Create a copy of the board for solving
	tempBoard := make([][]int, g.Size)
	for i := range tempBoard {
		tempBoard[i] = make([]int, g.Size)
		copy(tempBoard[i], g.Board[i])
	}

	// Try to find first solution
	if !solveSudoku(g, tempBoard, 0, 0) {
		return false // No solution exists
	}

	// Store first solution
	firstSolution := make([][]int, g.Size)
	for i := range firstSolution {
		firstSolution[i] = make([]int, g.Size)
		copy(firstSolution[i], tempBoard[i])
	}

	// Reset board and try to find second solution
	for i := range tempBoard {
		copy(tempBoard[i], g.Board[i])
	}

	// Start from the next possible move
	return !findAnotherSolution(g, tempBoard, 0, 0, firstSolution)
}

// findAnotherSolution attempts to find a different solution than the first one
func findAnotherSolution(g *types.Game, board [][]int, row, col int, firstSolution [][]int) bool {
	if row == g.Size {
		// Check if this solution is different from the first one
		for i := 0; i < g.Size; i++ {
			for j := 0; j < g.Size; j++ {
				if board[i][j] != firstSolution[i][j] {
					return true
				}
			}
		}
		return false
	}

	nextRow, nextCol := getNextCell(row, col, g.Size)

	if board[row][col] != 0 {
		return findAnotherSolution(g, board, nextRow, nextCol, firstSolution)
	}

	for num := 1; num <= g.Size; num++ {
		if IsValidMove(g, board, row, col, num) {
			board[row][col] = num
			if findAnotherSolution(g, board, nextRow, nextCol, firstSolution) {
				return true
			}
			board[row][col] = 0
		}
	}

	return false
}

// solveSudoku solves the Sudoku puzzle using backtracking
func solveSudoku(g *types.Game, board [][]int, row, col int) bool {
	if row == g.Size {
		return true
	}

	nextRow, nextCol := getNextCell(row, col, g.Size)

	if board[row][col] != 0 {
		return solveSudoku(g, board, nextRow, nextCol)
	}

	for num := 1; num <= g.Size; num++ {
		if IsValidMove(g, board, row, col, num) {
			board[row][col] = num
			if solveSudoku(g, board, nextRow, nextCol) {
				return true
			}
			board[row][col] = 0
		}
	}

	return false
}

// getNextCell returns the next cell position to check
func getNextCell(row, col, size int) (int, int) {
	col++
	if col == size {
		col = 0
		row++
	}
	return row, col
}

// CheckWin verifies if the current board state represents a winning condition
func CheckWin(g *types.Game) bool {
	// Check for empty cells
	for i := 0; i < g.Size; i++ {
		for j := 0; j < g.Size; j++ {
			if g.Board[i][j] == 0 {
				return false
			}
		}
	}

	// Check each row
	for row := 0; row < g.Size; row++ {
		nums := make(map[int]bool)
		for col := 0; col < g.Size; col++ {
			num := g.Board[row][col]
			if num < 1 || num > g.Size || nums[num] {
				return false
			}
			nums[num] = true
		}
	}

	// Check each column
	for col := 0; col < g.Size; col++ {
		nums := make(map[int]bool)
		for row := 0; row < g.Size; row++ {
			num := g.Board[row][col]
			if num < 1 || num > g.Size || nums[num] {
				return false
			}
			nums[num] = true
		}
	}

	// Check each 3x3 box
	for boxRow := 0; boxRow < g.Size; boxRow += 3 {
		for boxCol := 0; boxCol < g.Size; boxCol += 3 {
			nums := make(map[int]bool)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					num := g.Board[boxRow+i][boxCol+j]
					if num < 1 || num > g.Size || nums[num] {
						return false
					}
					nums[num] = true
				}
			}
		}
	}

	return true
}
