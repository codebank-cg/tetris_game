package model

import "fmt"

// AutoPlayer represents the autonomous Tetris player AI.
type AutoPlayer struct {
	enabled        bool             // is auto-play active
	speedLevel     int              // 1-5 speed setting
	targetDecision *MoveDecision    // current target move
	moveIndex      int              // current step in move execution
	weights        map[string]float64 // heuristic weights for board evaluation
}

// MoveDecision represents the AI's decision for piece placement.
type MoveDecision struct {
	rotations int     // target rotation (0-3)
	targetX   int     // target X position (0-9)
	softDrops int     // number of soft drops needed
	score     float64 // evaluation score
}

// GetRotations returns target rotation.
func (m *MoveDecision) GetRotations() int {
	return m.rotations
}

// GetTargetX returns target X position.
func (m *MoveDecision) GetTargetX() int {
	return m.targetX
}

// GetSoftDrops returns number of soft drops.
func (m *MoveDecision) GetSoftDrops() int {
	return m.softDrops
}

// GetScore returns evaluation score.
func (m *MoveDecision) GetScore() float64 {
	return m.score
}

// DefaultWeights returns the default heuristic weights for board evaluation.
func DefaultWeights() map[string]float64 {
	return map[string]float64{
		"aggregateHeight": -0.10,
		"holes":           -0.15,
		"bumpiness":       -0.05,
		"wells":           -0.05,
	}
}

// NewAutoPlayer creates a new AutoPlayer with default settings.
func NewAutoPlayer() *AutoPlayer {
	return &AutoPlayer{
		enabled:    false,
		speedLevel: 1,
		weights:    DefaultWeights(),
	}
}

// GetWeights returns a copy of the AutoPlayer's heuristic weights.
func (ap *AutoPlayer) GetWeights() map[string]float64 {
	copy := make(map[string]float64)
	for k, v := range ap.weights {
		copy[k] = v
	}
	return copy
}

// SetWeights updates the AutoPlayer's heuristic weights.
func (ap *AutoPlayer) SetWeights(weights map[string]float64) {
	for k, v := range weights {
		ap.weights[k] = v
	}
}

// SetSpeedLevel sets the speed level with clamping to 1-5.
func (ap *AutoPlayer) SetSpeedLevel(level int) {
	if level < 1 {
		level = 1
	}
	if level > 5 {
		level = 5
	}
	ap.speedLevel = level
}

// GetSpeedLevel returns the current speed level.
func (ap *AutoPlayer) GetSpeedLevel() int {
	return ap.speedLevel
}

// Toggle flips the enabled state.
func (ap *AutoPlayer) Toggle() {
	ap.enabled = !ap.enabled
}

// IsEnabled returns true if auto-play is active.
func (ap *AutoPlayer) IsEnabled() bool {
	return ap.enabled
}

// CycleSpeed cycles speed level: 1→2→3→4→5→1.
func (ap *AutoPlayer) CycleSpeed() {
	ap.speedLevel++
	if ap.speedLevel > 5 {
		ap.speedLevel = 1
	}
}

// IsValid validates the MoveDecision fields.
func (m *MoveDecision) IsValid() bool {
	if m.rotations < 0 || m.rotations > 3 {
		return false
	}
	if m.targetX < 0 || m.targetX > 9 {
		return false
	}
	return true
}

// Reset clears the decision state.
func (m *MoveDecision) Reset() {
	m.rotations = 0
	m.targetX = 0
	m.softDrops = 0
	m.score = 0.0
}

// String returns a readable format for debugging.
func (m *MoveDecision) String() string {
	return formatMoveDecision(*m)
}

// formatMoveDecision formats a MoveDecision as a string.
func formatMoveDecision(m MoveDecision) string {
	return fmt.Sprintf("MoveDecision{rot:%d, x:%d, drops:%d, score:%.2f}",
		m.rotations, m.targetX, m.softDrops, m.score)
}

// CalculateSoftDrops calculates drops needed to land at targetX.
func CalculateSoftDrops(board *Board, piece *Tetromino, targetX int) int {
	if piece == nil {
		return 0
	}

	// Start from current Y and count down until collision
	testY := piece.Y
	for testY > 0 {
		testY--
		if !isValidPositionForDrop(board, piece, targetX, testY) {
			break
		}
	}
	return piece.Y - testY - 1
}

