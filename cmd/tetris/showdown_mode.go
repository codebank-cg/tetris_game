package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/oc-garden/tetris_game/internal/model"
	"github.com/rivo/tview"
)

// runShowdownMode launches the two-bot showdown mode.
// presetA and presetB are weight preset names (e.g. "aggressive", "conservative").
func runShowdownMode(app *tview.Application, presetA, presetB string) {
	ss := model.NewShowdownState(presetA, presetB)

	colors := map[int]tcell.Color{
		0: tcell.ColorDefault,
		1: tcell.GetColor("#00FFFF"),
		2: tcell.GetColor("#FFFF00"),
		3: tcell.GetColor("#FF00FF"),
		4: tcell.GetColor("#00FF00"),
		5: tcell.GetColor("#FF0000"),
		6: tcell.GetColor("#0000FF"),
		7: tcell.GetColor("#FFA500"),
	}

	// Board A (left)
	boardBoxA := tview.NewBox()
	boardBoxA.SetBackgroundColor(tcell.ColorBlack)
	boardBoxA.SetBorder(true)
	boardBoxA.SetBorderAttributes(tcell.AttrBold)
	boardBoxA.SetTitle(fmt.Sprintf(" BOT A: %s ", presetA)).SetTitleAlign(tview.AlignCenter)
	boardBoxA.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		drawBoard(screen, x, y, ss.BotA.State, colors, false)
		return x, y, width, height
	})

	// Stats panel (center)
	statsBox := tview.NewBox()
	statsBox.SetBackgroundColor(tcell.ColorBlack)
	statsBox.SetBorder(true)
	statsBox.SetBorderAttributes(tcell.AttrBold)
	statsBox.SetTitle(" SHOWDOWN ").SetTitleAlign(tview.AlignCenter)
	statsBox.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		drawStats(screen, x+1, y+1, ss, presetA, presetB)
		return x, y, width, height
	})

	// Board B (right)
	boardBoxB := tview.NewBox()
	boardBoxB.SetBackgroundColor(tcell.ColorBlack)
	boardBoxB.SetBorder(true)
	boardBoxB.SetBorderAttributes(tcell.AttrBold)
	boardBoxB.SetTitle(fmt.Sprintf(" BOT B: %s ", presetB)).SetTitleAlign(tview.AlignCenter)
	boardBoxB.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		drawBoard(screen, x, y, ss.BotB.State, colors, false)
		return x, y, width, height
	})

	// Layout: Board A (22) | Stats (26) | Board B (22) = 70 total
	mainLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(boardBoxA, 22, 0, false).
		AddItem(statsBox, 26, 0, false).
		AddItem(boardBoxB, 22, 0, false)
	mainLayout.SetBackgroundColor(tcell.ColorBlack)

	// Center vertically (boards are 22 rows tall including border)
	outerLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 0, 1, false).
		AddItem(mainLayout, 22, 0, false).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 0, 1, false)
	outerLayout.SetBackgroundColor(tcell.ColorBlack)

	app.SetRoot(outerLayout, true).EnableMouse(false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			app.Stop()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				app.Stop()
				return nil
			case 'r', 'R':
				ss.Reset()
				return nil
			case '+', '=':
				if ss.SpeedLevel < 5 {
					ss.SpeedLevel++
				}
				return nil
			case '-', '_':
				if ss.SpeedLevel > 1 {
					ss.SpeedLevel--
				}
				return nil
			}
		}
		return event
	})

	go func() {
		for {
			ss.Tick(time.Now())
			app.Draw()
			time.Sleep(16 * time.Millisecond)
		}
	}()
}

