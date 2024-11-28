package types

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"time"
)

// Difficulty represents the game difficulty level
type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

// GameState represents the Sudoku game state
type GameState struct {
	Board      [][]int    `json:"board"`
	Time       int        `json:"time"`
	Difficulty Difficulty `json:"difficulty"`
}

// Game represents the Sudoku game structure
type Game struct {
	Size        int
	Difficulty  Difficulty
	Board       [][]int
	Solution    [][]int
	Entries     [][]*SquareEntry
	Window      fyne.Window
	TimeLabel   *widget.Label
	Timer       *time.Timer
	StartTime   time.Time
	Running     bool
	Hints       int
	ElapsedTime int
}

// SquareEntry is a custom entry widget that maintains a square shape
type SquareEntry struct {
	widget.Entry
}

func NewSquareEntry() *SquareEntry {
	entry := &SquareEntry{}
	entry.ExtendBaseWidget(entry)
	entry.TextStyle = fyne.TextStyle{Bold: true}
	return entry
}

func (e *SquareEntry) MinSize() fyne.Size {
	s := e.Entry.MinSize()
	if s.Width > s.Height {
		return fyne.NewSize(s.Width, s.Width)
	}
	return fyne.NewSize(s.Height, s.Height)
}
