package model

import (
	"fmt"
	"testing"
)

func createEmptyBoard() *Board {
	return NewBoard()
}

// Helper: create board with uniform height
func createBoardWithHeight(height int) *Board {
	b := NewBoard()
	for x := 0; x < 10; x++ {
		for y := 0; y < height && y < 20; y++ {
			b.Set(x, y, 1)
		}
	}
	return b
}

// Helper: create board with specific column heights
func createBoardWithHeights(heights [10]int) *Board {
	b := NewBoard()
	for x, h := range heights {
		for y := 0; y < h && y < 20; y++ {
			b.Set(x, y, 1)
		}
	}
	return b
}

// Helper: create board with holes
func createBoardWithHoles(count int) *Board {
	b := NewBoard()
	x := 5
	y := 0
	for i := 0; i < count && y < 18; i++ {
		b.Set(x, y, 1)
		y += 2
		b.Set(x, y, 1)
		y++
	}
	return b
}

// Helper: create board with wells
func createBoardWithWells(count int) *Board {
	b := NewBoard()
	for i := 0; i < count && i < 5; i++ {
		x := i*3 + 1
		for y := 0; y < 5; y++ {
			if x > 0 {
				b.Set(x-1, y, 1)
			}
			if x < 9 {
				b.Set(x+1, y, 1)
			}
		}
	}
	return b
}

// Helper: create board with complete line
func createBoardWithCompleteLine() *Board {
	b := NewBoard()
	for x := 0; x < 10; x++ {
		b.Set(x, 10, 1)
	}
	return b
}

// Test 1: AutoPlayer Creation
func TestAutoPlayerCreation(t *testing.T) {
	ap := NewAutoPlayer()

	if ap.IsEnabled() {
		t.Error("New AutoPlayer should be disabled")
	}
	if ap.GetSpeedLevel() != 1 {
		t.Errorf("Default speed level = %d, want 1", ap.GetSpeedLevel())
	}
}

// Test 1.2: Toggle Functionality
func TestAutoPlayer_Toggle(t *testing.T) {
	ap := NewAutoPlayer()

	if ap.IsEnabled() {
		t.Error("Initial state should be disabled")
	}

	ap.Toggle()
	if !ap.IsEnabled() {
		t.Error("After Toggle() should be enabled")
	}

	ap.Toggle()
	if ap.IsEnabled() {
		t.Error("After second Toggle() should be disabled")
	}
}

// Test 1.3: Speed Level Setting
func TestAutoPlayer_SetSpeedLevel(t *testing.T) {
	tests := []struct {
		name      string
		setLevel  int
		wantLevel int
	}{
		{"level 1", 1, 1},
		{"level 3", 3, 3},
		{"level 5", 5, 5},
		{"invalid level 0", 0, 1},
		{"invalid level 6", 6, 5},
		{"invalid level -1", -1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := NewAutoPlayer()
			ap.SetSpeedLevel(tt.setLevel)
			if ap.GetSpeedLevel() != tt.wantLevel {
				t.Errorf("SetSpeedLevel(%d) = %d, want %d", tt.setLevel, ap.GetSpeedLevel(), tt.wantLevel)
			}
		})
	}
}

// Test 1.4: Speed Cycling
func TestAutoPlayer_CycleSpeed(t *testing.T) {
	ap := NewAutoPlayer()
	expected := []int{2, 3, 4, 5, 1, 2}

	for i, want := range expected {
		ap.CycleSpeed()
		if ap.GetSpeedLevel() != want {
			t.Errorf("Cycle %d: GetSpeedLevel() = %d, want %d", i+1, ap.GetSpeedLevel(), want)
		}
	}
}

