package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// GameBoard displays the Tetris game board.
type GameBoard struct {
	*tview.Box
	colors map[int]tcell.Color
}

// NewGameBoard creates a new game board component.
func NewGameBoard() *GameBoard {
	gb := &GameBoard{
		Box: tview.NewBox(),
		colors: map[int]tcell.Color{
			0: tcell.ColorDefault,
			1: tcell.GetColor("cyan"),
			2: tcell.GetColor("yellow"),
			3: tcell.GetColor("purple"),
			4: tcell.GetColor("green"),
			5: tcell.GetColor("red"),
			6: tcell.GetColor("blue"),
			7: tcell.GetColor("orange"),
		},
	}
	return gb
}

// SetGameState sets the game state for rendering.
func (gb *GameBoard) SetGameState(state interface{}) {
	// Game state is managed externally, board just renders
}

// Draw renders the game board.
func (gb *GameBoard) Draw(screen tcell.Screen) {
	gb.Box.Draw(screen)
}
