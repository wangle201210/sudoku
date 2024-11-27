package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

type GameState struct {
	Board      [][]int    `json:"board"`
	Size       int        `json:"size"`
	Time       int        `json:"time"`
	Difficulty Difficulty `json:"difficulty"`
}

type SudokuGame struct {
	size        int
	difficulty  Difficulty
	board       [][]int
	solution    [][]int
	entries     [][]*widget.Entry
	window      fyne.Window
	timeLabel   *widget.Label
	timer       *time.Timer
	startTime   time.Time
	running     bool
	hints       int
	elapsedTime int
}

func newGame(difficulty Difficulty) *SudokuGame {
	g := &SudokuGame{
		size:       9,
		difficulty: difficulty,
		board:      make([][]int, 9),
		solution:   make([][]int, 9),
		entries:    make([][]*widget.Entry, 9),
		running:    false,
		hints:      3,
	}

	for i := range g.board {
		g.board[i] = make([]int, 9)
		g.solution[i] = make([]int, 9)
		g.entries[i] = make([]*widget.Entry, 9)
	}

	g.generateValidBoard()
	return g
}

func (g *SudokuGame) generateValidBoard() {
	g.generateLargeSolution()

	for i := range g.board {
		copy(g.board[i], g.solution[i])
	}

	cellsToRemove := 0
	switch g.difficulty {
	case Easy:
		cellsToRemove = 40
	case Medium:
		cellsToRemove = 50
	case Hard:
		cellsToRemove = 60
	}

	removed := 0
	maxAttempts := g.size * g.size * 2

	for removed < cellsToRemove && maxAttempts > 0 {
		row := rand.Intn(g.size)
		col := rand.Intn(g.size)

		if g.board[row][col] != 0 {
			temp := g.board[row][col]
			g.board[row][col] = 0

			if g.hasUniqueSolution() {
				removed++
			} else {
				g.board[row][col] = temp
			}
		}
		maxAttempts--
	}
}

func (g *SudokuGame) generateLargeSolution() {
	for i := range g.solution {
		for j := range g.solution[i] {
			g.solution[i][j] = 0
		}
	}

	g.fillSudoku(g.solution, 0, 0)
}

func (g *SudokuGame) fillSudoku(board [][]int, row, col int) bool {
	if row == g.size {
		return true
	}

	nextRow := row
	nextCol := col + 1
	if nextCol == g.size {
		nextRow = row + 1
		nextCol = 0
	}

	if board[row][col] != 0 {
		return g.fillSudoku(board, nextRow, nextCol)
	}

	nums := make([]int, 9)
	for i := range nums {
		nums[i] = i + 1
	}
	rand.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})

	for _, num := range nums {
		if g.isValidMove(board, row, col, num) {
			board[row][col] = num
			if g.fillSudoku(board, nextRow, nextCol) {
				return true
			}
			board[row][col] = 0
		}
	}

	return false
}

func (g *SudokuGame) isValidMove(board [][]int, row, col, num int) bool {
	for j := 0; j < g.size; j++ {
		if board[row][j] == num {
			return false
		}
	}

	for i := 0; i < g.size; i++ {
		if board[i][col] == num {
			return false
		}
	}

	startRow := (row / 3) * 3
	startCol := (col / 3) * 3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[startRow+i][startCol+j] == num {
				return false
			}
		}
	}

	return true
}

func (g *SudokuGame) hasUniqueSolution() bool {
	tempBoard := make([][]int, g.size)
	for i := range tempBoard {
		tempBoard[i] = make([]int, g.size)
		copy(tempBoard[i], g.board[i])
	}

	if !g.solveSudoku(tempBoard, 0, 0) {
		return false
	}

	return !g.hasAnotherSolution(g.board, 0, 0, tempBoard)
}

func (g *SudokuGame) hasAnotherSolution(board [][]int, row, col int, firstSolution [][]int) bool {
	if row == g.size {
		for i := 0; i < g.size; i++ {
			for j := 0; j < g.size; j++ {
				if board[i][j] != firstSolution[i][j] {
					return true
				}
			}
		}
		return false
	}

	nextRow := row
	nextCol := col + 1
	if nextCol == g.size {
		nextRow = row + 1
		nextCol = 0
	}

	if board[row][col] != 0 {
		return g.hasAnotherSolution(board, nextRow, nextCol, firstSolution)
	}

	for num := 1; num <= g.size; num++ {
		if g.isValidMove(board, row, col, num) {
			board[row][col] = num
			if g.hasAnotherSolution(board, nextRow, nextCol, firstSolution) {
				return true
			}
			board[row][col] = 0
		}
	}

	return false
}

