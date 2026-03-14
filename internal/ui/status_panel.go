package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

// StatusPanel displays game statistics.
type StatusPanel struct {
	*tview.Flex
	scoreText *tview.TextView
	levelText *tview.TextView
	linesText *tview.TextView
	nextPiece *tview.TextView
}

// NewStatusPanel creates a new status panel.
func NewStatusPanel() *StatusPanel {
	sp := &StatusPanel{
		Flex:      tview.NewFlex().SetDirection(tview.FlexRow),
		scoreText: tview.NewTextView().SetText("Score: 0"),
		levelText: tview.NewTextView().SetText("Level: 1"),
		linesText: tview.NewTextView().SetText("Lines: 0"),
		nextPiece: tview.NewTextView().SetText("Next:\n"),
	}

	sp.AddItem(sp.scoreText, 1, 1, false)
	sp.AddItem(sp.levelText, 1, 1, false)
	sp.AddItem(sp.linesText, 1, 1, false)
	sp.AddItem(sp.nextPiece, 6, 1, false)

	return sp
}

// Update updates the status panel with current game state.
func (sp *StatusPanel) Update(score, level, lines int) {
	sp.scoreText.SetText(fmt.Sprintf("Score: %d", score))
	sp.levelText.SetText(fmt.Sprintf("Level: %d", level))
	sp.linesText.SetText(fmt.Sprintf("Lines: %d", lines))
}

// SetNextPiece displays the next piece preview.
func (sp *StatusPanel) SetNextPiece(piece string) {
	sp.nextPiece.SetText("Next:\n" + piece)
}