// drawStats renders the stats panel content at (x, y) — inner content area.
func drawStats(screen tcell.Screen, x, y int, ss *model.ShowdownState, presetA, presetB string) {
	white := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	yellow := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack)
	cyan := tcell.StyleDefault.Foreground(tcell.GetColor("#00FFFF")).Background(tcell.ColorBlack)
	green := tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
	red := tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlack)
	bold := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack).Bold(true)

	sa := ss.BotA.State
	sb := ss.BotB.State

	row := 0

	// Headers
	drawText(screen, x, y+row, fmt.Sprintf("%-12s%-12s", " BOT A", " BOT B"), bold)
	row++
	drawText(screen, x, y+row, "------------------------", white)
	row++

	// Score
	drawText(screen, x, y+row, "Score:", yellow)
	drawText(screen, x+7, y+row, fmt.Sprintf("%-7d", sa.Score), white)
	drawText(screen, x+14, y+row, fmt.Sprintf("%d", sb.Score), white)
	row++

	// Level
	drawText(screen, x, y+row, "Level:", yellow)
	drawText(screen, x+7, y+row, fmt.Sprintf("%-7d", sa.Level), white)
	drawText(screen, x+14, y+row, fmt.Sprintf("%d", sb.Level), white)
	row++

	// Lines
	drawText(screen, x, y+row, "Lines:", yellow)
	drawText(screen, x+7, y+row, fmt.Sprintf("%-7d", sa.LinesCleared), white)
	drawText(screen, x+14, y+row, fmt.Sprintf("%d", sb.LinesCleared), white)
	row++

	// Tetrises
	drawText(screen, x, y+row, "4-line:", yellow)
	drawText(screen, x+7, y+row, fmt.Sprintf("%-7d", sa.TetrisCount), white)
	drawText(screen, x+14, y+row, fmt.Sprintf("%d", sb.TetrisCount), white)
	row++

	// Height
	heightA := model.GetAggregateHeight(sa.Board) / 10
	heightB := model.GetAggregateHeight(sb.Board) / 10
	drawText(screen, x, y+row, "Height:", yellow)
	drawText(screen, x+7, y+row, fmt.Sprintf("%-7d", heightA), white)
	drawText(screen, x+14, y+row, fmt.Sprintf("%d", heightB), white)
	row++

	// Holes
	holesA := model.CountHoles(sa.Board)
	holesB := model.CountHoles(sb.Board)
	drawText(screen, x, y+row, "Holes:", yellow)
	drawText(screen, x+7, y+row, fmt.Sprintf("%-7d", holesA), white)
	drawText(screen, x+14, y+row, fmt.Sprintf("%d", holesB), white)
	row++

	drawText(screen, x, y+row, "------------------------", white)
	row++

	// Speed
	drawText(screen, x, y+row, fmt.Sprintf("Speed: %d  (+/- to change)", ss.SpeedLevel), cyan)
	row++

	drawText(screen, x, y+row, "------------------------", white)
	row++

	// Leading indicator
	if !ss.Running && ss.WinnerBot != "" {
		switch ss.WinnerBot {
		case "DRAW":
			drawText(screen, x, y+row, "*** DRAW! ***           ", bold)
		case "A":
			drawText(screen, x, y+row, "*** BOT A WINS! ***     ", green)
		case "B":
			drawText(screen, x, y+row, "*** BOT B WINS! ***     ", green)
		}
		row++
		drawText(screen, x, y+row, "Press R to restart      ", yellow)
		row++
	} else {
		diff := sa.Score - sb.Score
		if diff > 0 {
			drawText(screen, x, y+row, fmt.Sprintf("A LEADS by %d", diff), green)
		} else if diff < 0 {
			drawText(screen, x, y+row, fmt.Sprintf("B LEADS by %d", -diff), green)
		} else {
			drawText(screen, x, y+row, "TIED", white)
		}
		row++
		// Status
		statusA := "playing"
		if sa.GameOver {
			statusA = "DEAD"
		}
		statusB := "playing"
		if sb.GameOver {
			statusB = "DEAD"
		}
		drawText(screen, x, y+row, fmt.Sprintf("A:%-8sB:%-8s", statusA, statusB), white)
		row++
	}

	drawText(screen, x, y+row, "------------------------", white)
	row++

	// Session best
	if ss.Winner.Runs > 0 {
		drawText(screen, x, y+row, fmt.Sprintf("Best: Bot%s  %d", ss.Winner.BestBot, ss.Winner.BestScore), cyan)
	} else {
		drawText(screen, x, y+row, "Best: ---               ", cyan)
	}
	row++
	drawText(screen, x, y+row, fmt.Sprintf("Runs: %d", ss.Winner.Runs), white)
	row++

	drawText(screen, x, y+row, "------------------------", white)
	row++

	// Presets reminder
	drawText(screen, x, y+row, fmt.Sprintf("A:%-11sB:%s", presetA, presetB), white)
	row++

	// Controls
	drawText(screen, x, y+row, "R:restart  Q:quit       ", red)
}