func (g *SudokuGame) solveSudoku(board [][]int, row, col int) bool {
	if row == g.size {
		return true
	}

	nextRow := row
	nextCol := col + 1
	if nextCol == g.size {
		nextRow = row + 1
		nextCol = 0
	}

	if board[row][col] != 0 {
		return g.solveSudoku(board, nextRow, nextCol)
	}

	for num := 1; num <= g.size; num++ {
		if g.isValidMove(board, row, col, num) {
			board[row][col] = num
			if g.solveSudoku(board, nextRow, nextCol) {
				return true
			}
			board[row][col] = 0
		}
	}

	return false
}

func (g *SudokuGame) checkWin() bool {
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			if g.board[i][j] == 0 {
				return false
			}
		}
	}

	for i := 0; i < g.size; i++ {
		seen := make(map[int]bool)
		for j := 0; j < g.size; j++ {
			if seen[g.board[i][j]] {
				return false
			}
			seen[g.board[i][j]] = true
		}
	}

	for j := 0; j < g.size; j++ {
		seen := make(map[int]bool)
		for i := 0; i < g.size; i++ {
			if seen[g.board[i][j]] {
				return false
			}
			seen[g.board[i][j]] = true
		}
	}

	boxSize := 3
	for blockRow := 0; blockRow < g.size; blockRow += boxSize {
		for blockCol := 0; blockCol < g.size; blockCol += boxSize {
			seen := make(map[int]bool)
			for i := 0; i < boxSize; i++ {
				for j := 0; j < boxSize; j++ {
					val := g.board[blockRow+i][blockCol+j]
					if seen[val] {
						return false
					}
					seen[val] = true
				}
			}
		}
	}

	return true
}

func (g *SudokuGame) giveHint() {
	if g.hints <= 0 {
		dialog.ShowInformation("No Hints", "You have used all your hints!", g.window)
		return
	}

	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			if g.board[i][j] == 0 {
				g.board[i][j] = g.solution[i][j]
				g.entries[i][j].SetText(fmt.Sprintf("%d", g.solution[i][j]))
				g.entries[i][j].Disable()
				g.hints--
				return
			}
		}
	}
}

func (g *SudokuGame) saveGame() {
	state := GameState{
		Board:      g.board,
		Size:       g.size,
		Time:       g.elapsedTime,
		Difficulty: g.difficulty,
	}

	data, err := json.Marshal(state)
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}

	savePath := filepath.Join(homeDir, "sudoku_save.json")
	err = ioutil.WriteFile(savePath, data, 0644)
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}

	dialog.ShowInformation("Success", "Game saved successfully!", g.window)
}

func (g *SudokuGame) loadGame() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		dialog.ShowError(err, g.window)
		return false
	}

	savePath := filepath.Join(homeDir, "sudoku_save.json")
	data, err := ioutil.ReadFile(savePath)
	if err != nil {
		dialog.ShowError(err, g.window)
		return false
	}

	var state GameState
	err = json.Unmarshal(data, &state)
	if err != nil {
		dialog.ShowError(err, g.window)
		return false
	}

	g.board = state.Board
	g.size = state.Size
	g.elapsedTime = state.Time
	g.difficulty = state.Difficulty

	return true
}

func (g *SudokuGame) startTimer() {
	g.startTime = time.Now()
	g.running = true

	g.timer = time.NewTimer(time.Second)
	go func() {
		for g.running {
			<-g.timer.C
			if g.running {
				g.updateTimer()
				g.timer.Reset(time.Second)
			}
		}
	}()
}

func (g *SudokuGame) stopTimer() {
	g.running = false
	if g.timer != nil {
		g.timer.Stop()
	}
}

func (g *SudokuGame) updateTimer() {
	if g.timeLabel != nil {
		elapsed := time.Since(g.startTime)
		minutes := int(elapsed.Minutes())
		seconds := int(elapsed.Seconds()) % 60

		g.elapsedTime = minutes*60 + seconds

		if g.window != nil {
			g.window.Canvas().Refresh(g.timeLabel)
			g.timeLabel.SetText(fmt.Sprintf("Time: %d:%02d", minutes, seconds))
		}
	}
}