// isValidPositionForDrop checks if piece can be at position.
func isValidPositionForDrop(board *Board, piece *Tetromino, x, y int) bool {
	matrix := piece.GetMatrix()
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if matrix[row][col] != 0 {
				boardX := x + col
				boardY := y - row
				if !board.IsWithinBounds(boardX, boardY) {
					return false
				}
				if boardY >= 0 && !board.IsEmpty(boardX, boardY) {
					return false
				}
			}
		}
	}
	return true
}

// getColHeight returns height of pieces in column.
func getColHeight(board *Board, col int) int {
	for y := 19; y >= 0; y-- {
		if board.Get(col, y) != 0 {
			return y + 1
		}
	}
	return 0
}

// GetAggregateHeight returns sum of all column heights.
func GetAggregateHeight(board *Board) int {
	total := 0
	for x := 0; x < 10; x++ {
		total += getColHeight(board, x)
	}
	return total
}

// countCompleteLines returns number of full lines.
func countCompleteLines(board *Board) int {
	count := 0
	for y := 0; y < 20; y++ {
		if board.IsLineFull(y) {
			count++
		}
	}
	return count
}

// evaluateLineClears returns exponential bonus for multi-line clears.
// LINE CLEARS ARE THE HIGHEST PRIORITY - weights heavily favor clearing lines.
// Ratios: 1-line=1×, 2-line=10×, 3-line=50×, 4-line=150×
// This makes line clears the dominant scoring factor above all else.
func evaluateLineClears(lines int) float64 {
	switch lines {
	case 1:
		return 1.00 // Base value (increased from 0.40)
	case 2:
		return 10.00 // 10× single-line (strong 2-line incentive)
	case 3:
		return 50.00 // 50× single-line (very strong 3-line incentive)
	case 4:
		return 150.00 // 150× single-line (Tetris is ABSOLUTE highest priority!)
	default:
		return 0.0
	}
}

// CountHoles returns count of empty cells with blocks above.
func CountHoles(board *Board) int {
	holes := 0
	for x := 0; x < 10; x++ {
		foundBlock := false
		for y := 19; y >= 0; y-- {
			if board.Get(x, y) != 0 {
				foundBlock = true
			} else if foundBlock {
				holes++
			}
		}
	}
	return holes
}

// CalculateBumpiness returns sum of height differences between adjacent columns.
func CalculateBumpiness(board *Board) int {
	bumpiness := 0
	prevHeight := getColHeight(board, 0)
	for x := 1; x < 10; x++ {
		height := getColHeight(board, x)
		diff := height - prevHeight
		if diff < 0 {
			diff = -diff
		}
		bumpiness += diff
		prevHeight = height
	}
	return bumpiness
}

// countWells returns count of columns with depth 2+ empty spaces.
func countWells(board *Board) int {
	wells := 0
	for x := 0; x < 10; x++ {
		// Check if this column forms a well (empty between filled columns)
		if x > 0 && x < 9 {
			// Count consecutive empty cells from bottom
			depth := 0
			for y := 0; y < 20; y++ {
				if board.Get(x, y) == 0 {
					depth++
				} else {
					break
				}
			}
			// Check if surrounded by blocks
			if depth >= 2 {
				// Check if there are blocks on both sides at some height
				hasLeftWall := false
				hasRightWall := false
				for y := 0; y < depth && y < 20; y++ {
					if board.Get(x-1, y) != 0 {
						hasLeftWall = true
					}
					if board.Get(x+1, y) != 0 {
						hasRightWall = true
					}
				}
				if hasLeftWall && hasRightWall {
					wells++
				}
			}
		}
	}
	return wells
}

// evalAggregateHeight returns aggregate height score (lower is better).
func evalAggregateHeight(board *Board) int {
	return GetAggregateHeight(board)
}

// evalHoles returns hole count (lower is better).
func evalHoles(board *Board) int {
	return CountHoles(board)
}

// evalBumpiness returns bumpiness score (lower is better).
func evalBumpiness(board *Board) int {
	return CalculateBumpiness(board)
}

// evalWells returns well count (lower is better).
func evalWells(board *Board) int {
	return countWells(board)
}

// evaluateBoard calculates weighted score for a board position.
func evaluateBoard(board *Board, piece *Tetromino, x, rotations int, weights map[string]float64) float64 {
	// Evaluate current board state
	aggHeight := float64(evalAggregateHeight(board))
	holes := float64(evalHoles(board))
	bumpiness := float64(evalBumpiness(board))
	wells := float64(evalWells(board))
	lines := countCompleteLines(board)

	// Apply weights with exponential line-clear bonus
	score := weights["aggregateHeight"]*aggHeight +
		weights["holes"]*holes +
		weights["bumpiness"]*bumpiness +
		weights["wells"]*wells +
		evaluateLineClears(lines)

	return score
}

