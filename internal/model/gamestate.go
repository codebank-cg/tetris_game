package model

// GameState represents the current state of the game.
type GameState struct {
	Board          *Board
	CurrentPiece   *Tetromino
	NextPiece      *Tetromino
	HoldPiece      *Tetromino
	Randomizer     *Randomizer
	Score          int
	Level          int
	LinesCleared   int
	TetrisCount    int
	PieceCount     int
	CanHold        bool
	GameOver       bool
	Paused         bool
	DropInterval   int
	ClearedLines   []int
	ClearAnimFrame int
	ClearAnimIndex int
	autoMoveStep   int
}

// NewGameState creates a new game state.
func NewGameState() *GameState {
	gs := &GameState{
		Board:        NewBoard(),
		Randomizer:   NewRandomizer(),
		Level:        1,
		CanHold:      true,
		DropInterval: 1500,
	}
	gs.NextPiece = NewTetromino(gs.Randomizer.NextPiece())
	gs.spawnPiece()
	return gs
}

// spawnPiece spawns the next piece.
func (gs *GameState) spawnPiece() {
	gs.CurrentPiece = gs.NextPiece
	gs.NextPiece = NewTetromino(gs.Randomizer.NextPiece())
	gs.CurrentPiece.X = 3
	gs.CurrentPiece.Y = 18
	gs.CanHold = true
	gs.PieceCount++
}

// MovePiece moves the current piece.
func (gs *GameState) MovePiece(dx, dy int) bool {
	if gs.CurrentPiece == nil || gs.Paused || gs.GameOver {
		return false
	}
	newX := gs.CurrentPiece.X + dx
	newY := gs.CurrentPiece.Y + dy
	if gs.isValidPosition(gs.CurrentPiece, newX, newY) {
		gs.CurrentPiece.Move(dx, dy)
		return true
	}
	return false
}

// RotatePiece rotates the current piece clockwise.
func (gs *GameState) RotatePiece() bool {
	if gs.CurrentPiece == nil || gs.Paused || gs.GameOver {
		return false
	}
	gs.CurrentPiece.RotateClockwise()
	if !gs.isValidPosition(gs.CurrentPiece, gs.CurrentPiece.X, gs.CurrentPiece.Y) {
		gs.CurrentPiece.RotateCounterClockwise()
		return false
	}
	return true
}

// RotatePieceCounter rotates the current piece counter-clockwise.
func (gs *GameState) RotatePieceCounter() bool {
	if gs.CurrentPiece == nil || gs.Paused || gs.GameOver {
		return false
	}
	gs.CurrentPiece.RotateCounterClockwise()
	if !gs.isValidPosition(gs.CurrentPiece, gs.CurrentPiece.X, gs.CurrentPiece.Y) {
		gs.CurrentPiece.RotateClockwise()
		return false
	}
	return true
}

// HoldCurrentPiece swaps the current piece with the hold piece.
func (gs *GameState) HoldCurrentPiece() bool {
	if !gs.CanHold || gs.CurrentPiece == nil || gs.Paused || gs.GameOver {
		return false
	}
	if gs.HoldPiece == nil {
		gs.HoldPiece = gs.CurrentPiece
		gs.spawnPiece()
	} else {
		gs.CurrentPiece, gs.HoldPiece = gs.HoldPiece, gs.CurrentPiece
		gs.CurrentPiece.X = 3
		gs.CurrentPiece.Y = 18
		gs.CurrentPiece.Rotation = 0
		gs.CurrentPiece.updateMatrix()
	}
	gs.CanHold = false
	return true
}

// DropPiece drops the piece instantly.
func (gs *GameState) DropPiece() int {
	if gs.CurrentPiece == nil || gs.Paused || gs.GameOver {
		return 0
	}
	dropped := 0
	for gs.MovePiece(0, -1) {
		dropped++
	}
	gs.lockPiece()
	return dropped
}

// SoftDrop moves the piece down one cell.
func (gs *GameState) SoftDrop() bool {
	if gs.MovePiece(0, -1) {
		return true
	}
	gs.lockPiece()
	return false
}

// GetGhostY returns the Y position where the piece would land if dropped.
func (gs *GameState) GetGhostY() int {
	if gs == nil || gs.CurrentPiece == nil || gs.GameOver || gs.Paused {
		return -1
	}

	ghostY := gs.CurrentPiece.Y
	maxIterations := gs.CurrentPiece.Y + 1

	for ghostY > 0 && maxIterations > 0 {
		testY := ghostY - 1
		if testY < 0 {
			break
		}
		if !gs.isValidPosition(gs.CurrentPiece, gs.CurrentPiece.X, testY) {
			break
		}
		ghostY = testY
		maxIterations--
	}
	return ghostY
}

// isValidPosition checks if piece position is valid.
func (gs *GameState) isValidPosition(piece *Tetromino, x, y int) bool {
	matrix := piece.GetMatrix()
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if matrix[row][col] != 0 {
				boardX := x + col
				boardY := y - row
				if !gs.Board.IsWithinBounds(boardX, boardY) {
					return false
				}
				if boardY >= 0 && !gs.Board.IsEmpty(boardX, boardY) {
					return false
				}
			}
		}
	}
	return true
}

