package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/wangle201210/sudoku/internal/game"
	"github.com/wangle201210/sudoku/internal/types"
	"github.com/wangle201210/sudoku/internal/ui"
)

func main() {
	// Create new application
	a := app.New()

	// Create main window
	w := a.NewWindow("Sudoku")

	// Create game state
	g := &types.Game{
		Size:       9,
		Difficulty: types.Easy,
		Board:      make([][]int, 9),
		Solution:   make([][]int, 9),
		Entries:    make([][]*types.SquareEntry, 9),
		Window:     w,
		Running:    false,
	}

	// Initialize arrays
	for i := range g.Board {
		g.Board[i] = make([]int, 9)
		g.Solution[i] = make([]int, 9)
		g.Entries[i] = make([]*types.SquareEntry, 9)
	}

	// Initialize timer label
	g.TimeLabel = widget.NewLabel("Time: 0s")

	// Generate initial board
	game.GenerateValidBoard(g)

	// Create initial game content
	gameContent := ui.CreateUI(g)

	// Create main container first
	mainContainer := container.NewVBox()

	// Create difficulty select
	difficultySelect := widget.NewSelect([]string{"Easy", "Medium", "Hard"}, func(selected string) {
		var difficulty types.Difficulty
		switch selected {
		case "Medium":
			difficulty = types.Medium
		case "Hard":
			difficulty = types.Hard
		default:
			difficulty = types.Easy
		}
		g.Difficulty = difficulty
		ui.StopTimer(g)
		game.GenerateValidBoard(g)
		ui.StartTimer(g)
		if len(mainContainer.Objects) > 1 {
			mainContainer.Objects[1] = ui.CreateUI(g)
			mainContainer.Refresh()
		}
	})

	// Create top toolbar with mainContainer reference
	toolbar := container.NewHBox(
		difficultySelect,
		widget.NewButton("New Game", func() {
			ui.StopTimer(g)
			game.GenerateValidBoard(g)
			ui.StartTimer(g)
			if len(mainContainer.Objects) > 1 {
				mainContainer.Objects[1] = ui.CreateUI(g)
				mainContainer.Refresh()
			}
		}),
		widget.NewButton("Check", func() {
			if game.CheckWin(g) {
				dialog.ShowInformation("Congratulations", "You won!", g.Window)
				ui.StopTimer(g)
			} else {
				dialog.ShowInformation("Info", "Keep going!", g.Window)
			}
		}),
	)

	// Add components to main container
	mainContainer.Add(toolbar)
	mainContainer.Add(gameContent)

	// Set default difficulty
	difficultySelect.SetSelected("Easy")

	// Start initial timer
	ui.StartTimer(g)

	// Set window content
	w.SetContent(mainContainer)

	// Set window size
	w.Resize(fyne.NewSize(400, 500))

	// Show and run
	w.ShowAndRun()
}