// enumerateMoves generates all valid moves for current piece.
func enumerateMoves(gameState *GameState, piece *Tetromino) []MoveDecision {
	moves := []MoveDecision{}

	for rot := 0; rot < 4; rot++ {
		for x := 0; x < 10; x++ {
			if isValidMove(gameState, piece, x, rot) {
				drops := calculateDropsForMove(gameState, piece, x, rot)
				moves = append(moves, MoveDecision{
					rotations: rot,
					targetX:   x,
					softDrops: drops,
					score:     0,
				})
			}
		}
	}

	return moves
}

// FindBestMove finds the optimal move for current game state.
func FindBestMove(gameState *GameState, weights map[string]float64) *MoveDecision {
	if gameState.CurrentPiece == nil {
		return nil
	}

	moves := enumerateMoves(gameState, gameState.CurrentPiece)
	if len(moves) == 0 {
		return nil
	}

	var bestMove *MoveDecision
	bestScore := -999999.0

	for i := range moves {
		score := simulateAndEvaluate(gameState, &moves[i], weights)
		moves[i].score = score

		if score > bestScore || (score == bestScore && bestMove != nil && shouldPreferMove(moves[i], *bestMove)) {
			bestScore = score
			bestMove = &moves[i]
		}
	}

	return bestMove
}

// FindBestMoveWithNext finds the optimal move using two-piece lookahead.
// Evaluates current piece + next piece sequences to find combo setups.
// Falls back to FindBestMove() if next piece is not available.
// Returns nil if game state is invalid (nil pieces, game over, or no valid moves).
func FindBestMoveWithNext(gameState *GameState, weights map[string]float64) *MoveDecision {
	if gameState == nil || gameState.CurrentPiece == nil {
		return nil
	}

	// If no next piece available, use single-piece evaluation
	if gameState.NextPiece == nil {
		return FindBestMove(gameState, weights)
	}

	moves := enumerateMoves(gameState, gameState.CurrentPiece)
	if len(moves) == 0 {
		return nil
	}

	var bestMove *MoveDecision
	bestScore := -999999.0

	for i := range moves {
		// Use two-piece lookahead evaluation
		score := EvaluateTwoPieceSequence(gameState, &moves[i], gameState.NextPiece, weights)
		moves[i].score = score

		if score > bestScore || (score == bestScore && bestMove != nil && shouldPreferMove(moves[i], *bestMove)) {
			bestScore = score
			bestMove = &moves[i]
		}
	}

	return bestMove
}

// simulateAndEvaluate simulates move and returns score.
func simulateAndEvaluate(gameState *GameState, move *MoveDecision, weights map[string]float64) float64 {
	testBoard := cloneBoard(gameState.Board)
	testPiece := clonePiece(gameState.CurrentPiece)

	for i := 0; i < move.rotations; i++ {
		testPiece.RotateClockwise()
	}
	testPiece.X = move.targetX
	testPiece.Y -= move.softDrops

	placePieceOnBoard(testBoard, testPiece)

	return evaluateBoard(testBoard, testPiece, move.targetX, move.rotations, weights)
}

