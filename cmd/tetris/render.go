package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/oc-garden/tetris_game/internal/model"
)

// drawBoard renders a Tetris board at the given box position.
// x, y are the outer box coordinates (as provided by tview SetDrawFunc).
// showGhost: render ghost piece (use false for showdown, true for single-player when enabled).
func drawBoard(screen tcell.Screen, x, y int, game *model.GameState, colors map[int]tcell.Color, showGhost bool) {
	boardX := x + 1
	boardY := y + 1

	// Draw placed cells
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

	// Draw ghost piece
	if showGhost && game.CurrentPiece != nil && !game.GameOver {
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

	// Draw current piece
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

	// Draw game-over overlay
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
}

// drawText is a helper to render a string at (x, y) with the given style.
func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}