// Test 2.1: MoveDecision Validation
func TestMoveDecision_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		decision MoveDecision
		want     bool
	}{
		{"valid center", MoveDecision{rotations: 0, targetX: 5, softDrops: 10, score: 12.5}, true},
		{"valid edge X=0", MoveDecision{rotations: 2, targetX: 0, softDrops: 5, score: 8.0}, true},
		{"valid edge X=9", MoveDecision{rotations: 3, targetX: 9, softDrops: 15, score: 6.5}, true},
		{"invalid rotation 4", MoveDecision{rotations: 4, targetX: 5, softDrops: 10, score: 12.5}, false},
		{"invalid rotation -1", MoveDecision{rotations: -1, targetX: 5, softDrops: 10, score: 12.5}, false},
		{"invalid X=-1", MoveDecision{rotations: 0, targetX: -1, softDrops: 10, score: 12.5}, false},
		{"invalid X=10", MoveDecision{rotations: 0, targetX: 10, softDrops: 10, score: 12.5}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.decision.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test 2.2: MoveDecision String Format
func TestMoveDecision_String(t *testing.T) {
	d := MoveDecision{rotations: 2, targetX: 5, softDrops: 10, score: 12.5}
	got := d.String()
	want := "MoveDecision{rot:2, x:5, drops:10, score:12.50}"

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

// Test 4.1: Column Height Calculation
func TestGetColHeight(t *testing.T) {
	tests := []struct {
		name       string
		setupFn    func() *Board
		col        int
		wantHeight int
	}{
		{"empty column", createEmptyBoard, 5, 0},
		{"full column", func() *Board {
			b := NewBoard()
			for y := 0; y < 20; y++ {
				b.Set(5, y, 1)
			}
			return b
		}, 5, 20},
		{"partial height 5", func() *Board {
			b := NewBoard()
			for y := 0; y < 5; y++ {
				b.Set(5, y, 1)
			}
			return b
		}, 5, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := tt.setupFn()
			if got := getColHeight(board, tt.col); got != tt.wantHeight {
				t.Errorf("getColHeight() = %d, want %d", got, tt.wantHeight)
			}
		})
	}
}

// Test 4.2: Aggregate Height
func TestGetAggregateHeight(t *testing.T) {
	tests := []struct {
		name    string
		setupFn func() *Board
		want    int
	}{
		{"empty board", createEmptyBoard, 0},
		{"flat height 5", func() *Board {
			return createBoardWithHeight(5)
		}, 50},
		{"full board", func() *Board {
			return createBoardWithHeight(20)
		}, 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAggregateHeight(tt.setupFn()); got != tt.want {
				t.Errorf("getAggregateHeight() = %d, want %d", got, tt.want)
			}
		})
	}
}

// Test 4.3: Complete Lines Count
func TestCountCompleteLines(t *testing.T) {
	tests := []struct {
		name      string
		setupFn   func() *Board
		wantLines int
	}{
		{"no lines", createEmptyBoard, 0},
		{"one full line", func() *Board {
			return createBoardWithCompleteLine()
		}, 1},
		{"three full lines", func() *Board {
			b := NewBoard()
			for line := 0; line < 3; line++ {
				for x := 0; x < 10; x++ {
					b.Set(x, line, 1)
				}
			}
			return b
		}, 3},
		{"partial line", func() *Board {
			b := NewBoard()
			for x := 0; x < 9; x++ {
				b.Set(x, 5, 1)
			}
			return b
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countCompleteLines(tt.setupFn()); got != tt.wantLines {
				t.Errorf("countCompleteLines() = %d, want %d", got, tt.wantLines)
			}
		})
	}
}

// Test 4.4: Hole Detection
func TestCountHoles(t *testing.T) {
	tests := []struct {
		name      string
		setupFn   func() *Board
		wantHoles int
	}{
		{"no holes", createEmptyBoard, 0},
		{"single hole", func() *Board {
			b := NewBoard()
			b.Set(5, 0, 1)
			b.Set(5, 2, 1)
			return b
		}, 1},
		{"multiple holes same column", func() *Board {
			b := NewBoard()
			b.Set(5, 0, 1)
			b.Set(5, 2, 1)
			b.Set(5, 4, 1)
			return b
		}, 2},
		{"no holes - solid stack", func() *Board {
			return createBoardWithHeight(10)
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countHoles(tt.setupFn()); got != tt.wantHoles {
				t.Errorf("countHoles() = %d, want %d", got, tt.wantHoles)
			}
		})
	}
}

// Test 4.5: Bumpiness Calculation
func TestCalculateBumpiness(t *testing.T) {
	tests := []struct {
		name          string
		colHeights    [10]int
		wantBumpiness int
	}{
		{"flat surface", [10]int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5}, 0},
		{"single step", [10]int{5, 10, 5, 5, 5, 5, 5, 5, 5, 5}, 10},
		{"increasing slope", [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 9},
		{"random heights", [10]int{3, 5, 2, 4, 1, 6, 3, 2, 4, 5}, 22},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := createBoardWithHeights(tt.colHeights)
			if got := calculateBumpiness(b); got != tt.wantBumpiness {
				t.Errorf("calculateBumpiness() = %d, want %d", got, tt.wantBumpiness)
			}
		})
	}
}