// EvaluateTwoPieceSequence evaluates a move sequence: current piece + next piece.
// Returns combined score of both pieces' placements, including multi-line bonuses.
// This enables the AI to plan setups where current piece enables better next piece placement.
func EvaluateTwoPieceSequence(gameState *GameState, currentMove *MoveDecision, nextPiece *Tetromino, weights map[string]float64) float64 {
	if gameState == nil || currentMove == nil || nextPiece == nil {
		return -999999.0
	}

	testBoard := cloneBoard(gameState.Board)
	testPiece := clonePiece(gameState.CurrentPiece)

	for i := 0; i < currentMove.rotations; i++ {
		testPiece.RotateClockwise()
	}
	testPiece.X = currentMove.targetX
	testPiece.Y -= currentMove.softDrops

	placePieceOnBoard(testBoard, testPiece)

	linesClearedByCurrent := countCompleteLines(testBoard)

	// Clear complete lines from testBoard so the next piece is placed on the correct board state.
	for y := 19; y >= 0; y-- {
		if testBoard.IsLineFull(y) {
			testBoard.ClearLine(y)
		}
	}

	bestNextScore := -999999.0
	bestComboLines := 0
	nextMoves := enumerateMovesForBoard(testBoard, nextPiece)

	if len(nextMoves) == 0 {
		return -999999.0
	}

	for i := range nextMoves {
		nextBoard := cloneBoard(testBoard)
		nextTestPiece := clonePiece(nextPiece)

		for j := 0; j < nextMoves[i].rotations; j++ {
			nextTestPiece.RotateClockwise()
		}
		nextTestPiece.X = nextMoves[i].targetX
		nextTestPiece.Y = 18

		for nextTestPiece.Y > 0 {
			nextTestPiece.Y--
			if !isValidPositionForBoard(nextBoard, nextTestPiece, nextTestPiece.X, nextTestPiece.Y) {
				nextTestPiece.Y++
				break
			}
		}

		placePieceOnBoard(nextBoard, nextTestPiece)

		comboLines := countCompleteLines(nextBoard)

		nextScore := evaluateBoard(nextBoard, nextTestPiece, nextMoves[i].targetX, nextMoves[i].rotations, weights)

		if comboLines >= 2 {
			if comboLines == 4 {
				nextScore += evaluateLineClears(4) * 3.0
			} else if comboLines == 3 {
				nextScore += evaluateLineClears(3) * 2.5
			} else {
				nextScore += evaluateLineClears(comboLines) * 2.0
			}
		} else if comboLines == 1 {
			nextScore += evaluateLineClears(1) * 1.2
		}

		if comboLines > bestComboLines {
			nextScore += float64(comboLines) * 5.0
		}

		if nextScore > bestNextScore {
			bestNextScore = nextScore
			bestComboLines = comboLines
		}
	}

	currentScore := evaluateBoard(testBoard, testPiece, currentMove.targetX, currentMove.rotations, weights) + evaluateLineClears(linesClearedByCurrent)

	totalScore := currentScore*0.4 + bestNextScore*0.6

	if bestComboLines >= 2 {
		totalScore += float64(bestComboLines) * 10.0
	}

	return totalScore
}

// enumerateMovesForBoard generates all valid moves for a piece on a given board.
// Helper for two-piece lookahead evaluation.
func enumerateMovesForBoard(board *Board, piece *Tetromino) []MoveDecision {
	moves := []MoveDecision{}

	for rot := 0; rot < 4; rot++ {
		for x := 0; x < 10; x++ {
			// Create test piece for validation
			testPiece := clonePiece(piece)
			for i := 0; i < rot; i++ {
				testPiece.RotateClockwise()
			}
			testPiece.X = x
			testPiece.Y = 18

			// Check if this position is valid
			if isValidPositionForBoard(board, testPiece, x, 18) {
				drops := 0
				// Count how many drops are possible
				for y := 18; y > 0; y-- {
					if !isValidPositionForBoard(board, testPiece, x, y-1) {
						break
					}
					drops++
				}

				moves = append(moves, MoveDecision{
					rotations: rot,
					targetX:   x,
					softDrops: drops,
					score:     0,
				})
			}
		}
	}

	return moves
}

// isValidPositionForBoard checks if piece can be placed on board at position.
func isValidPositionForBoard(board *Board, piece *Tetromino, x, y int) bool {
	matrix := piece.GetMatrix()
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if matrix[row][col] != 0 {
				boardX := x + col
				boardY := y - row
				if !board.IsWithinBounds(boardX, boardY) {
					return false
				}
				if boardY >= 0 && !board.IsEmpty(boardX, boardY) {
					return false
				}
			}
		}
	}
	return true
}

// cloneBoard creates a copy of the board.
func cloneBoard(board *Board) *Board {
	clone := NewBoard()
	for y := 0; y < 20; y++ {
		for x := 0; x < 10; x++ {
			clone.cells[y][x] = board.cells[y][x]
		}
	}
	return clone
}

// placePieceOnBoard locks piece to board.
func placePieceOnBoard(board *Board, piece *Tetromino) {
	matrix := piece.GetMatrix()
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if matrix[row][col] != 0 {
				boardX := piece.X + col
				boardY := piece.Y - row
				if boardY >= 0 && board.IsWithinBounds(boardX, boardY) {
					board.Set(boardX, boardY, piece.Color)
				}
			}
		}
	}
}

// shouldPreferMove returns true if move1 is preferable to move2 (tie-breaker).
func shouldPreferMove(move1, move2 MoveDecision) bool {
	if move1.targetX < move2.targetX {
		return true
	}
	if move1.rotations < move2.rotations {
		return true
	}
	return false
}

