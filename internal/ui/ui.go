package ui

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/wangle201210/sudoku/internal/game"
	"github.com/wangle201210/sudoku/internal/types"
)

const (
	cellSize    = 40 // cell size in pixels
	borderWidth = 2  // border width in pixels
)

var (
	borderColor  = color.NRGBA{R: 0, G: 0, B: 0, A: 255}       // black for borders
	blockBgColor = color.NRGBA{R: 240, G: 240, B: 240, A: 255} // light gray for background
)

// CreateUI creates the game UI
func CreateUI(g *types.Game) fyne.CanvasObject {
	// Initialize game if needed
	if g.Board == nil {
		g.Size = 9
		g.Difficulty = types.Easy
		game.GenerateValidBoard(g)
		g.Entries = make([][]*types.SquareEntry, g.Size)
		for i := range g.Entries {
			g.Entries[i] = make([]*types.SquareEntry, g.Size)
		}
	}

	// Create main container with padding
	mainContainer := container.NewVBox()

	// Create timer container
	timerContainer := container.NewHBox()
	g.TimeLabel = widget.NewLabel("Time: 0s")
	timerContainer.Add(g.TimeLabel)

	// Create game board
	gameBoard := createGameBoard(g)

	// Add all elements to main container with padding
	mainContainer.Add(container.NewPadded(timerContainer))
	mainContainer.Add(container.NewPadded(gameBoard))

	// Start timer
	StartTimer(g)

	return mainContainer
}

// refreshBoard updates the game board display
func refreshBoard(g *types.Game) {
	for i := 0; i < g.Size; i++ {
		for j := 0; j < g.Size; j++ {
			entry := g.Entries[i][j]
			if g.Board[i][j] != 0 {
				entry.SetText(fmt.Sprintf("%d", g.Board[i][j]))
				entry.Disable()
			} else {
				entry.SetText("")
				entry.Enable()
			}
		}
	}
}

// createGameBoard creates the game board UI
func createGameBoard(g *types.Game) fyne.CanvasObject {
	// Create game board container
	gameBoard := container.NewGridWithColumns(3)

	// Create 9 3x3 block containers
	for blockRow := 0; blockRow < 3; blockRow++ {
		for blockCol := 0; blockCol < 3; blockCol++ {
			// Create 3x3 block container
			blockContainer := container.NewGridWithColumns(3)

			// Create background and border
			background := canvas.NewRectangle(blockBgColor)
			border := canvas.NewRectangle(borderColor)
			border.StrokeWidth = float32(borderWidth)
			border.StrokeColor = borderColor
			border.FillColor = color.Transparent

			// Combine background, border and content
			blockWithBorder := container.NewMax(background, border, blockContainer)

			// Fill 3x3 block with cells
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					row := blockRow*3 + i
					col := blockCol*3 + j

					// Create entry with fixed size
					entry := types.NewSquareEntry()

					// Set input validation
					entry.Validator = func(s string) error {
						if s == "" {
							return nil
						}
						if len(s) > 1 {
							return fmt.Errorf("single digit only")
						}
						num, err := strconv.Atoi(s)
						if err != nil || num < 1 || num > 9 {
							return fmt.Errorf("enter 1-9")
						}
						return nil
					}

					// Set input handler
					entry.OnChanged = createOnChangedHandler(g, &entry.Entry, row, col)

					// Set original numbers
					if g.Board[row][col] != 0 {
						entry.SetText(fmt.Sprintf("%d", g.Board[row][col]))
						entry.Disable()
					}

					g.Entries[row][col] = entry
					blockContainer.Add(entry)
				}
			}
			gameBoard.Add(blockWithBorder)
		}
	}

	return gameBoard
}

// createOnChangedHandler creates input change handler
func createOnChangedHandler(g *types.Game, entry *widget.Entry, row, col int) func(string) {
	return func(s string) {
		// Clear the cell if input is empty
		if s == "" {
			g.Board[row][col] = 0
			return
		}

		// Validate input
		num, err := strconv.Atoi(s)
		if err != nil || num < 1 || num > g.Size {
			entry.SetText("")
			return
		}

		// Store the current value
		currentValue := g.Board[row][col]

		// Temporarily set the cell to 0 to check if the new number is valid
		g.Board[row][col] = 0
		if !game.IsValidMove(g, g.Board, row, col, num) {
			// Restore the previous value and clear the input
			g.Board[row][col] = currentValue
			entry.SetText("")
			dialog.ShowInformation("Info", "Invalid move", g.Window)
			return
		}

		// Set the new value
		g.Board[row][col] = num

		// Check for win condition
		if game.CheckWin(g) {
			dialog.ShowInformation("Congratulations", "You won!", g.Window)
			StopTimer(g)
		}
	}
}

// UpdateTimer updates the timer display
func UpdateTimer(g *types.Game) {
	if !g.Running {
		return
	}

	elapsed := time.Since(g.StartTime)
	g.ElapsedTime = int(elapsed.Seconds())
	g.TimeLabel.SetText(fmt.Sprintf("Time: %ds", g.ElapsedTime))
}

// StartTimer starts the game timer
func StartTimer(g *types.Game) {
	// Stop existing timer if any
	StopTimer(g)

	// Reset elapsed time
	g.ElapsedTime = 0
	g.TimeLabel.SetText("Time: 0s")

	// Start new timer
	g.Running = true
	g.StartTime = time.Now()
	g.Timer = time.NewTimer(time.Second)

	go func() {
		for {
			select {
			case <-g.Timer.C:
				if !g.Running {
					return
				}
				UpdateTimer(g)
				g.Timer.Reset(time.Second)
			}
		}
	}()
}

// StopTimer stops the game timer
func StopTimer(g *types.Game) {
	g.Running = false
	if g.Timer != nil {
		g.Timer.Stop()
	}
}