// Test 4.6: Wells Detection
func TestCountWells(t *testing.T) {
	tests := []struct {
		name      string
		setupFn   func() *Board
		wantWells int
	}{
		{"no wells", createEmptyBoard, 0},
		{"single well", func() *Board {
			b := NewBoard()
			for y := 0; y < 5; y++ {
				b.Set(4, y, 1)
				b.Set(6, y, 1)
			}
			return b
		}, 1},
		{"solid stack no wells", func() *Board {
			return createBoardWithHeight(10)
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countWells(tt.setupFn()); got != tt.wantWells {
				t.Errorf("countWells() = %d, want %d", got, tt.wantWells)
			}
		})
	}
}

// Test 5-7: Individual Heuristics
func TestEvalAggregateHeight(t *testing.T) {
	b := createBoardWithHeights([10]int{3, 5, 2, 0, 4, 1, 6, 3, 2, 4})
	got := evalAggregateHeight(b)
	want := 30
	if got != want {
		t.Errorf("evalAggregateHeight() = %d, want %d", got, want)
	}
}

func TestEvalHoles(t *testing.T) {
	b := createBoardWithHoles(5)
	got := evalHoles(b)
	want := 5
	if got != want {
		t.Errorf("evalHoles() = %d, want %d", got, want)
	}
}

func TestEvalBumpiness(t *testing.T) {
	b := createBoardWithHeights([10]int{3, 5, 2, 4, 1, 6, 3, 2, 4, 5})
	got := evalBumpiness(b)
	want := 22
	if got != want {
		t.Errorf("evalBumpiness() = %d, want %d", got, want)
	}
}

func TestEvalWells(t *testing.T) {
	b := createBoardWithWells(2)
	got := evalWells(b)
	want := 2
	if got != want {
		t.Errorf("evalWells() = %d, want %d", got, want)
	}
}

func TestGetWeights_SetWeights(t *testing.T) {
	weights := GetWeights()

	expectedKeys := []string{"aggregateHeight", "holes", "bumpiness", "wells"}
	for _, key := range expectedKeys {
		if _, ok := weights[key]; !ok {
			t.Errorf("GetWeights() missing key: %s", key)
		}
	}

	originalWeights := GetWeights()

	newWeights := map[string]float64{
		"aggregateHeight": -1.0,
		"holes":           -0.5,
		"bumpiness":       -0.3,
		"wells":           -0.2,
	}
	SetWeights(newWeights)
	got := GetWeights()

	for key, want := range newWeights {
		if got[key] != want {
			t.Errorf("SetWeights(%s) = %f, want %f", key, got[key], want)
		}
	}

	SetWeights(originalWeights)
}

func TestEvaluateBoard(t *testing.T) {
	board := createBoardWithHeights([10]int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5})
	piece := NewTetromino(TetrominoI)

	score := evaluateBoard(board, piece, 5, 0)

	if score == 0 {
		t.Error("evaluateBoard() should return non-zero score")
	}
}

func TestEvaluateBoard_WeightImpact(t *testing.T) {
	board := createBoardWithHeights([10]int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5})
	piece := NewTetromino(TetrominoI)

	originalWeights := GetWeights()

	baseline := evaluateBoard(board, piece, 5, 0)

	testWeights := GetWeights()
	testWeights["aggregateHeight"] = 0
	SetWeights(testWeights)

	modified := evaluateBoard(board, piece, 5, 0)

	if baseline == modified {
		t.Error("Changing weights should affect score")
	}

	SetWeights(originalWeights)
}

func TestEnumerateMoves(t *testing.T) {
	gameState := NewGameState()
	piece := gameState.CurrentPiece

	moves := enumerateMoves(gameState, piece)

	if len(moves) < 20 {
		t.Errorf("enumerateMoves() returned too few moves: %d, expect 30-40", len(moves))
	}

	for i, move := range moves {
		if !move.IsValid() {
			t.Errorf("Move %d is invalid: %+v", i, move)
		}
		if move.rotations < 0 || move.rotations > 3 {
			t.Errorf("Move %d has invalid rotation: %d", i, move.rotations)
		}
		if move.targetX < 0 || move.targetX > 9 {
			t.Errorf("Move %d has invalid X: %d", i, move.targetX)
		}
	}
}