// lockPiece locks the current piece to the board.
func (gs *GameState) lockPiece() {
	if gs.CurrentPiece == nil {
		return
	}
	matrix := gs.CurrentPiece.GetMatrix()
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if matrix[row][col] != 0 {
				boardX := gs.CurrentPiece.X + col
				boardY := gs.CurrentPiece.Y - row
				if boardY >= 0 && gs.Board.IsWithinBounds(boardX, boardY) {
					gs.Board.Set(boardX, boardY, gs.CurrentPiece.Color)
				}
			}
		}
	}
	gs.clearLines()
	gs.spawnPiece()
	if !gs.isValidPosition(gs.CurrentPiece, gs.CurrentPiece.X, gs.CurrentPiece.Y) {
		gs.GameOver = true
	}
}

// clearLines clears full lines and updates score.
func (gs *GameState) clearLines() {
	fullLines := []int{}
	for y := 0; y < 20; y++ {
		if gs.Board.IsLineFull(y) {
			fullLines = append(fullLines, y)
		}
	}
	if len(fullLines) > 0 {
		// Sort highest to lowest for animation order (top to bottom)
		for i := 0; i < len(fullLines)/2; i++ {
			j := len(fullLines) - 1 - i
			fullLines[i], fullLines[j] = fullLines[j], fullLines[i]
		}
		gs.ClearedLines = fullLines
		gs.ClearAnimFrame = 1
		gs.ClearAnimIndex = 0
	}
}

// UpdateClearAnimation updates the line clear animation.
// Flashes lines one by one, then clears all at once.
// Returns true if animation is still playing, false if completed.
// When lines are actually cleared (animation completes), returns completed=true.
func (gs *GameState) UpdateClearAnimation() (completed bool) {
	if gs.ClearAnimFrame <= 0 || len(gs.ClearedLines) == 0 {
		return false
	}

	gs.ClearAnimFrame++

	// Each line flashes for 12 frames (600ms) - slower, more dramatic
	if gs.ClearAnimFrame > 12 {
		gs.ClearAnimIndex++
		gs.ClearAnimFrame = 1

		// All lines flashed? Now actually clear them all at once
		if gs.ClearAnimIndex >= len(gs.ClearedLines) {
			if len(gs.ClearedLines) == 4 {
				gs.TetrisCount++
			}
			for _, line := range gs.ClearedLines {
				gs.Board.ClearLine(line)
			}
			gs.LinesCleared += len(gs.ClearedLines)
			gs.UpdateScore(len(gs.ClearedLines))
			gs.Level = gs.LinesCleared / 10 // Level up every 10 lines (official Game Boy)
			gs.ClearedLines = []int{}
			gs.ClearAnimIndex = 0
			gs.ClearAnimFrame = 0
			return true // Lines were just cleared
		}
	}
	return false // Animation still in progress
}

// GetCurrentClearedLine returns the line currently being animated.
func (gs *GameState) GetCurrentClearedLine() int {
	if gs.ClearAnimIndex >= 0 && gs.ClearAnimIndex < len(gs.ClearedLines) {
		return gs.ClearedLines[gs.ClearAnimIndex]
	}
	return -1
}

// IsClearAnimating returns true if line clear animation is playing.
func (gs *GameState) IsClearAnimating() bool {
	return gs.ClearAnimFrame > 0
}

// UpdateScore updates the score based on lines cleared.
// Original Nintendo scoring system (NES, Game Boy, SNES):
// 1 line:  40 × (level + 1)
// 2 lines: 100 × (level + 1)
// 3 lines: 300 × (level + 1)
// 4 lines: 1200 × (level + 1)
func (gs *GameState) UpdateScore(lines int) {
	baseScores := map[int]int{
		1: 40,
		2: 100,
		3: 300,
		4: 1200,
	}
	if score, ok := baseScores[lines]; ok {
		gs.Score += score * (gs.Level + 1)
	}
}

// IsPause toggles the pause state.
func (gs *GameState) Pause() {
	if !gs.GameOver {
		gs.Paused = !gs.Paused
	}
}

// Reset resets the game state.
func (gs *GameState) Reset() {
	gs.Board = NewBoard()
	gs.Randomizer = NewRandomizer()
	gs.Score = 0
	gs.Level = 1
	gs.LinesCleared = 0
	gs.TetrisCount = 0
	gs.PieceCount = 0
	gs.GameOver = false
	gs.Paused = false
	gs.CanHold = true
	gs.HoldPiece = nil
	gs.ClearedLines = []int{}
	gs.ClearAnimFrame = 0
	gs.ClearAnimIndex = 0
	gs.autoMoveStep = 0
	gs.NextPiece = NewTetromino(gs.Randomizer.NextPiece())
	gs.spawnPiece()
}

// IncreaseLevel increases the game level (faster drop).
func (gs *GameState) IncreaseLevel() {
	if gs.Level < 20 {
		gs.Level++
	}
}

// DecreaseLevel decreases the game level (slower drop).
func (gs *GameState) DecreaseLevel() {
	if gs.Level > 1 {
		gs.Level--
	}
}

// GetDropInterval returns the drop interval in milliseconds based on current level.
// Level 1: 1500ms, Level 10: 600ms, Level 20: 100ms
func (gs *GameState) GetDropInterval() int {
	// Formula: interval = 1500 - (level - 1) * 100, minimum 100ms
	interval := 1500 - (gs.Level-1)*100
	if interval < 100 {
		interval = 100
	}
	return interval
}
