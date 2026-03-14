package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/oc-garden/tetris_game/internal/model"
)

// RenderAutoPlayIndicator displays AUTO-PLAY status on screen.
// Positioned to left of game board (x=0) to avoid overlap.
func RenderAutoPlayIndicator(screen tcell.Screen, autoPlayer *model.AutoPlayer) {
	style := tcell.StyleDefault.Foreground(tcell.ColorGray)
	text := "AUTO-PLAY: OFF"

	if autoPlayer.IsEnabled() {
		style = tcell.StyleDefault.Foreground(tcell.ColorLime).Bold(true)
		text = "AUTO-PLAY: ON"
	}

	// Render text vertically at x=0 (left of board border at x=2)
	for i, r := range []rune(text) {
		screen.SetContent(0, i+1, r, nil, style)
	}
}

// RenderSpeedLevel displays current speed level.
// Positioned to left of game board (x=0) to avoid overlap.
func RenderSpeedLevel(screen tcell.Screen, autoPlayer *model.AutoPlayer) {
	style := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	text := "SPEED: "

	if autoPlayer.IsEnabled() {
		text += string('0' + byte(autoPlayer.GetSpeedLevel()))
	} else {
		text += "-"
	}

	// Render text vertically at x=0 (left of board border at x=2)
	for i, r := range []rune(text) {
		screen.SetContent(0, i+8, r, nil, style)
	}
}

// RenderDecisionPanel displays AI's current decision.
func RenderDecisionPanel(screen tcell.Screen, decision *model.MoveDecision) {
	baseX := 26
	baseY := 20

	style := tcell.StyleDefault.Foreground(tcell.ColorAqua)

	if decision == nil {
		text := "AI: Thinking..."
		for i, r := range []rune(text) {
			screen.SetContent(baseX, baseY+i, r, nil, style)
		}
		return
	}

	texts := []string{
		"AI Decision:",
		"Target X:",
		"Rotation:",
		"Drops:",
		"Score:",
	}

	for i, text := range texts {
		for j, r := range []rune(text) {
			screen.SetContent(baseX+j, baseY+i, r, nil, style)
		}
	}

	values := []string{
		"",
		string('0' + byte(decision.GetTargetX())),
		string('0' + byte(decision.GetRotations())),
		string('0' + byte(decision.GetSoftDrops())),
	}

	for i, val := range values[1:] {
		for j, r := range []rune(val) {
			screen.SetContent(baseX+12+j, baseY+i+1, r, nil, style)
		}
	}
}