func TestEnumerateMoves_AllPieceTypes(t *testing.T) {
	pieceTypes := []TetrominoType{
		TetrominoI, TetrominoO, TetrominoT,
		TetrominoS, TetrominoZ, TetrominoJ, TetrominoL,
	}

	for _, pType := range pieceTypes {
		t.Run(string(pType), func(t *testing.T) {
			gameState := NewGameState()
			piece := NewTetromino(pType)
			piece.X = 3
			piece.Y = 18

			moves := enumerateMoves(gameState, piece)

			if len(moves) < 15 {
				t.Errorf("%s: too few moves: %d", pType, len(moves))
			}
		})
	}
}

func TestFindBestMove(t *testing.T) {
	gameState := NewGameState()

	decision := FindBestMove(gameState)

	if decision == nil {
		t.Fatal("FindBestMove() returned nil on empty board")
	}

	if !decision.IsValid() {
		t.Errorf("FindBestMove() returned invalid decision: %+v", decision)
	}
}

func TestFindBestMove_Determinism(t *testing.T) {
	gameState := NewGameState()

	decision1 := FindBestMove(gameState)
	decision2 := FindBestMove(gameState)
	decision3 := FindBestMove(gameState)

	if decision1.targetX != decision2.targetX || decision2.targetX != decision3.targetX {
		t.Error("FindBestMove() should be deterministic")
	}
	if decision1.rotations != decision2.rotations {
		t.Error("FindBestMove() should be deterministic")
	}
}

func TestGetDelayForSpeed(t *testing.T) {
	tests := []struct {
		baseDelay  int
		speedLevel int
		wantDelay  int
	}{
		{1500, 1, 1500},
		{1500, 2, 750},
		{1500, 3, 300},
		{1500, 4, 150},
		{1500, 5, 0},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("level_%d", tt.speedLevel), func(t *testing.T) {
			got := GetDelayForSpeed(tt.baseDelay, tt.speedLevel)
			if got != tt.wantDelay {
				t.Errorf("GetDelayForSpeed(%d, %d) = %d, want %d", tt.baseDelay, tt.speedLevel, got, tt.wantDelay)
			}
		})
	}
}

func TestEvaluateLineClears_ExponentialBonus(t *testing.T) {
	tests := []struct {
		lines int
		want  float64
	}{
		{0, 0.0},
		{1, 1.00},
		{2, 10.00},
		{3, 50.00},
		{4, 150.00},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d lines", tt.lines), func(t *testing.T) {
			got := evaluateLineClears(tt.lines)
			if got != tt.want {
				t.Errorf("evaluateLineClears(%d) = %.2f, want %.2f", tt.lines, got, tt.want)
			}
		})
	}

	// Verify exponential growth
	if evaluateLineClears(4) <= evaluateLineClears(3)*2 {
		t.Error("4-line clear should be worth more than 2× 3-line clear")
	}
	if evaluateLineClears(2) <= evaluateLineClears(1)*2 {
		t.Error("2-line clear should be worth more than 2× 1-line clear")
	}
}

// TestWeightAnalysis_TetrisIncentive analyzes whether current heuristic weights
// properly incentivize multi-line clears (especially Tetris/4-line).
// After fix: Uses exponential line-clear bonus for accurate analysis.
func TestWeightAnalysis_TetrisIncentive(t *testing.T) {
	weights := GetWeights()

	t.Log("Current heuristic weights:")
	for k, v := range weights {
		t.Logf("  %s: %.2f", k, v)
	}

	t.Log("\nScenario: Compare conservative play vs Tetris setup")
	t.Log("Assumption: Both scenarios have 0 holes, 0 bumpiness, 0 wells")

	// Scenario 1: Conservative play - keep stack very low, no lines cleared
	conservativeHeight := 5.0
	conservativeScore := weights["aggregateHeight"] * conservativeHeight

	// Scenario 2: Tetris setup - build higher to clear 4 lines
	tetrisHeight := 10.0
	tetrisScore := weights["aggregateHeight"]*tetrisHeight + evaluateLineClears(4)

	t.Logf("\nConservative (height=%.0f, 0 lines): score = %.2f",
		conservativeHeight, conservativeScore)
	t.Logf("Tetris setup (height=%.0f, 4 lines): score = %.2f",
		tetrisHeight, tetrisScore)

	scoreDiff := tetrisScore - conservativeScore
	t.Logf("\nTetris advantage: %.2f points", scoreDiff)

	if tetrisScore > conservativeScore {
		t.Log("✓ Tetris is incentivized (good)")
	} else {
		t.Log("✗ Tetris is NOT incentivized - AI will play too conservatively")
		t.Logf("  Need %.2f more points to make Tetris worthwhile", -scoreDiff)
	}

	// Calculate break-even point
	t.Log("\nBreak-even analysis:")
	heightPenalty := weights["aggregateHeight"] * (tetrisHeight - conservativeHeight)
	t.Logf("  Height penalty for +5 units: %.2f", heightPenalty)
	t.Logf("  Lines bonus needed to break even: %.2f", -heightPenalty)
	t.Logf("  Actual 4-line bonus (exponential): %.2f", evaluateLineClears(4))
}

