package game

import (
	"math/rand"

	"github.com/wangle201210/sudoku/internal/types"
)

// GenerateValidBoard generates a valid 9x9 Sudoku board
func GenerateValidBoard(g *types.Game) {
	// Initialize empty board
	g.Size = 9
	g.Board = make([][]int, g.Size)
	g.Solution = make([][]int, g.Size)
	for i := range g.Board {
		g.Board[i] = make([]int, g.Size)
		g.Solution[i] = make([]int, g.Size)
	}

	// Generate complete solution
	generateSolution(g)

	// Copy solution to game board
	for i := range g.Board {
		copy(g.Board[i], g.Solution[i])
	}

	// Remove numbers based on difficulty
	cellsToRemove := 0
	switch g.Difficulty {
	case types.Easy:
		cellsToRemove = 40 // Leave ~41 numbers
	case types.Medium:
		cellsToRemove = 50 // Leave ~31 numbers
	case types.Hard:
		cellsToRemove = 60 // Leave ~21 numbers
	}

	// Remove numbers while ensuring unique solution
	removed := 0
	maxAttempts := g.Size * g.Size * 2
	for removed < cellsToRemove && maxAttempts > 0 {
		row := rand.Intn(g.Size)
		col := rand.Intn(g.Size)

		if g.Board[row][col] != 0 {
			temp := g.Board[row][col]
			g.Board[row][col] = 0

			// Check if puzzle still has unique solution
			if HasUniqueSolution(g) {
				removed++
			} else {
				g.Board[row][col] = temp
			}
		}
		maxAttempts--
	}
}

// generateSolution generates a complete valid Sudoku solution
func generateSolution(g *types.Game) {
	// Clear the solution grid
	for i := range g.Solution {
		for j := range g.Solution[i] {
			g.Solution[i][j] = 0
		}
	}

	// Fill diagonal 3x3 boxes first (can be filled independently)
	fillDiagonalBoxes(g)

	// Fill remaining cells
	fillRemaining(g, 0, 3)
}

// fillDiagonalBoxes fills the three diagonal 3x3 boxes with valid numbers
func fillDiagonalBoxes(g *types.Game) {
	for box := 0; box < 9; box += 3 {
		nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		rand.Shuffle(len(nums), func(i, j int) {
			nums[i], nums[j] = nums[j], nums[i]
		})

		index := 0
		for i := box; i < box+3; i++ {
			for j := box; j < box+3; j++ {
				g.Solution[i][j] = nums[index]
				index++
			}
		}
	}
}

// fillRemaining fills the remaining cells recursively
func fillRemaining(g *types.Game, row, col int) bool {
	if col >= g.Size && row < g.Size-1 {
		row++
		col = 0
	}
	if row >= g.Size && col >= g.Size {
		return true
	}
	if row < 3 {
		if col < 3 {
			col = 3
		}
	} else if row < 6 {
		if col == (row/3)*3 {
			col += 3
		}
	} else {
		if col == 6 {
			row++
			col = 0
			if row >= g.Size {
				return true
			}
		}
	}

	for num := 1; num <= g.Size; num++ {
		if IsValidMove(g, g.Solution, row, col, num) {
			g.Solution[row][col] = num
			if fillRemaining(g, row, col+1) {
				return true
			}
			g.Solution[row][col] = 0
		}
	}
	return false
}