// executeRotations rotates piece to target rotation.
func executeRotations(gameState *GameState, targetRotations int) {
	current := gameState.CurrentPiece.Rotation
	rotationsNeeded := (targetRotations - current + 4) % 4

	for i := 0; i < rotationsNeeded; i++ {
		gameState.RotatePiece()
	}
}

// executeHorizontalMove moves piece horizontally to targetX.
func executeHorizontalMove(gameState *GameState, targetX int) {
	for gameState.CurrentPiece.X < targetX {
		if !gameState.MovePiece(1, 0) {
			break
		}
	}
	for gameState.CurrentPiece.X > targetX {
		if !gameState.MovePiece(-1, 0) {
			break
		}
	}
}

// executeDrop performs soft drops.
func executeDrop(gameState *GameState, drops int) {
	for i := 0; i < drops; i++ {
		if !gameState.MovePiece(0, -1) {
			break
		}
	}
}

// GetDelayForSpeed calculates delay based on speed level.
func GetDelayForSpeed(baseDelay int, speedLevel int) int {
	multipliers := map[int]float64{
		1: 1.0,
		2: 0.5,
		3: 0.2,
		4: 0.1,
		5: 0.0,
	}
	mult, ok := multipliers[speedLevel]
	if !ok {
		mult = 1.0
	}
	return int(float64(baseDelay) * mult)
}

// shouldHoldPiece determines if holding current piece is beneficial.
func shouldHoldPiece(gameState *GameState, decision *MoveDecision) bool {
	if !gameState.CanHold {
		return false
	}
	if gameState.HoldPiece == nil {
		return true
	}
	return false
}

// executeHold performs hold action.
func executeHold(gameState *GameState) {
	gameState.HoldCurrentPiece()
}

// IsInDropPhase checks if the current move is in the drop phase.
func IsInDropPhase(gameState *GameState) bool {
	return gameState.autoMoveStep == 2
}

// ExecuteMove executes a complete move decision step by step.
// Returns true if move is complete, false if more steps needed.
func ExecuteMove(gameState *GameState, decision *MoveDecision) bool {
	if decision == nil || gameState.CurrentPiece == nil {
		return true
	}

	if gameState.autoMoveStep == 0 {
		if decision.rotations > 0 {
			gameState.RotatePiece()
			decision.rotations--
			return false
		}
		gameState.autoMoveStep = 1
	}

	if gameState.autoMoveStep == 1 {
		if gameState.CurrentPiece.X != decision.targetX {
			if gameState.CurrentPiece.X < decision.targetX {
				gameState.MovePiece(1, 0)
			} else {
				gameState.MovePiece(-1, 0)
			}
			return false
		}
		gameState.autoMoveStep = 2
	}

	if gameState.autoMoveStep == 2 {
		if decision.softDrops > 0 {
			gameState.MovePiece(0, -1)
			decision.softDrops--
			return false
		}
		gameState.lockPiece()
		gameState.autoMoveStep = 0
		return true
	}

	return true
}

// isValidMove checks if a move is valid.
func isValidMove(gameState *GameState, piece *Tetromino, x, rot int) bool {
	testPiece := clonePiece(piece)

	for i := 0; i < rot; i++ {
		testPiece.RotateClockwise()
	}
	testPiece.X = x

	matrix := testPiece.GetMatrix()
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if matrix[row][col] != 0 {
				boardX := x + col
				boardY := testPiece.Y - row
				if !gameState.Board.IsWithinBounds(boardX, boardY) {
					return false
				}
				if boardY >= 0 && !gameState.Board.IsEmpty(boardX, boardY) {
					return false
				}
			}
		}
	}
	return true
}

// calculateDropsForMove calculates soft drops needed.
func calculateDropsForMove(gameState *GameState, piece *Tetromino, x, rot int) int {
	testPiece := clonePiece(piece)
	for i := 0; i < rot; i++ {
		testPiece.RotateClockwise()
	}
	testPiece.X = x

	drops := 0
	for gameState.isValidPosition(testPiece, testPiece.X, testPiece.Y-1) {
		testPiece.Y--
		drops++
	}
	return drops
}

// clonePiece creates a copy of a piece.
func clonePiece(piece *Tetromino) *Tetromino {
	return &Tetromino{
		Type:     piece.Type,
		Color:    piece.Color,
		X:        piece.X,
		Y:        piece.Y,
		Rotation: piece.Rotation,
		Matrix:   piece.Matrix,
	}
}