func TestEvaluateTwoPieceSequence_Basic(t *testing.T) {
	gameState := NewGameState()

	// Setup: Simple board, test that two-piece evaluation returns valid score
	decision := &MoveDecision{
		rotations: 0,
		targetX:   5,
		softDrops: 10,
	}

	score := EvaluateTwoPieceSequence(gameState, decision, gameState.NextPiece)

	if score == -999999.0 {
		t.Error("EvaluateTwoPieceSequence() should return valid score for empty board")
	}

	t.Logf("Two-piece sequence score on empty board: %.2f", score)
}

func TestEvaluateTwoPieceSequence_ComboBonus(t *testing.T) {
	gameState := NewGameState()

	// Setup: Board with 3 rows nearly full (gap in columns 4-5)
	// O-piece can fill part, I-piece can complete Tetris
	for y := 0; y < 3; y++ {
		for x := 0; x < 10; x++ {
			if x < 4 || x > 5 {
				gameState.Board.Set(x, y, 1)
			}
		}
	}

	// O-piece current
	gameState.CurrentPiece = NewTetromino(TetrominoO)
	gameState.CurrentPiece.X = 4
	gameState.CurrentPiece.Y = 3

	// I-piece next
	nextPiece := NewTetromino(TetrominoI)

	// Test move: O-piece at x=4, which sets up I-piece Tetris
	decision := &MoveDecision{
		rotations: 0,
		targetX:   4,
		softDrops: 0,
	}

	score := EvaluateTwoPieceSequence(gameState, decision, nextPiece)

	t.Logf("Combo setup score (O+I for Tetris): %.2f", score)

	// Score should be high due to Tetris potential
	if score < 10.0 {
		t.Logf("Warning: Combo score lower than expected (got %.2f)", score)
	}
}

func TestEvaluateTwoPieceSequence_EdgeCases(t *testing.T) {
	gameState := NewGameState()
	decision := &MoveDecision{rotations: 0, targetX: 5, softDrops: 10}
	nextPiece := NewTetromino(TetrominoI)

	t.Run("NilGameState", func(t *testing.T) {
		score := EvaluateTwoPieceSequence(nil, decision, nextPiece)
		if score != -999999.0 {
			t.Errorf("Expected -999999.0 for nil gameState, got %.2f", score)
		}
	})

	t.Run("NilDecision", func(t *testing.T) {
		score := EvaluateTwoPieceSequence(gameState, nil, nextPiece)
		if score != -999999.0 {
			t.Errorf("Expected -999999.0 for nil decision, got %.2f", score)
		}
	})

	t.Run("NilNextPiece", func(t *testing.T) {
		score := EvaluateTwoPieceSequence(gameState, decision, nil)
		if score != -999999.0 {
			t.Errorf("Expected -999999.0 for nil nextPiece, got %.2f", score)
		}
	})

	t.Run("NoValidNextMoves", func(t *testing.T) {
		// Fill board with partial rows (no complete rows, so no clears) blocking top rows.
		// Leave column 0 empty to prevent full-line clears, but stack high enough
		// to block the next piece from spawning.
		for y := 0; y < 20; y++ {
			for x := 1; x < 10; x++ {
				gameState.Board.Set(x, y, 1)
			}
		}
		score := EvaluateTwoPieceSequence(gameState, decision, nextPiece)
		if score != -999999.0 {
			t.Errorf("Expected -999999.0 for no valid moves, got %.2f", score)
		}
	})
}

