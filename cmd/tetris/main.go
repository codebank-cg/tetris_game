package main

import (
	"fmt"
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

		for row := 0; row < 20; row++ {
			screenY := boardY + (19 - row)
			for col := 0; col < 10; col++ {
				cellX := boardX + col*2
				color := game.Board.Get(col, row)
				if color != 0 {
					baseColor := colors[color]
					screen.SetContent(cellX, screenY, '█', nil, tcell.StyleDefault.Foreground(baseColor).Background(tcell.ColorBlack))
					screen.SetContent(cellX+1, screenY, '█', nil, tcell.StyleDefault.Foreground(baseColor).Background(tcell.ColorBlack))
				}
			}
		}

		if game.CurrentPiece != nil && !game.GameOver && ghostEnabled && !autoPlayer.IsEnabled() {
			ghostY := game.GetGhostY()
			if ghostY >= 0 && ghostY < game.CurrentPiece.Y && ghostY < 20 {
				matrix := game.CurrentPiece.GetMatrix()
				if matrix != nil {
					ghostStyle := tcell.StyleDefault.Foreground(tcell.GetColor("#90EE90")).Background(tcell.ColorBlack)
					for row := 0; row < 4; row++ {
						for col := 0; col < 4; col++ {
							if matrix[row][col] != 0 {
								pieceY := ghostY - row
								if pieceY >= 0 && pieceY < 20 {
									bx := boardX + (game.CurrentPiece.X+col)*2
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

		if game.CurrentPiece != nil && !game.GameOver {
			matrix := game.CurrentPiece.GetMatrix()
			pieceStyle := tcell.StyleDefault.Foreground(colors[game.CurrentPiece.Color]).Background(tcell.ColorBlack)
			for row := 0; row < 4; row++ {
				for col := 0; col < 4; col++ {
					if matrix[row][col] != 0 {
						bx := boardX + (game.CurrentPiece.X+col)*2
						by := boardY + (19 - (game.CurrentPiece.Y - row))
						if by >= boardY && by < boardY+20 && bx >= boardX && bx < boardX+20 {
							screen.SetContent(bx, by, '█', nil, pieceStyle)
							screen.SetContent(bx+1, by, '█', nil, pieceStyle)
						}
					}
				}
			}
		}

		if game.GameOver {
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
		if game.NextPiece == nil {
			return x, y, width, height
		}

		boxWidth := width - 2
		boxHeight := height - 2
		matrix := game.NextPiece.GetMatrix()
		pieceColor := colors[game.NextPiece.Color]
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
					ghostEnabled = false
					ghostText.SetText("GHOST: OFF")
				} else {
					autoText.SetText("AUTO-PLAY: OFF")
					speedText.SetText("SPEED:   -")
					autoParamsText.SetText("")
				}
				return nil
			case 'z', 'Z':
				if !game.Paused && !game.GameOver {
					game.RotatePieceCounter()
				}
				return nil
			case 'g', 'G':
				if !autoPlayer.IsEnabled() {
					ghostEnabled = !ghostEnabled
					if ghostEnabled {
						ghostText.SetText("GHOST: ON ")
					} else {
						ghostText.SetText("GHOST: OFF")
					}
				}
				return nil
			case ' ':
				if !game.Paused && !game.GameOver {
					dropped := game.DropPiece()
					game.UpdateScore(dropped)
				}
				return nil
			case 'p', 'P':
				game.Paused = !game.Paused
				app.Draw()
				return nil
			case 'r', 'R':
				game = model.NewGameState()
				if autoPlayer.IsEnabled() {
					autoText.SetText("AUTO-PLAY: ON ")
					speedText.SetText(fmt.Sprintf("SPEED:   %d", autoPlayer.GetSpeedLevel()))
					ghostEnabled = false
					ghostText.SetText("GHOST: OFF")
				} else {
					autoText.SetText("AUTO-PLAY: OFF")
					speedText.SetText("SPEED:   -")
					if !ghostEnabled {
						ghostText.SetText("GHOST: OFF")
					}
				}
				autoParamsText.SetText("TARGET:  X: -  R:-\nSCORE:   -. -")
				return nil
			}
		case tcell.KeyLeft:
			if !game.Paused && !game.GameOver {
				game.MovePiece(-1, 0)
			}
			return nil
		case tcell.KeyRight:
			if !game.Paused && !game.GameOver {
				game.MovePiece(1, 0)
			}
			return nil
		case tcell.KeyDown:
			if !game.Paused && !game.GameOver {
				game.SoftDrop()
				lastDrop = time.Now()
			}
			return nil
		case tcell.KeyUp:
			if !game.Paused && !game.GameOver {
				game.RotatePiece()
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

			if autoPlayer.IsEnabled() && !game.Paused && !game.GameOver && !game.IsClearAnimating() && game.CurrentPiece != nil {
				delay := model.GetDelayForSpeed(game.GetDropInterval(), autoPlayer.GetSpeedLevel())
				if model.IsInDropPhase(game) {
					delay = delay / 3
				}
				if time.Since(lastActionTime) > time.Duration(delay)*time.Millisecond {
					if currentDecision == nil {
						currentDecision = model.FindBestMoveWithNext(game)
					}
					if currentDecision != nil {
						if model.ExecuteMove(game, currentDecision) {
							currentDecision = nil
						}
						lastActionTime = time.Now()
						boardChanged = true
					} else {
						game.SoftDrop()
						lastActionTime = time.Now()
						boardChanged = true
					}
				}
			}

			if !autoPlayer.IsEnabled() && time.Since(lastDrop) > time.Duration(game.GetDropInterval())*time.Millisecond && !game.Paused && !game.GameOver && !game.IsClearAnimating() {
				game.SoftDrop()
				lastDrop = time.Now()
				boardChanged = true
			}

			if game.IsClearAnimating() {
				linesCleared := game.UpdateClearAnimation()
				if linesCleared {
					musicPlayer.PlayLineClearBeep()
					boardChanged = true
				}
			}

			if game.GameOver && !gameOverDetected {
				gameOverDetected = true
				musicPlayer.Stop()
				boardChanged = true
			}

			scoreChanged := game.Score != lastScore || game.Level != lastLevel || game.LinesCleared != lastLines
			currentPiecePos := 0
			pieceMoved := false
			if game.CurrentPiece != nil {
				currentPiecePos = game.CurrentPiece.X*100 + game.CurrentPiece.Y
				pieceMoved = currentPiecePos != lastPiecePos
			}

			if autoPlayer.IsEnabled() && !game.GameOver && game.CurrentPiece != nil {
				decision := model.FindBestMoveWithNext(game)
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
					game.Score, game.Level, game.LinesCleared))
				lastScore = game.Score
				lastLevel = game.Level
				lastLines = game.LinesCleared
			}

			if autoPlayer.IsEnabled() {
				speedText.SetText(fmt.Sprintf("SPEED:   %d", autoPlayer.GetSpeedLevel()))
			} else {
				speedText.SetText("SPEED:   -")
			}

			if !autoPlayer.IsEnabled() {
				if ghostEnabled {
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
