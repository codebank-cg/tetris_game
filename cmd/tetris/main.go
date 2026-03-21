package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/oc-garden/tetris_game/internal/audio"
	"github.com/oc-garden/tetris_game/internal/model"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	musicPlayer := audio.NewMusicPlayer()
	go musicPlayer.PlayKorobeiniki()

	var mu sync.Mutex // protects game, ghostEnabled, gameOverDetected
	game := model.NewGameState()
	autoPlayer := model.NewAutoPlayer()
	ghostEnabled := false

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

	gameBoard := tview.NewBox()
	gameBoard.SetBorder(true)
	gameBoard.SetBorderAttributes(tcell.AttrBold)
	gameBoard.SetTitle(" TETRIS ").SetTitleAlign(tview.AlignCenter)
	gameBoard.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		boardX := x + 1
		boardY := y + 1

		mu.Lock()
		g := game
		ge := ghostEnabled
		mu.Unlock()

		for row := 0; row < 20; row++ {
			screenY := boardY + (19 - row)
			for col := 0; col < 10; col++ {
				cellX := boardX + col*2
				color := g.Board.Get(col, row)
				if color != 0 {
					baseColor := colors[color]
					screen.SetContent(cellX, screenY, '█', nil, tcell.StyleDefault.Foreground(baseColor).Background(tcell.ColorBlack))
					screen.SetContent(cellX+1, screenY, '█', nil, tcell.StyleDefault.Foreground(baseColor).Background(tcell.ColorBlack))
				}
			}
		}

		if g.CurrentPiece != nil && !g.GameOver && ge && !autoPlayer.IsEnabled() {
			ghostY := g.GetGhostY()
			if ghostY >= 0 && ghostY <= g.CurrentPiece.Y && ghostY < 20 {
				matrix := g.CurrentPiece.GetMatrix()
				if matrix != nil {
					ghostStyle := tcell.StyleDefault.Foreground(tcell.GetColor("#90EE90")).Background(tcell.ColorBlack)
					for row := 0; row < 4; row++ {
						for col := 0; col < 4; col++ {
							if matrix[row][col] != 0 {
								pieceY := ghostY - row
								if pieceY >= 0 && pieceY < 20 {
									bx := boardX + (g.CurrentPiece.X+col)*2
									by := boardY + (19 - pieceY)
									screen.SetContent(bx, by, '░', nil, ghostStyle)
									screen.SetContent(bx+1, by, '░', nil, ghostStyle)
								}
							}
						}
					}
				}
			}
		}

		if g.CurrentPiece != nil && !g.GameOver {
			matrix := g.CurrentPiece.GetMatrix()
			pieceStyle := tcell.StyleDefault.Foreground(colors[g.CurrentPiece.Color]).Background(tcell.ColorBlack)
			for row := 0; row < 4; row++ {
				for col := 0; col < 4; col++ {
					if matrix[row][col] != 0 {
						bx := boardX + (g.CurrentPiece.X+col)*2
						by := boardY + (19 - (g.CurrentPiece.Y - row))
						if by >= boardY && by < boardY+20 && bx >= boardX && bx < boardX+20 {
							screen.SetContent(bx, by, '█', nil, pieceStyle)
							screen.SetContent(bx+1, by, '█', nil, pieceStyle)
						}
					}
				}
			}
		}

		if g.GameOver {
			gameOverText := "GAME OVER"
			startX := boardX + (20-len(gameOverText))/2
			startY := boardY + 9
			for i, ch := range gameOverText {
				style := tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlack).Bold(true)
				screen.SetContent(startX+i, startY, ch, nil, style)
			}

			restartText := "Press R"
			startX = boardX + (20-len(restartText))/2
			for i, ch := range restartText {
				style := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack)
				screen.SetContent(startX+i, startY+2, ch, nil, style)
			}
		}

		return x, y, width, height
	})

	nextBox := tview.NewBox()
	nextBox.SetBorder(true)
	nextBox.SetBorderAttributes(tcell.AttrBold)
	nextBox.SetTitle(" NEXT ").SetTitleAlign(tview.AlignCenter)
	nextBox.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		mu.Lock()
		g := game
		mu.Unlock()

		if g.NextPiece == nil {
			return x, y, width, height
		}

		boxWidth := width - 2
		boxHeight := height - 2
		matrix := g.NextPiece.GetMatrix()
		pieceColor := colors[g.NextPiece.Color]
		pieceStyle := tcell.StyleDefault.Foreground(pieceColor).Background(tcell.ColorBlack)

		minCol, maxCol := 4, 0
		minRow, maxRow := 4, 0
		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				if matrix[row][col] != 0 {
					if col < minCol {
						minCol = col
					}
					if col > maxCol {
						maxCol = col
					}
					if row < minRow {
						minRow = row
					}
					if row > maxRow {
						maxRow = row
					}
				}
			}
		}

		pieceWidth := (maxCol - minCol + 1) * 2
		pieceHeight := maxRow - minRow + 1
		offsetX := (boxWidth - pieceWidth) / 2
		offsetY := (boxHeight - pieceHeight) / 2

		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				if matrix[row][col] != 0 {
					screenX := x + 1 + offsetX + (col-minCol)*2
					screenY := y + 1 + offsetY + (maxRow - row)
					screen.SetContent(screenX, screenY, '█', nil, pieceStyle)
					screen.SetContent(screenX+1, screenY, '█', nil, pieceStyle)
				}
			}
		}

		return x, y, width, height
	})

	infoBox := tview.NewTextView().SetDynamicColors(true)
	infoBox.SetTextColor(tcell.ColorWhite)
	infoBox.SetBackgroundColor(tcell.ColorBlack)
	infoBox.SetText(
		"[#FFFF00::b]Score:[::-] 0\n" +
			"[#FFFF00::b]Level:[::-] 1\n" +
			"[#FFFF00::b]Lines:[::-] 0\n\n" +
			"[::b]Controls:[::-]\n" +
			"[#FFFF00]←→[white] Move\n" +
			"[#FFFF00]↑[white] Rotate\n" +
			"[#FFFF00]Z[white] Rotate CCW\n" +
			"[#FFFF00]↓[white] Soft Drop\n" +
			"[#FFFF00]Space[white] Hard Drop\n" +
			"[#FFFF00]G[white] Ghost[white] Mode\n" +
			"[#FFFF00]P[white] Pause\n" +
			"[#FFFF00]R[white] Restart\n" +
			"[#FFFF00]Q[white] Quit",
	)
	infoBox.SetBorder(true)
	infoBox.SetBorderAttributes(tcell.AttrBold)
	infoBox.SetTitle(" INFO ")

	autoText := tview.NewTextView()
	autoText.SetTextColor(tcell.ColorWhite)
	autoText.SetBackgroundColor(tcell.ColorBlack)
	autoText.SetText("AUTO-PLAY: OFF")

	speedText := tview.NewTextView()
	speedText.SetTextColor(tcell.ColorYellow)
	speedText.SetBackgroundColor(tcell.ColorBlack)
	speedText.SetText("SPEED:   -")

	ghostText := tview.NewTextView()
	ghostText.SetTextColor(tcell.GetColor("#00FFFF"))
	ghostText.SetBackgroundColor(tcell.ColorBlack)
	ghostText.SetText("GHOST: OFF")

	autoParamsText := tview.NewTextView()
	autoParamsText.SetTextColor(tcell.GetColor("#00CED1"))
	autoParamsText.SetBackgroundColor(tcell.ColorBlack)
	autoParamsText.SetText("")

	leftPanelContent := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(autoText, 2, 0, false).
		AddItem(speedText, 2, 0, false).
		AddItem(ghostText, 2, 0, false).
		AddItem(autoParamsText, 5, 0, false)
	leftPanelContent.SetBackgroundColor(tcell.ColorBlack)

	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(leftPanelContent, 0, 1, false)
	leftPanel.SetBorder(true)
	leftPanel.SetBorderAttributes(tcell.AttrBold)
	leftPanel.SetTitle(" AUTO-PLAY ")
	leftPanel.SetBackgroundColor(tcell.ColorBlack)

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nextBox, 6, 0, false).
		AddItem(infoBox, 16, 0, false)
	rightPanel.SetBackgroundColor(tcell.ColorBlack)

	rightPanelWithSpacing := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(rightPanel, 22, 0, false).
		AddItem(tview.NewBox(), 0, 1, false)
	rightPanelWithSpacing.SetBackgroundColor(tcell.ColorBlack)

	mainLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftPanel, 22, 0, false).
		AddItem(gameBoard, 22, 0, false).
		AddItem(rightPanelWithSpacing, 18, 0, false)
	mainLayout.SetBackgroundColor(tcell.ColorBlack)

	mainLayoutWithSpacing := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(mainLayout, 22, 0, false).
		AddItem(tview.NewBox(), 0, 1, false)
	mainLayoutWithSpacing.SetBackgroundColor(tcell.ColorBlack)

	app.SetRoot(mainLayoutWithSpacing, true).EnableMouse(false)

	lastDrop := time.Now()
	lastActionTime := time.Now()

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
			case 'a', 'A':
				autoPlayer.Toggle()
				if autoPlayer.IsEnabled() {
					autoText.SetText("AUTO-PLAY: ON ")
					speedText.SetText(fmt.Sprintf("SPEED:   %d", autoPlayer.GetSpeedLevel()))
					mu.Lock()
					ghostEnabled = false
					mu.Unlock()
					ghostText.SetText("GHOST: OFF")
				} else {
					autoText.SetText("AUTO-PLAY: OFF")
					speedText.SetText("SPEED:   -")
					autoParamsText.SetText("")
				}
				return nil
			case 'z', 'Z':
				mu.Lock()
				g := game
				mu.Unlock()
				if !g.Paused && !g.GameOver {
					g.RotatePieceCounter()
				}
				return nil
			case 'g', 'G':
				if !autoPlayer.IsEnabled() {
					mu.Lock()
					ghostEnabled = !ghostEnabled
					ge := ghostEnabled
					mu.Unlock()
					if ge {
						ghostText.SetText("GHOST: ON ")
					} else {
						ghostText.SetText("GHOST: OFF")
					}
				}
				return nil
			case ' ':
				mu.Lock()
				g := game
				mu.Unlock()
				if !g.Paused && !g.GameOver {
					dropped := g.DropPiece()
					// NES Tetris: 2 points per row hard-dropped
					g.Score += dropped * 2
				}
				return nil
			case 'p', 'P':
				mu.Lock()
				game.Paused = !game.Paused
				mu.Unlock()
				app.Draw()
				return nil
			case 'r', 'R':
				mu.Lock()
				game = model.NewGameState()
				mu.Unlock()
				musicPlayer.Stop()
				musicPlayer.Restart()
				go musicPlayer.PlayKorobeiniki()
				if autoPlayer.IsEnabled() {
					autoText.SetText("AUTO-PLAY: ON ")
					speedText.SetText(fmt.Sprintf("SPEED:   %d", autoPlayer.GetSpeedLevel()))
					mu.Lock()
					ghostEnabled = false
					mu.Unlock()
					ghostText.SetText("GHOST: OFF")
				} else {
					autoText.SetText("AUTO-PLAY: OFF")
					speedText.SetText("SPEED:   -")
					mu.Lock()
					if !ghostEnabled {
						mu.Unlock()
						ghostText.SetText("GHOST: OFF")
					} else {
						mu.Unlock()
					}
				}
				autoParamsText.SetText("TARGET:  X: -  R:-\nSCORE:   -. -")
				return nil
			}
		case tcell.KeyLeft:
			mu.Lock()
			g := game
			mu.Unlock()
			if !g.Paused && !g.GameOver {
				g.MovePiece(-1, 0)
			}
			return nil
		case tcell.KeyRight:
			mu.Lock()
			g := game
			mu.Unlock()
			if !g.Paused && !g.GameOver {
				g.MovePiece(1, 0)
			}
			return nil
		case tcell.KeyDown:
			mu.Lock()
			g := game
			mu.Unlock()
			if !g.Paused && !g.GameOver {
				g.SoftDrop()
				lastDrop = time.Now()
			}
			return nil
		case tcell.KeyUp:
			mu.Lock()
			g := game
			mu.Unlock()
			if !g.Paused && !g.GameOver {
				g.RotatePiece()
			}
			return nil
		}
		return event
	})

	go func() {
		gameOverDetected := false
		lastScore := 0
		lastLevel := 0
		lastLines := 0
		lastPiecePos := 0
		lastAutoParams := ""
		lastAutoPlayerState := false
		currentDecision := (*model.MoveDecision)(nil)

		for {
			boardChanged := false

			if autoPlayer.IsEnabled() != lastAutoPlayerState {
				lastAutoPlayerState = autoPlayer.IsEnabled()
				if autoPlayer.IsEnabled() {
					lastActionTime = time.Now()
					boardChanged = true
				}
			}

			mu.Lock()
			g := game
			mu.Unlock()

			if autoPlayer.IsEnabled() && !g.Paused && !g.GameOver && !g.IsClearAnimating() && g.CurrentPiece != nil {
				delay := model.GetDelayForSpeed(g.GetDropInterval(), autoPlayer.GetSpeedLevel())
				if model.IsInDropPhase(g) {
					delay = delay / 3
				}
				if time.Since(lastActionTime) > time.Duration(delay)*time.Millisecond {
					if currentDecision == nil {
						currentDecision = model.FindBestMoveWithNext(g)
					}
					if currentDecision != nil {
						if model.ExecuteMove(g, currentDecision) {
							currentDecision = nil
						}
						lastActionTime = time.Now()
						boardChanged = true
					} else {
						g.SoftDrop()
						lastActionTime = time.Now()
						boardChanged = true
					}
				}
			}

			if !autoPlayer.IsEnabled() && time.Since(lastDrop) > time.Duration(g.GetDropInterval())*time.Millisecond && !g.Paused && !g.GameOver && !g.IsClearAnimating() {
				g.SoftDrop()
				lastDrop = time.Now()
				boardChanged = true
			}

			if g.IsClearAnimating() {
				linesCleared := g.UpdateClearAnimation()
				if linesCleared {
					go musicPlayer.PlayLineClearBeep()
					boardChanged = true
				}
			}

			mu.Lock()
			if g.GameOver && !gameOverDetected {
				gameOverDetected = true
				mu.Unlock()
				musicPlayer.Stop()
				boardChanged = true
			} else {
				mu.Unlock()
			}

			scoreChanged := g.Score != lastScore || g.Level != lastLevel || g.LinesCleared != lastLines
			currentPiecePos := 0
			pieceMoved := false
			if g.CurrentPiece != nil {
				currentPiecePos = g.CurrentPiece.X*100 + g.CurrentPiece.Y
				pieceMoved = currentPiecePos != lastPiecePos
			}

			// Reuse currentDecision for display — avoid redundant FindBestMoveWithNext call
			if autoPlayer.IsEnabled() && !g.GameOver && g.CurrentPiece != nil {
				decision := currentDecision
				if decision == nil {
					decision = model.FindBestMoveWithNext(g)
				}
				if decision != nil {
					autoParams := fmt.Sprintf("TARGET:  X:%2d  R:%d\nSCORE:   %.1f",
						decision.GetTargetX(), decision.GetRotations(), decision.GetScore())
					if autoParams != lastAutoParams {
						autoParamsText.SetText(autoParams)
						lastAutoParams = autoParams
					}
				} else {
					if lastAutoParams != "NO VALID MOVE" {
						autoParamsText.SetText("NO VALID MOVE")
						lastAutoParams = "NO VALID MOVE"
					}
				}
			}

			if scoreChanged || boardChanged {
				infoBox.SetText(fmt.Sprintf(
					"[#FFFF00::b]Score:[::-] %d\n"+
						"[#FFFF00::b]Level:[::-] %d\n"+
						"[#FFFF00::b]Lines:[::-] %d\n\n"+
						"[::b]Controls:[::-]\n"+
						"[#FFFF00]←→[white] Move\n"+
						"[#FFFF00]↑[white] Rotate\n"+
						"[#FFFF00]Z[white] Rotate CCW\n"+
						"[#FFFF00]↓[white] Soft Drop\n"+
						"[#FFFF00]Space[white] Hard Drop\n"+
						"[#FFFF00]G[white] Ghost[white] Mode\n"+
						"[#FFFF00]P[white] Pause\n"+
						"[#FFFF00]R[white] Restart\n"+
						"[#FFFF00]Q[white] Quit",
					g.Score, g.Level, g.LinesCleared))
				lastScore = g.Score
				lastLevel = g.Level
				lastLines = g.LinesCleared
			}

			if autoPlayer.IsEnabled() {
				speedText.SetText(fmt.Sprintf("SPEED:   %d", autoPlayer.GetSpeedLevel()))
			} else {
				speedText.SetText("SPEED:   -")
			}

			mu.Lock()
			ge := ghostEnabled
			mu.Unlock()

			if !autoPlayer.IsEnabled() {
				if ge {
					ghostText.SetText("GHOST: ON ")
				} else {
					ghostText.SetText("GHOST: OFF")
				}
			}

			if pieceMoved {
				lastPiecePos = currentPiecePos
			}

			if boardChanged || scoreChanged || pieceMoved {
				app.Draw()
			}

			time.Sleep(16 * time.Millisecond)
		}
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