func (g *SudokuGame) createUI() fyne.CanvasObject {
	difficultySelect := widget.NewSelect([]string{"Easy", "Medium", "Hard"}, nil)
	difficultySelect.SetSelected("Medium")

	mainGrid := container.NewGridWithColumns(1)

	cellSize := 40

	for bigRow := 0; bigRow < 9; bigRow += 3 {
		rowContainer := container.NewHBox()
		for bigCol := 0; bigCol < 9; bigCol += 3 {
			subGrid := container.NewGridWithColumns(3)

			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					row := bigRow + i
					col := bigCol + j

					entry := widget.NewEntry()
					entry.Resize(fyne.NewSize(float32(cellSize), float32(cellSize)))
					entry.TextStyle = fyne.TextStyle{Bold: true}

					if g.board[row][col] != 0 {
						entry.SetText(fmt.Sprintf("%d", g.board[row][col]))
						entry.Disable()
					}

					g.entries[row][col] = entry

					final_row, final_col := row, col
					entry.OnChanged = func(text string) {
						if text == "" {
							g.board[final_row][final_col] = 0
							return
						}

						num, err := strconv.Atoi(text)
						if err != nil || num < 1 || num > 9 {
							entry.SetText("")
							return
						}

						if !g.isValid(final_row, final_col, num) {
							entry.SetText("")
							return
						}

						g.board[final_row][final_col] = num

						if g.checkWin() {
							g.stopTimer()
							dialog.ShowInformation("Congratulations!", "You've solved the puzzle!", g.window)
						}
					}

					centeredEntry := container.NewCenter(entry)
					border := container.NewBorder(nil, nil, nil, nil, centeredEntry)
					subGrid.Add(container.NewPadded(border))
				}
			}

			borderedSubGrid := container.NewBorder(nil, nil, nil, nil,
				container.NewPadded(
					canvas.NewRectangle(theme.DisabledButtonColor()),
					subGrid,
				),
			)
			rowContainer.Add(borderedSubGrid)
		}
		mainGrid.Add(rowContainer)
	}

	newGameBtn := widget.NewButtonWithIcon("New Game", theme.MediaPlayIcon(), func() {
		g.stopTimer()
		newG := newGame(g.difficulty)
		newG.window = g.window
		g.window.SetContent(newG.createUI())
		newG.startTimer()
	})

	difficultySelect.OnChanged = func(d string) {
		var newDiff Difficulty
		switch d {
		case "Easy":
			newDiff = Easy
		case "Hard":
			newDiff = Hard
		default:
			newDiff = Medium
		}
		g.stopTimer()
		newG := newGame(newDiff)
		newG.window = g.window
		g.window.SetContent(newG.createUI())
		newG.startTimer()
	}

	g.timeLabel = widget.NewLabel("Time: 0:00")

	// hintBtn := widget.NewButtonWithIcon("Hint", theme.HelpIcon(), func() {
	// 	g.giveHint()
	// })

	content := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Difficulty:"),
			difficultySelect,
			newGameBtn,
			// hintBtn,
			g.timeLabel,
		),
		mainGrid,
	)

	g.startTimer()

	return content
}

func (g *SudokuGame) isValid(row, col, num int) bool {
	for x := 0; x < g.size; x++ {
		if x != col && g.board[row][x] == num {
			return false
		}
	}

	for x := 0; x < g.size; x++ {
		if x != row && g.board[x][col] == num {
			return false
		}
	}

	boxSize := 3
	startRow := (row / boxSize) * boxSize
	startCol := (col / boxSize) * boxSize

	return g.isValidBox(startRow, startCol, num)
}

func (g *SudokuGame) isValidBox(startRow, startCol, num int) bool {
	boxSize := 3
	endRow := startRow + boxSize
	endCol := startCol + boxSize

	if endRow > g.size {
		endRow = g.size
	}
	if endCol > g.size {
		endCol = g.size
	}

	for i := startRow; i < endRow; i++ {
		for j := startCol; j < endCol; j++ {
			if g.board[i][j] == num {
				return false
			}
		}
	}
	return true
}

func main() {
	rand.Seed(time.Now().UnixNano())

	myApp := app.New()
	myWindow := myApp.NewWindow("Sudoku")
	myWindow.Resize(fyne.NewSize(400, 500))

	game := newGame(Medium)
	game.window = myWindow

	myWindow.SetContent(game.createUI())
	myWindow.ShowAndRun()
}