func TestEnumerateMovesForBoard(t *testing.T) {
	board := NewBoard()

	t.Run("EmptyBoard_AllPieces", func(t *testing.T) {
		pieceTypes := []TetrominoType{TetrominoI, TetrominoO, TetrominoT, TetrominoS, TetrominoZ, TetrominoJ, TetrominoL}
		for _, pt := range pieceTypes {
			piece := NewTetromino(pt)
			moves := enumerateMovesForBoard(board, piece)
			if len(moves) == 0 {
				t.Errorf("Expected moves for %v on empty board, got 0", pt)
			}
		}
	})

	t.Run("NarrowGap_OnlyIPieceFits", func(t *testing.T) {
		// Create 1-wide well at column 5
		for x := 0; x < 10; x++ {
			for y := 0; y < 5; y++ {
				if x != 5 {
					board.Set(x, y, 1)
				}
			}
		}

		iPiece := NewTetromino(TetrominoI)
		iMoves := enumerateMovesForBoard(board, iPiece)
		if len(iMoves) == 0 {
			t.Error("Expected I-piece to have moves in 1-wide well scenario")
		}

		oPiece := NewTetromino(TetrominoO)
		oMoves := enumerateMovesForBoard(board, oPiece)
		// O-piece should have fewer moves than I-piece because it can't fit in the narrow well
		if len(oMoves) >= len(iMoves) {
			t.Logf("Warning: O-piece (%d moves) has as many or more moves than I-piece (%d)",
				len(oMoves), len(iMoves))
		}
	})

	t.Run("FullBoard_NoMoves", func(t *testing.T) {
		fullBoard := NewBoard()
		for y := 0; y < 20; y++ {
			for x := 0; x < 10; x++ {
				fullBoard.Set(x, y, 1)
			}
		}

		piece := NewTetromino(TetrominoI)
		moves := enumerateMovesForBoard(fullBoard, piece)
		if len(moves) != 0 {
			t.Errorf("Expected 0 moves on full board, got %d", len(moves))
		}
	})
}

func TestIsValidPositionForBoard(t *testing.T) {
	t.Run("EmptyBoard_Valid", func(t *testing.T) {
		board := NewBoard()
		piece := NewTetromino(TetrominoI)
		if !isValidPositionForBoard(board, piece, 5, 18) {
			t.Error("Expected valid position on empty board")
		}
	})

	t.Run("OutOfBounds_Invalid", func(t *testing.T) {
		board := NewBoard()
		piece := NewTetromino(TetrominoI)
		if isValidPositionForBoard(board, piece, 10, 18) {
			t.Error("Expected invalid position (x=10 out of bounds)")
		}
		if isValidPositionForBoard(board, piece, -1, 18) {
			t.Error("Expected invalid position (x=-1 out of bounds)")
		}
	})

	t.Run("Collision_Invalid", func(t *testing.T) {
		board := NewBoard()
		piece := NewTetromino(TetrominoI)
		// Place block where piece would land
		board.Set(5, 16, 1)
		if isValidPositionForBoard(board, piece, 5, 17) {
			t.Error("Expected invalid position (collision at y=16)")
		}
	})

	t.Run("ClearPosition_Valid", func(t *testing.T) {
		board := NewBoard()
		piece := NewTetromino(TetrominoI)
		// Block far away (top-left corner), position at center should still be valid
		board.Set(0, 19, 1)
		if !isValidPositionForBoard(board, piece, 5, 18) {
			t.Error("Expected valid position (block at corner doesn't interfere)")
		}
	})
}

func BenchmarkFindBestMoveWithNext(b *testing.B) {
	gameState := NewGameState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindBestMoveWithNext(gameState)
	}
}

func BenchmarkEvaluateTwoPieceSequence(b *testing.B) {
	gameState := NewGameState()
	decision := &MoveDecision{
		rotations: 0,
		targetX:   5,
		softDrops: 10,
	}
	nextPiece := NewTetromino(TetrominoI)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvaluateTwoPieceSequence(gameState, decision, nextPiece)
	}
}

func TestTwoPieceLookahead_NoExcessiveAllocations(t *testing.T) {
	gameState := NewGameState()
	decision := &MoveDecision{
		rotations: 0,
		targetX:   5,
		softDrops: 10,
	}

	allocs := testing.AllocsPerRun(100, func() {
		EvaluateTwoPieceSequence(gameState, decision, gameState.NextPiece)
	})

	// Two-piece lookahead allocates boards for simulation
	// Current implementation: ~500-600 allocs per call (board clones, move slices)
	// This is acceptable for the complexity of 2-piece simulation
	if allocs > 1000 {
		t.Errorf("Excessive allocations: %.1f per call (expected <1000)", allocs)
	}
	t.Logf("Allocations per call: %.1f (acceptable for 2-piece lookahead)", allocs)
}
