package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// AutoPlayerPanel displays autoplay status and AI decisions.
type AutoPlayerPanel struct {
	*tview.Flex
	autoText     *tview.TextView
	speedText    *tview.TextView
	decisionView *tview.TextView
}

// NewAutoPlayerPanel creates a new autoplayer panel.
func NewAutoPlayerPanel() *AutoPlayerPanel {
	app := &AutoPlayerPanel{
		Flex:         tview.NewFlex().SetDirection(tview.FlexRow),
		autoText:     tview.NewTextView().SetText("AUTO-PLAY: OFF").SetTextColor(tcell.ColorGray),
		speedText:    tview.NewTextView().SetText("SPEED: -").SetTextColor(tcell.ColorYellow),
		decisionView: tview.NewTextView(),
	}

	app.AddItem(app.autoText, 1, 1, false)
	app.AddItem(app.speedText, 1, 1, false)
	app.AddItem(app.decisionView, 5, 1, false)

	return app
}

// UpdateAuto updates the autoplay status.
func (app *AutoPlayerPanel) UpdateAuto(enabled bool, speedLevel int) {
	if enabled {
		app.autoText.SetText("AUTO-PLAY: ON").SetTextColor(tcell.ColorLime)
		app.speedText.SetText(fmt.Sprintf("SPEED: %d", speedLevel))
	} else {
		app.autoText.SetText("AUTO-PLAY: OFF").SetTextColor(tcell.ColorGray)
		app.speedText.SetText("SPEED: -")
	}
}

// UpdateDecision displays the AI's current decision.
func (app *AutoPlayerPanel) UpdateDecision(targetX, rotations, drops int, score float64) {
	if targetX < 0 {
		app.decisionView.SetText("AI: Thinking...")
	} else {
		app.decisionView.SetText(
			fmt.Sprintf("Target X: %d\nRotation: %d\nDrops: %d\nScore: %.1f",
				targetX, rotations, drops, score))
	}
}
